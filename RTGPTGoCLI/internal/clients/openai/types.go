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
	Type string `json:"type"`
	Item OAIConversationItemMetadata `json:"item"`
}

type OAIConversationItemMetadata struct {
	Type string `json:"type"`
	Role string `json:"role"`
	Content []OAIConversationItemContent `json:"content"`
}

type OAIConversationItemContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type OAIResponsePayload struct {
	Type     string          `json:"type"`
	Response OAIResponseMetadata `json:"response"`
}

type OAIResponseMetadata struct {
	Conversation string `json:"conversation,omitempty"`
	Instructions string `json:"instructions,omitempty"`
}

type OAIStreamingEvent struct {
	Type string `json:"type"`
}

type OAISessionCreatedEventPayload struct {
	Type    string          `json:"type"`
	Session OAISessionCreatedMetadata `json:"session"`
}

type OAISessionCreatedMetadata struct {
	ID string `json:"id"`
}

type OAIResponseCreatedEventPayload struct {
	Type     string           `json:"type"`
	Response OAIResponseCreatedMetadata `json:"response"`
}

type OAIResponseCreatedMetadata struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type OAIResponseOutPutTextDeltaPayload struct {
	Type      string `json:"type"`
	ItemId    string `json:"item_id"`
	Delta     string `json:"delta"`
	SeqNumber int    `json:"sequence_number"`
}

type OAIResponseFailedEventPayload struct {
	Type     string `json:"type"`
	Response OAIResponseFailedMetadata `json:"response"`
}

type OAIResponseFailedMetadata struct {
	Error OAIResponseFailedErrorData `json:"error"`
}

type OAIResponseFailedErrorData struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type OAIResponseErrorPayload struct {
	Type    string      `json:"type"`
	EventID string      `json:"event_id"`
	Error   OAIErrorMetadata `json:"error"`
}

type OAIErrorMetadata struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Param   string `json:"param,omitempty"`
}
