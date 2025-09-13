package handler

import (
	"RTGPTGoCLI/RTGPTGoCLI/internal/common"
	"RTGPTGoCLI/RTGPTGoCLI/internal/functions"
	"RTGPTGoCLI/RTGPTGoCLI/internal/functions/math"
	"RTGPTGoCLI/RTGPTGoCLI/pkg/errorhandler"
	"context"
	"encoding/json"
	"fmt"
)

func NewHandler() *FunctionHandler {
	return &FunctionHandler{
		functions: make(FunctionsType),
	}
}

func (fh *FunctionHandler) LoadFunctions() *errorhandler.AppError {
	// Load all custom functions to handler

	functionsToLoad := []functions.FunctionInterface{
		&math.FunctionMultiply{},
		// Add more functions here as you create them
	}

	for _, fn := range functionsToLoad {
		if appErr := fh.registerFunction(fn); appErr != nil {
			return appErr
		}
	}
	return nil
}

func (fh *FunctionHandler) GenerateOpenAITools() []interface{} {
	// Generate OpenAI tools for session configuration
	tools := make([]interface{}, 0, len(fh.functions))
	for _, fn := range fh.functions {
		tools = append(tools, fn.ConvertToOpenAITool())
	}
	return tools
}

func (fh *FunctionHandler) Execute(ctx context.Context, name string, argumentsJSON string) (interface{}, *errorhandler.AppError) {
	// Execute function
	fn, err := fh.getFunction(name)
	if err != nil {
		return nil, err
	}

	var params functions.FunctionParams
	if err := json.Unmarshal([]byte(argumentsJSON), &params); err != nil {
		return nil, common.NewErrJsonUnmarshalAppError(err)
	}

	result, appErr := fn.Execute(ctx, params)
	return result, appErr
}

func (fh *FunctionHandler) getFunction(name string) (functions.FunctionInterface, *errorhandler.AppError) {
	// Get function from handler
	fn, exists := fh.functions[name]
	if !exists {
		return nil, errorhandler.NewAppError(errorhandler.WarningLevel, fmt.Sprintf(FunctionDoesntExistsErr, name), nil)
	}
	
	return fn, nil
}

func (fh *FunctionHandler) registerFunction(fn functions.FunctionInterface) *errorhandler.AppError {
	// Register function to handler
	functionName := fn.GetMetadata().Name

	if functionName == "" {
		return errorhandler.NewAppError(errorhandler.WarningLevel, FunctionEmptyNameErr, nil)
	}

	if _, exists := fh.functions[functionName]; exists {
		return errorhandler.NewAppError(errorhandler.WarningLevel, fmt.Sprintf(FunctionAlreadyExistsErr, functionName), nil)
	}

	fh.functions[functionName] = fn
	return nil
}
