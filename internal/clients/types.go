package clients

import (
	"RTGPTGoCLI/pkg/errorhandler"
	"context"
)

type MessageEvent struct {
	// Message event struct
	Type string
	Text string
	Done bool
}

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

type ServiceClientConnection interface {
	// Service client connection interface
	ClientConnection
	GetMessageChannel() <-chan MessageEvent
	SendMessage(ctx context.Context, message string) *errorhandler.AppError
}
