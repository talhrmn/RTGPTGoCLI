package websocket

import (
	"RTGPTGoCLI/internal/config"
	"RTGPTGoCLI/pkg/errorhandler"
	"RTGPTGoCLI/pkg/logger"
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

func NewWebSocketClient(cfg *config.Config) *WebSocketClient {
	// Initialize websocket client
	return &WebSocketClient{
		config:     cfg,
		connection: nil,
		cancel: nil,
		reconnectOnce:     sync.Once{},
		cleanUpOnce:       sync.Once{},

		url:        fmt.Sprintf(WSUrlBuild, cfg.BaseURL, cfg.Model),
		headers: http.Header{
			WSAuthHeader: []string{WSBearerPrefix + cfg.APIKey},
		},

		mu:                sync.RWMutex{},
		sendChannel:       make(chan []byte, cfg.ChannelBuffer),
		messageChannel:    make(chan []byte, cfg.ChannelBuffer),
		errorChannel:      make(chan errorhandler.AppError, cfg.ChannelBuffer),
		connected:         false,
	}
}

func (wsc *WebSocketClient) Connect(ctx context.Context) error {
	// Connect to websocket
	if err := wsc.connectOrRetry(ctx); err != nil {
		return fmt.Errorf(WSConnectionErr, err)
	}
	return nil
}


func (wsc *WebSocketClient) Disconnect() error {
	// Disconnect from websocket
    var disconnectErr error
    
    wsc.cleanUpOnce.Do(func() {
        logger.Debug(WSDisconnectedMsg)
        
        if !wsc.IsConnected() {
            return
        }
		
        if wsc.cancel != nil {
            wsc.cancel()
        }

        wsc.mu.Lock()
        defer wsc.mu.Unlock()

        close(wsc.sendChannel)
        close(wsc.messageChannel)
        close(wsc.errorChannel)

        
        if wsc.connection != nil {
			// Close WebSocket connection
            if err := wsc.connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
                disconnectErr = fmt.Errorf(WSCloseErr, err)
            }

            if err := wsc.connection.Close(); err != nil {
                disconnectErr = fmt.Errorf(WSCloseErr, err)
            }
        }
        
        wsc.setConnected(false)
    })
    
    return disconnectErr
}

func (wsc *WebSocketClient) IsConnected() bool {
	// Return if websocket is connected
	wsc.mu.RLock()
	defer wsc.mu.RUnlock()
	return wsc.connected
}

func (wsc *WebSocketClient) GetErrorChannel() <-chan errorhandler.AppError {
	// Return error channel
	return wsc.errorChannel
}

func (wsc *WebSocketClient) SendMessage(ctx context.Context, message []byte) error {
	// Send message to websocket
	if !wsc.IsConnected() {
		return errors.New(WSConnectionIsClosedErr)
	}

	select {
	case <-ctx.Done():
		return nil
	case wsc.sendChannel <- message:
		return nil
	}
}

func (wsc *WebSocketClient) GetMessageChannel() <-chan []byte {
	// Return message channel
	return wsc.messageChannel
}

func (wsc *WebSocketClient) connectOrRetry(ctx context.Context) error {
	// Connect (or retry connecting) to the WebSocket server
	if wsc.cancel != nil {
		wsc.cancel()
		time.Sleep(time.Duration(wsc.config.Timeout) * time.Second)
	}

	connectionContext, cancel := context.WithCancel(ctx)
	wsc.cancel = cancel

	if wsc.connection != nil {
		wsc.connection.Close()
	}

	connection, _, err := websocket.DefaultDialer.Dial(wsc.url, wsc.headers)
	if err != nil {
		return err
	}

	wsc.connection = connection
	wsc.setConnected(true)
	logger.Debug(WSConnectedMsg)

	go wsc.readRoutine(connectionContext)
	go wsc.writeRoutine(connectionContext)

	return nil
}

func (wsc *WebSocketClient) setConnected(state bool) {
	// Set connection state
	wsc.mu.Lock()
	defer wsc.mu.Unlock()
	wsc.connected = state
}

func (wsc *WebSocketClient) readRoutine(ctx context.Context) {
	// Read messages from websocket in go routine
	defer func() {
		if r := recover(); r != nil {
			wsc.setConnected(false)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			if !wsc.IsConnected() {
				wsc.errorChannel <- *errorhandler.NewAppError(errorhandler.WarningLevel, WSClientClosedErr, errors.New(WSClientClosedErr))
				return
			}

			_, response, err := wsc.connection.ReadMessage()
			if err != nil {
				wsc.setConnected(false)
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					wsc.errorChannel <- *errorhandler.NewAppError(errorhandler.WarningLevel, WSConnectionIsClosedErr, errors.New(WSConnectionIsClosedErr))
					return
				}

				wsc.reconnectOnce.Do(func() {
					go wsc.handleReconnection(context.Background())
				})
				return
			}

			select {
			case wsc.messageChannel <- response:
			case <-ctx.Done():
				return
			default:
				continue
			}
		}
	}
}

func (wsc *WebSocketClient) writeRoutine(ctx context.Context) {
	// Write messages to websocket in go routine
	defer func() {
		if r := recover(); r != nil {
			wsc.setConnected(false)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case prompt := <-wsc.sendChannel:
			if !wsc.IsConnected() {
				wsc.errorChannel <- errorhandler.AppError{
					Level: errorhandler.WarningLevel,
					Message: WSClientClosedErr,
					Error: errors.New(WSClientClosedErr),
				}
				continue
			}

			err := wsc.connection.WriteMessage(websocket.TextMessage, prompt)
			if err != nil {
				wsc.errorChannel <- errorhandler.AppError{
					Level: errorhandler.WarningLevel,
					Message: WSWriteErr,
					Error: err,
				}
				wsc.setConnected(false)

				wsc.reconnectOnce.Do(func() {
					go wsc.handleReconnection(context.Background())
				})
				return
			}
		}
	}
}

func (wsc *WebSocketClient) handleReconnection(ctx context.Context) {
	logger.Debug(WSReconnectingMsg)
	wsc.setConnected(false)

	// Add a small delay before starting reconnection attempts
	time.Sleep(time.Duration(wsc.config.Timeout) * time.Millisecond)

	for attempt := 1; attempt <= wsc.config.Retries; attempt++ {
		delay := time.Duration(attempt) * time.Duration(wsc.config.Timeout) * time.Second
		if delay > time.Duration(wsc.config.Timeout) * time.Second {
			delay = time.Duration(wsc.config.Timeout) * time.Second
		}

		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(delay)
		}

		if err := wsc.connectOrRetry(ctx); err != nil {
			wsc.errorChannel <- *errorhandler.NewAppError(errorhandler.WarningLevel, fmt.Sprintf(WSReconnectionAttemptFailedErr, attempt, err), err)
			continue
		}

		logger.Debug(WSReconnectionSuccessMsg)
		// Reset reconnection flag for next time
		wsc.reconnectOnce = sync.Once{}
		return
	}

	// Failed to reconnect after all attempts
	select {
	case wsc.errorChannel <- *errorhandler.NewAppError(errorhandler.ErrorLevel, fmt.Sprintf(WSReconnectionFailedErr, wsc.config.Retries), fmt.Errorf(WSReconnectionFailedErr, wsc.config.Retries)):
	default:
	}
}
