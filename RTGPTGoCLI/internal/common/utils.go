package common

import (
	"encoding/json"
	"fmt"
)

func PrintJSON(data []byte) {
	var message interface{}
	json.Unmarshal(data, &message)
	prettyJSON, _ := json.MarshalIndent(message, "", "  ")
	fmt.Println("Sending message:", string(prettyJSON))
}
