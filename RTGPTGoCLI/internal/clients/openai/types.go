package openai

import (
	"RTGPTGoCLI/RTGPTGoCLI/internal/clients"
	"RTGPTGoCLI/RTGPTGoCLI/internal/config"
	"RTGPTGoCLI/RTGPTGoCLI/pkg/errorhandler"
	"sync"
)

type OpenAIClientInterface interface {
	// OpenAIClient interface
	clients.ServiceClientConnection
	GetAvailableFunctions() []string
}

type OpenAIClient struct {
	// OpenAIClient struct
	config *config.Config
	wsc    clients.WebClientConnection
	mu     sync.RWMutex
	cleanUpOnce sync.Once

	sessionID   string
	responseID  string
	isStreaming bool

	messageChannel chan clients.MessageEvent
	errorChannel   chan errorhandler.AppError	
}

type OAISessionConfigPayload struct {
	// OpenAI session config struct
	Type    string             `json:"type"`
	Session OAISessionConfigMetadata `json:"session"`
}

type OAISessionConfigMetadata struct {
	// OpenAI session config metadata struct
	Type             string   `json:"type,omitempty"`
	OutputModalities []string `json:"output_modalities,omitempty"`
	Instructions     string   `json:"instructions,omitempty"`
	Tools	            []interface{} `json:"tools,omitempty"`
	ToolChoice       string   `json:"tool_choice,omitempty"`
}

type OAIConversationPayload struct {
	// OpenAI conversation payload struct
	Type string `json:"type"`
	Item OAIConversationItemMetadata `json:"item"`
}

type OAIConversationItemMetadata struct {
	// OpenAI conversation item metadata struct
	Type string `json:"type"`
	Role string `json:"role"`
	Content []OAIConversationItemContent `json:"content"`
}

type OAIConversationItemContent struct {
	// OpenAI conversation item content struct
	Type string `json:"type"`
	Text string `json:"text"`
}

type OAIResponsePayload struct {
	// OpenAI response payload struct
	Type     string          `json:"type"`
	Response OAIResponseMetadata `json:"response"`
}

type OAIResponseMetadata struct {
	// OpenAI response metadata struct
	Conversation string `json:"conversation,omitempty"`
	Instructions string `json:"instructions,omitempty"`
}

type OAIStreamingEvent struct {
	// OpenAI streaming event struct
	Type string `json:"type"`
}

type OAISessionCreatedEventPayload struct {
	// OpenAI session created event payload struct
	Type    string          `json:"type"`
	Session OAISessionCreatedMetadata `json:"session"`
}

type OAISessionCreatedMetadata struct {
	// OpenAI session created metadata struct
	ID string `json:"id"`
}

type OAIResponseCreatedEventPayload struct {
	// OpenAI response created event payload struct
	Type     string           `json:"type"`
	Response OAIResponseCreatedMetadata `json:"response"`
}

type OAIResponseCreatedMetadata struct {
	// OpenAI response created metadata struct
	ID     string `json:"id"`
	Status string `json:"status"`
}

type OAIResponseOutPutTextDeltaPayload struct {
	// OpenAI response output text delta payload struct
	Type      string `json:"type"`
	ItemId    string `json:"item_id"`
	Delta     string `json:"delta"`
	SeqNumber int    `json:"sequence_number"`
}

type OAIResponseFailedEventPayload struct {
	// OpenAI response failed event payload struct
	Type     string `json:"type"`
	Response OAIResponseFailedMetadata `json:"response"`
}

type OAIResponseFailedMetadata struct {
	// OpenAI response failed metadata struct
	Error OAIResponseFailedErrorData `json:"error"`
}

type OAIResponseFailedErrorData struct {
	// OpenAI response failed error data struct
	Code    string `json:"code"`
	Message string `json:"message"`
}

type OAIResponseErrorPayload struct {
	// OpenAI response error payload struct
	Type    string      `json:"type"`
	EventID string      `json:"event_id"`
	Error   OAIErrorMetadata `json:"error"`
}

type OAIErrorMetadata struct {
	// OpenAI error metadata struct
	Code    string `json:"code"`
	Message string `json:"message"`
	Param   string `json:"param,omitempty"`
}
