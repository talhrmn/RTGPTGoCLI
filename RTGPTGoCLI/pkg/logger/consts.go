package logger

import "log"

// Logger prefix constants
const (
	INFO_PREFIX    = "[INFO] "
	DEBUG_PREFIX   = "[DEBUG] "
	WARNING_PREFIX = "[WARNING] "
	ERROR_PREFIX   = "[ERROR] "
)

// Logger variables
var (
	infoLogger    *log.Logger
	debugLogger   *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
	isDebugMode   bool
)
