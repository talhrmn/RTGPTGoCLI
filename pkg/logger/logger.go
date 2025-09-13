package logger

import (
	"log"
	"os"
)

func InitLoggers() {
	// Initialize loggers
	infoLogger = log.New(os.Stdout, INFO_PREFIX, log.LstdFlags)
	debugLogger = log.New(os.Stdout, DEBUG_PREFIX, log.LstdFlags)
	warningLogger = log.New(os.Stderr, WARNING_PREFIX, log.LstdFlags)
	errorLogger = log.New(os.Stderr, ERROR_PREFIX, log.LstdFlags)
	isDebugMode = false
}

func SetDebugMode(debug bool) {
	// Set debug mode
	isDebugMode = debug
}

func Info(msg string) {
	// Print info log
	infoLogger.Println(msg)
}

func Debug(msg string) {
	// Print debug log, only if debug mode is enabled
	if isDebugMode {
		debugLogger.Println(msg)
	}
}

func Warning(msg string) {
	// Print warning log, only if debug mode is enabled
	if isDebugMode {
		warningLogger.Println(msg)
	}
}

func Error(msg string) {
	// Print error log
	errorLogger.Println(msg)
}
