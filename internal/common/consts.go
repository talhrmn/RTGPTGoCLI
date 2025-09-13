package common

import (
	"RTGPTGoCLI/pkg/errorhandler"
	"encoding/json"
)

const (
	UnexpectedError = "unexpected error: %v"
	JsonUnmarshalError = "failed to unmarshal JSON: %v"
	JsonMarshalError = "failed to marshal JSON: %v"
)

func NewErrJsonUnmarshalAppError(err error) *errorhandler.AppError {
	return errorhandler.NewAppError(errorhandler.WarningLevel, JsonUnmarshalError, err)
}

func NewErrJsonMarshalAppError(err error) *errorhandler.AppError {
	return errorhandler.NewAppError(errorhandler.WarningLevel, JsonMarshalError, &json.MarshalerError{})
}
