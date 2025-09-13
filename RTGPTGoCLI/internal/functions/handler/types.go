package handler

import "RTGPTGoCLI/RTGPTGoCLI/internal/functions"

// Functions type
type FunctionsType map[string]functions.FunctionInterface

type FunctionHandler struct {
	// Function handler struct
	functions FunctionsType
}
