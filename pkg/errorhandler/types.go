package errorhandler

// ErrorHandler struct
type ErrorHandler struct {
	debug bool
}

// AppError struct
type AppError struct {
	Level string
	Message string
	Code string
	Error error
}

