package math

import (
	"RTGPTGoCLI/internal/functions"
	"RTGPTGoCLI/pkg/errorhandler"
	"RTGPTGoCLI/pkg/logger"
	"context"
	"fmt"
)

func (mtpFn *FunctionMultiply) Execute(ctx context.Context, params functions.FunctionParams) (interface{}, *errorhandler.AppError) {
	// Execute custom multiply function
	logger.Debug(fmt.Sprintf(ExecutingMultiplyWithParams, params))
	
	numbers, err := GetNumbersFromParams(params)
	if err != nil {
		return nil, err
	}

	if len(numbers) < 2 {
		return nil, errorhandler.NewAppError(errorhandler.WarningLevel, AtLeastTwoNumbersRequiredErr, nil)
	}

	result := mtpFn.calculateResult(numbers)
	return functions.FunctionResponse{
		Result: result,
		Operation: functions.FunctionOperationText,
		Inputs: numbers,
	}, nil
}

func (mtpFn *FunctionMultiply) GetMetadata() functions.FunctionPayload {
	// Get function metadata
	return functions.FunctionPayload{
		Name:        MultiplyFunctionNameText,
		Description: MultiplyFunctionDescriptionText,
		Parameters: []functions.FunctionParameterMetadata{
			{
				Name:        functions.FunctionNumbersParamText,
				Type:        functions.FunctionArrayTypeText,
				Description: MultiplyArrayDescriptionText,
				Required:    functions.FunctionRequiredValue,
			},
		},
	}
}

func (mtpFn *FunctionMultiply) ConvertToOpenAITool() functions.OpenAIToolsPayload {
	// Convert function to OpenAI tool
	return functions.OpenAIToolsPayload{
		Type: functions.FunctionTypeText,
		Name: MultiplyFunctionNameText,
		Description: MultiplyFunctionDescriptionText,
		Parameters: functions.ToolParametersMetadata{
			Type: functions.FunctionObjectTypeText,
			Properties: map[string]interface{}{
				functions.FunctionNumbersParamText: MultiplyParametersProperties{
					Type:        functions.FunctionArrayTypeText,
					Description: MultiplyArrayDescriptionText,
					Items: MultiplyPropertiesItem{
						Type: functions.FunctionNumberTypeText,
					},
				},
			},
			Required: []string{functions.FunctionNumbersParamText},
		},
	}
}

func (mtpFn *FunctionMultiply) calculateResult(numbers []float64) float64 {
	// Calculate custom multiply result
	result := 1.0
	for _, num := range numbers {
		result *= num
	}
	return result
}
