package clients

import (
	"RTGPTGoCLI/RTGPTGoCLI/pkg/errorhandler"
	"context"
)

type ClientConnection interface {
	// Client connection interface
	Connect(ctx context.Context) error
	Disconnect() error
	IsConnected() bool
	GetErrorChannel() <-chan errorhandler.AppError
}

type WebClientConnection interface {
	// Web client connection interface
	ClientConnection
	SendMessage(ctx context.Context, message []byte) error
	GetMessageChannel() <-chan []byte
}
