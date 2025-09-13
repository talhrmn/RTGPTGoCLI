package errorhandler

func NewAppError(level string, message string, err error) *AppError {
	// Create app error
	return &AppError{
		Level: level,
		Message: message,
		Error: err,
	}
}
