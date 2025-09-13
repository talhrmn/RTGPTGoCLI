package config

const (
	// Flag names
	ApiKeyFlag  FlagType = "api-key"
	BaseURLFlag FlagType = "base-url"
	ModelFlag   FlagType = "model"
	DebugFlag   FlagType = "debug"
	TimeoutFlag FlagType = "timeout"
	RetriesFlag FlagType = "retries"
	ChannelBufferFlag FlagType = "channel-buffer"
)

const (
	// Default values
	DefaultAPIKey = ""
	DefaultBaseURL = "api.openai.com/v1/realtime"
	DefaultTimeout = 30
	DefaultModel = "gpt-4o-mini-realtime-preview"
	DefaultDebug = false
	DefaultRetries = 3
	DefaultChannelBuffer = 100
)

const (
	// Error messages
	MissingOrErrorLoadingEnvFileErr = "no .env file found or error loading .env file, proceeding with existing environment variables."
	MissingRequiredFlagsOrEnvVarsErr = "missing the following required flags or environment variables: %v"
)

const (
	// Flag usage strings
	APIKeyFlagUsageText = "API key for authentication"
	BaseURLFlagUsageText = "Base URL for the API"
	ModelFlagUsageText = "Model to use for requests"
	TimeoutFlagUsageText = "Request timeout in seconds"
	RetriesFlagUsageText = "Number of retries for failed requests"
	ChannelBufferFlagUsageText = "Buffer size for channels"
	DebugFlagUsageText = "Enable debug mode"
)

const (
	// General strings
	ConfigInfoText = "\nCurrent Configuration:\n"
)
