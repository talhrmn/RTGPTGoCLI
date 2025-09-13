package app

import (
	"RTGPTGoCLI/internal/cli"
	"RTGPTGoCLI/internal/clients/openai"
	"RTGPTGoCLI/internal/clients/websocket"
	"RTGPTGoCLI/internal/config"
	"RTGPTGoCLI/pkg/errorhandler"
	"RTGPTGoCLI/pkg/logger"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func New(cfg *config.Config) *App {
	// Create app
	ctx, cancel := context.WithCancel(context.Background())

	wsc := websocket.NewWebSocketClient(cfg)
	oaiClient := openai.NewOAIClient(cfg, wsc)
	cli := cli.New(cfg, oaiClient)
	errHandler := errorhandler.NewErrorHandler(cfg.Debug)
	
	return &App{
		config:    cfg,
		cli:       cli,
		oaiClient: oaiClient,
		ctx:       ctx,
		cancel:    cancel,
		errorHandler: errHandler,
	}
}

func (app *App) Run() error {
	// Run app
	app.handleShutdown()

	if err := app.oaiClient.Connect(app.ctx); err != nil {
		appErr := *errorhandler.NewAppError(errorhandler.ErrorLevel, fmt.Sprintf(AppFailedToConnectToOAIErr, err), err)
		app.errorHandler.HandleError(appErr)
		return err
	}

	app.cli.Run(app.ctx, app.cancel)
	return nil
}

func (app *App) handleShutdown() {
	// Handle app shutdown
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChannel
		logger.Debug(AppShutdownMsg)
		app.cancel()

		if err := app.oaiClient.Disconnect(); err != nil {
			appErr := *errorhandler.NewAppError(errorhandler.ErrorLevel, fmt.Sprintf(AppFailedToDisconnectFromOAIErr, err), err)
			app.errorHandler.HandleError(appErr)
		}
	}()
}
