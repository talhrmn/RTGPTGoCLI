package functions

import (
	"RTGPTGoCLI/RTGPTGoCLI/pkg/errorhandler"
	"context"
)

// Function parameters type
type FunctionParams map[string]interface{}

type FunctionInterface interface {
	// Custom function execution
	Execute(ctx context.Context, params FunctionParams) (interface{}, *errorhandler.AppError)
	GetMetadata() FunctionPayload
	ConvertToOpenAITool() OpenAIToolsPayload
}

type OpenAIToolsPayload struct {
	// OpenAI tool metadata
	Type string `json:"type"`
	Name string `json:"name"`
	Description string `json:"description"`
	Parameters ToolParametersMetadata `json:"parameters"`
}

type ToolParametersMetadata struct {
	// Tool parameters metadata
	Type string `json:"type"`
	Properties map[string]interface{} `json:"properties"`
	Required []string `json:"required,omitempty"`
}


type FunctionPayload struct {
	// Function metadata
	Name        string
	Description string
	Parameters  []FunctionParameterMetadata
}

type FunctionParameterMetadata struct {
	// Function parameter metadata
	Name        string
	Type        string
	Description string
	Required    bool
}

type FunctionResponse struct {
	Result  float64 `json:"result"`
	Operation string `json:"operation"`
	Inputs    []float64 `json:"inputs"`
}
