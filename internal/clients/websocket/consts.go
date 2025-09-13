package websocket

const (
	// Websocket client constants
	WSAuthHeader   = "Authorization"
	WSBearerPrefix = "Bearer "
	WSUrlBuild     = "wss://%s?model=%s"
)

const (
	// Websocket client error messages
	WSConnectionErr = "failed to connect to websocket: %v"
	WSConnectionIsClosedErr = "websocket connection is closed"
	WSClientClosedErr = "websocket client closed"
	WSCloseErr = "websocket close error: %v"
	WSReadErr = "websocket read error: %v"
	WSWriteErr = "websocket write error: %v"
	WSReconnectionAttemptFailedErr = "reconnection attempt %d failed: %v\n"
	WSReconnectionFailedErr = "failed to reconnect after %d attempts"
)

const (
	// Websocket client info messages
	WSConnectedMsg = "WebSocket connected to the server."
	WSDisconnectedMsg = "WebSocket disconnected from the server."
	WSReconnectingMsg = "Reconnecting to the server..."
	WSReconnectionSuccessMsg = "Reconnection successful"
)
