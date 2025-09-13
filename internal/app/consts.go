package app

const (
	// Errors
	AppFailedToConnectToOAIErr = "Failed to connect to OpenAI: %v"
	AppFailedToDisconnectFromOAIErr = "Failed to disconnect from OpenAI: %v"
	AppErrorClosingOAIConnErr = "Error closing OAI client: %v"
)

const (
	// Messages
	AppShutdownMsg = "Received shutdown signal, exiting gracefully..."
)
