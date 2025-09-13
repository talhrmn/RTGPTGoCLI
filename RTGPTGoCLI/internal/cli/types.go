package cli

import (
	"RTGPTGoCLI/RTGPTGoCLI/internal/clients/openai"
	"RTGPTGoCLI/RTGPTGoCLI/internal/config"
	"RTGPTGoCLI/RTGPTGoCLI/pkg/errorhandler"
	"bufio"
)

type CLI struct {
	// Command Line Interface struct
	config    *config.Config
	scanner   *bufio.Scanner
	oaiClient openai.OpenAIClientInterface
	streamingChannel chan struct{}
	errorHandler *errorhandler.ErrorHandler
}
