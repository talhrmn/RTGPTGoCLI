package main

import (
	"RTGPTGoCLI/internal/app"
	"RTGPTGoCLI/internal/config"
	"RTGPTGoCLI/pkg/logger"
	"fmt"
	"os"
)

func main() {
	// Main program
	logger.InitLoggers()
	
	cfg, err := config.SetUp()
	if err != nil {
		logger.Error(fmt.Sprintf(FailedToLoadConfigErr, err))
		os.Exit(1)
	}

	logger.SetDebugMode(cfg.Debug)
	if cfg.Debug {
		if cfgString, err := cfg.GetConfigInfo(); err != nil {
			logger.Warning(fmt.Sprintf(FailedToGetConfigInfoErr, err))
		} else {
			logger.Debug(cfgString)
		}
	}

	application := app.New(cfg)
	if err := application.Run(); err != nil {
		logger.Error(fmt.Sprintf(FailedToRunApplicationErr, err))
		os.Exit(1)
	}

	logger.Info(CLIExitedSuccessfullyMsg)
}
