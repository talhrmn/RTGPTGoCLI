package errorhandler

import (
	"RTGPTGoCLI/RTGPTGoCLI/pkg/logger"
	"os"
)

func NewErrorHandler(debug bool) *ErrorHandler {
	logger.LoggerInit()
	logger.SetDebugMode(debug)
	return &ErrorHandler{
		debug: debug,
	}
}

func (h *ErrorHandler) HandleError(appErr AppError) {
	errString := h.getErrorString(appErr)
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

func NewAppError(level string, message string, err error) *AppError {
	return &AppError{
		Level: level,
		Message: message,
		Error: err,
	}
}

func (h *ErrorHandler) getErrorString(appErr AppError) string {
	return appErr.Message + "\n" + appErr.Error.Error()
}
