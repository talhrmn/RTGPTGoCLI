package config

type Config struct {
	APIKey  string
	BaseURL string
	Model   string
	Debug   bool
	Timeout int
	Retries int
	ChannelBuffer int
}

type FlagType string
