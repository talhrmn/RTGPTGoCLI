package app

import (
	"RTGPTGoCLI/RTGPTGoCLI/internal/cli"
	"RTGPTGoCLI/RTGPTGoCLI/internal/clients/openai"
	"RTGPTGoCLI/RTGPTGoCLI/internal/config"
	"RTGPTGoCLI/RTGPTGoCLI/pkg/errorhandler"
	"context"
)

type App struct {
	// App struct
	config    *config.Config
	cli       *cli.CLI
	oaiClient openai.OpenAIClientInterface
	ctx       context.Context
	cancel    context.CancelFunc
	errorHandler *errorhandler.ErrorHandler
}
