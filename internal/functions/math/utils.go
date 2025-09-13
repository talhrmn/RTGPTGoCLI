package math

import (
	"RTGPTGoCLI/internal/functions"
	"RTGPTGoCLI/pkg/errorhandler"
	"fmt"
	"strconv"
)

func GetNumbersFromParams(params functions.FunctionParams) ([]float64, *errorhandler.AppError) {
	numbersInterface, exists := params[functions.FunctionNumbersParamText]
	if !exists {
		return nil, errorhandler.NewAppError(
			errorhandler.WarningLevel, 
			fmt.Sprintf(functions.FunctionMissingRequiredParameterErr, 
				functions.FunctionNumbersParamText), 
			nil,
		)
	}

	numbers, ok := numbersInterface.([]interface{})
	if !ok {
		return nil, errorhandler.NewAppError(
			errorhandler.WarningLevel, 
			fmt.Sprintf(functions.FunctionInvalidParameterTypeErr, functions.FunctionNumbersParamText), 
			nil,
		)
	}

	floatNumbers := make([]float64, len(numbers))
	for i, numVal := range numbers {
		switch v := numVal.(type) {
			case string:
				num, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, errorhandler.NewAppError(
						errorhandler.WarningLevel, 
						fmt.Sprintf(functions.InvalidParameterErr, numVal),
						err,
					)
				}
				floatNumbers[i] = num
			default:
				floatNumbers[i] = numVal.(float64)
			}
	}

	return floatNumbers, nil
}
