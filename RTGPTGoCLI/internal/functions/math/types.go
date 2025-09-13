package math

// Custom multiply function struct
type FunctionMultiply struct{}

type MultiplyPayload struct {
	// Multiply request struct
	Numbers []float64 `json:"numbers"`
}

type MultiplyParametersProperties struct {
	// Multiply parameters properties
	Type string `json:"type"`
	Description string `json:"description"`
	Items MultiplyPropertiesItem `json:"items"`
}

type MultiplyPropertiesItem struct {
	// property item metadata
	Type string `json:"type"`
}	
	