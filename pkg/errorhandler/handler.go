package errorhandler

import (
	"RTGPTGoCLI/pkg/logger"
	"os"
)

func NewErrorHandler(debug bool) *ErrorHandler {
	// Create error handler
	logger.InitLoggers()
	logger.SetDebugMode(debug)
	return &ErrorHandler{
		debug: debug,
	}
}

func (eh *ErrorHandler) HandleError(appErr AppError) {
	// Handle error
	errString := eh.getErrorString(appErr)
	switch appErr.Level {
	case InfoLevel:
		logger.Info(errString)
	case DebugLevel:
		logger.Debug(errString)
	case WarningLevel:
		logger.Warning(errString)
	case ErrorLevel:
		logger.Error(errString)
		os.Exit(1)
	}
}

func (eh *ErrorHandler) getErrorString(appErr AppError) string {
	// Get app error string
	return appErr.Message + "\n" + appErr.Error.Error()
}
