package common

import (
	"encoding/json"
	"fmt"
)

func PrintJSON(data []byte) {
	// Print JSON data for debugging
	var message interface{}
	json.Unmarshal(data, &message)
	prettyJSON, _ := json.MarshalIndent(message, "", "  ")
	fmt.Println("JSON: ", string(prettyJSON))
}
