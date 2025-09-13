package websocket

const (
	// Websocket client constants
	WSAuthHeader   = "Authorization"
	WSBearerPrefix = "Bearer "
	WSUrlBuild     = "wss://%s?model=%s"
)

const (
	// Websocket client error messages
	ConnectionErr = "failed to connect to websocket: %v"
	ConnectionIsClosedErr = "websocket connection is closed"
	ClientClosedErr = "websocket client closed"
	CloseErr = "websocket close error: %v"
	ReadErr = "websocket read error: %v"
	WriteErr = "websocket write error: %v"
	ReconnectionAttemptFailedErr = "reconnection attempt %d failed: %v\n"
	ReconnectionFailedErr = "failed to reconnect after %d attempts"
)

const (
	// Websocket client info messages
	ConnectedMsg = "WebSocket connected to the server."
	DisconnectedMsg = "WebSocket disconnected from the server."
	ReconnectingMsg = "Reconnecting to the server..."
	ReconnectionSuccessMsg = "Reconnection successful"
)
