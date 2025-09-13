package app

import (
	"RTGPTGoCLI/internal/cli"
	"RTGPTGoCLI/internal/clients/openai"
	"RTGPTGoCLI/internal/config"
	"RTGPTGoCLI/pkg/errorhandler"
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
