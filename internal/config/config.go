package config

import (
	"RTGPTGoCLI/pkg/logger"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

func SetUp() (*Config, error) {
	// Initialize config
	cfg := &Config{}

	if err := cfg.loadVars(); err != nil {
		return nil, err
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (cfg *Config) loadVars() error {
	// Load config from environment variables and flags
	cfg.loadFromEnv()
	cfg.loadFromFlags()

	return nil
}

func (cfg *Config) loadFromEnv() {
	// Load config from environment variables
	cfg.setDefaults()
	cfg.loadEnvVars()
}

func (cfg *Config) setDefaults() {
	// Set default values
	cfg.APIKey = DefaultAPIKey
	cfg.BaseURL = DefaultBaseURL
	cfg.Timeout = DefaultTimeout
	cfg.Model = DefaultModel
	cfg.Debug = DefaultDebug
	cfg.Retries = DefaultRetries
	cfg.ChannelBuffer = DefaultChannelBuffer
}

func (cfg *Config) loadEnvVars() {
	// Load config from environment variables
	err := godotenv.Load()
	if err != nil {
		logger.Warning(MissingOrErrorLoadingEnvFileErr)
	}
	cfg.setStringEnvVar(ApiKeyFlag, &cfg.APIKey)
	cfg.setStringEnvVar(BaseURLFlag, &cfg.BaseURL)
	cfg.setStringEnvVar(ModelFlag, &cfg.Model)

	cfg.setIntEnvVar(TimeoutFlag, &cfg.Timeout)
	cfg.setIntEnvVar(RetriesFlag, &cfg.Retries)
	cfg.setIntEnvVar(ChannelBufferFlag, &cfg.ChannelBuffer)

	cfg.setBoolEnvVar(DebugFlag, &cfg.Debug)
}

func (cfg *Config) loadFromFlags() {
	// Load config from flags
	flag.StringVar(&cfg.APIKey, string(ApiKeyFlag), cfg.APIKey, APIKeyFlagUsageText)
	flag.StringVar(&cfg.BaseURL, string(BaseURLFlag), cfg.BaseURL, BaseURLFlagUsageText)
	flag.StringVar(&cfg.Model, string(ModelFlag), cfg.Model, ModelFlagUsageText)

	flag.IntVar(&cfg.Timeout, string(TimeoutFlag), cfg.Timeout, TimeoutFlagUsageText)
	flag.IntVar(&cfg.Retries, string(RetriesFlag), cfg.Retries, RetriesFlagUsageText)
	flag.IntVar(&cfg.ChannelBuffer, string(ChannelBufferFlag), cfg.ChannelBuffer, ChannelBufferFlagUsageText)

	flag.BoolVar(&cfg.Debug, string(DebugFlag), cfg.Debug, DebugFlagUsageText)
	flag.Parse()
}

func (cfg *Config) validate() error {
	// Validate only the required flags, that they exist and are valid
	requiredKeys := map[FlagType]*string{
		ApiKeyFlag:  &cfg.APIKey,
		BaseURLFlag: &cfg.BaseURL,
		ModelFlag:   &cfg.Model,
	}

	missingVars := []FlagType{}

	for name, valuePtr := range requiredKeys {
		if *valuePtr == "" {
			missingVars = append(missingVars, name)
		}
	}

	if len(missingVars) > 0 {
		return fmt.Errorf(MissingRequiredFlagsOrEnvVarsErr, missingVars)
	}

	return nil
}

func (flagName FlagType) envVar() string {
	// Convert flag name to environment variable name
	return strings.ToUpper(strings.ReplaceAll(string(flagName), "-", "_"))
}


func (cfg *Config) setStringEnvVar(key FlagType, cfgPtr *string) {
	// Set string environment variable value to flag
	envKey := key.envVar()
	if envVal := os.Getenv(envKey); envVal != "" {
		*cfgPtr = envVal
	}
}

func (cfg *Config) setIntEnvVar(key FlagType, cfgPtr *int) {
	// Set int environment variable value to flag
	envKey := key.envVar()
	if envVal := os.Getenv(envKey); envVal != "" {
		if intVal, _ := strconv.Atoi(envVal); intVal > 0 {
			*cfgPtr = intVal
		}
	}
}

func (cfg *Config) setBoolEnvVar(key FlagType, cfgPtr *bool) {
	// Set bool environment variable value to flag
	envKey := key.envVar()
	if envVal := os.Getenv(envKey); envVal != "" {
		switch envVal {
		case "true", "1":
			*cfgPtr = true
		case "false", "0":
			*cfgPtr = false
		}
	}
}

func (cfg *Config) GetConfigInfo() (string, error) {
	configString := strings.Builder{}
	configString.WriteString(ConfigInfoText)
	jsonBytes, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return "", err
	}
	configString.WriteString(string(jsonBytes))
	configString.WriteString("\n")
	return configString.String(), nil
}
