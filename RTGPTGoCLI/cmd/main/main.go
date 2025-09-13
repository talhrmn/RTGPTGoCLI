package main

import (
	"RTGPTGoCLI/RTGPTGoCLI/internal/config"
	"RTGPTGoCLI/RTGPTGoCLI/pkg/errorhandler"
	"RTGPTGoCLI/RTGPTGoCLI/pkg/logger"
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

	errorHandler := errorhandler.NewErrorHandler(cfg.Debug)

	logger.SetDebugMode(cfg.Debug)
	if cfg.Debug {
		if cfgString, err := cfg.GetConfigInfo(); err != nil {
			appErr := errorhandler.NewAppError(
				errorhandler.WarningLevel,
				fmt.Sprintf(FailedToGetConfigInfoErr, err),
				err,
			)
			errorHandler.HandleError(*appErr)
		} else {
			logger.Debug(cfgString)
		}
	}

	logger.Info(CLIExitedSuccessfullyMsg)
}
