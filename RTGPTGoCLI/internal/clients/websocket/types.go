package websocket

import (
	"RTGPTGoCLI/RTGPTGoCLI/internal/clients"
	"RTGPTGoCLI/RTGPTGoCLI/internal/config"
	"RTGPTGoCLI/RTGPTGoCLI/pkg/errorhandler"
	"context"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type WebSocketClientInterface interface {
	// WebSocketClient interface
	clients.WebClientConnection
}

type WebSocketClient struct {
	// Websocket client struct
	config     *config.Config
	connection *websocket.Conn
	reconnectOnce     sync.Once
	cleanUpOnce       sync.Once
	cancel context.CancelFunc

	url        string
	headers    http.Header

	mu             sync.RWMutex
	sendChannel    chan []byte
	messageChannel chan []byte
	errorChannel   chan errorhandler.AppError
	connected      bool
}
