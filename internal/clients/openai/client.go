package openai

import (
	"RTGPTGoCLI/internal/clients"
	"RTGPTGoCLI/internal/common"
	"RTGPTGoCLI/internal/config"
	"RTGPTGoCLI/internal/functions"
	"RTGPTGoCLI/internal/functions/handler"
	"RTGPTGoCLI/pkg/errorhandler"
	"RTGPTGoCLI/pkg/logger"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
)

func NewOAIClient(cfg *config.Config, wsc clients.WebClientConnection) *OpenAIClient {
	// Create new OpenAI client

	functionHandler := handler.NewHandler()
	if err := functionHandler.LoadFunctions(); err != nil {
		logger.Warning(fmt.Sprintf(OAILoadFunctionsErr, err))
	}

	return &OpenAIClient{
		config:           cfg,
		wsc:              wsc,
		mu:               sync.RWMutex{},
		cleanUpOnce:      sync.Once{},

		functionHandler: functionHandler,
		sessionID:        "",
		responseID:       "",
		isStreaming:      false,
		
		messageChannel:   make(chan clients.MessageEvent, cfg.ChannelBuffer),
		errorChannel:     make(chan errorhandler.AppError, cfg.ChannelBuffer),
	}
}

func (oaic *OpenAIClient) Connect(ctx context.Context) error {
	// Connect to OpenAI
	if err := oaic.wsc.Connect(ctx); err != nil {
		return err
	}

	if appErr := oaic.sendSessionConfig(ctx); appErr != nil {
		return appErr.Error
	}
	logger.Debug(OAISessionCreatedMsg)

	go oaic.processMessages(ctx)

	logger.Debug(OAIConnectedMsg)
	return nil
}

func (oaic *OpenAIClient) Disconnect() error {
	// Disconnect from OpenAI
	var disconnectErr error
	oaic.cleanUpOnce.Do(func() {
		logger.Debug(OAIDisconnectingMsg)
		if !oaic.IsConnected() {
			return
		}

		close(oaic.messageChannel)
		close(oaic.errorChannel)

		if err := oaic.wsc.Disconnect(); err != nil {
			disconnectErr = err
		}
	})
	logger.Debug(OAIDisconnectedMsg)
	return disconnectErr
}

func (oaic *OpenAIClient) IsConnected() bool {
	// Return if OpenAI is connected
	return oaic.wsc.IsConnected() && oaic.sessionID != ""
}

func (oaic *OpenAIClient) GetErrorChannel() <-chan errorhandler.AppError {
	// Return error channel
	return oaic.errorChannel
}

func (oaic *OpenAIClient) GetMessageChannel() <-chan clients.MessageEvent {
	// Return message channel
	return oaic.messageChannel
}

func (oaic *OpenAIClient) SendMessage(ctx context.Context, message string) *errorhandler.AppError {
	// Send message to OpenAI
	if oaic.getIsStreaming() {
		return errorhandler.NewAppError(errorhandler.WarningLevel, OAIMessageStreamInProgressMsg, nil)
	}

	conversationItem := OAIConversationPayload{
		Type: OAIConversationItemCreateEventType,
		Item: OAIConversationItemMetadata{
			Type: OAIConversationItemType,
			Role: OAIConversationItemRole,
			Content: []OAIConversationItemContent{
				{
					Type: OAIInputText,
					Text: message,
				},
			},
		},
	}

	if appErr := oaic.sendToWebSocket(ctx, conversationItem); appErr != nil {
		return appErr
	}

	messagePayload := OAIResponsePayload{
		Type: OAIResponseCreateEventType,
		Response: OAIResponseMetadata{
			Instructions: message,
		},
	}

	if appErr := oaic.sendToWebSocket(ctx, messagePayload); appErr != nil {
		return appErr
	}

	return nil
}

func (oaic *OpenAIClient) GetAvailableFunctions() []string {
	// Return available custom functions
	tools := oaic.functionHandler.GenerateOpenAITools()
	names := make([]string, len(tools))
	for i, tool := range tools {
		names[i] = tool.(map[string]interface{})[OAIFunctionFieldName].(string)
	}
	return names
}

func (oaic *OpenAIClient) sendToWebSocket(ctx context.Context, payload interface{}) *errorhandler.AppError {
	// Send payload to WebSocket
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return common.NewErrJsonMarshalAppError(err)
	}

	if err := oaic.wsc.SendMessage(ctx, payloadBytes); err != nil {
		return errorhandler.NewAppError(errorhandler.WarningLevel, OAISendMessageErr, err)
	}
	return nil
}

func (oaic *OpenAIClient) sendSessionConfig(ctx context.Context) *errorhandler.AppError {
	// Define and send session config
	tools := oaic.functionHandler.GenerateOpenAITools()
	sessionConfigPayload := OAISessionConfigPayload{
		Type: OAISessionUpdateEventType,
		Session: OAISessionConfigMetadata{
			Type: OAISessionTypeRealtimeText,
			OutputModalities:    []string{OAISessionModalitiesText}, 
			Instructions:  OAISessionInstructionsText,
			Tools:         tools,
			ToolChoice:   OAISessionToolsChoiceText,
		},
	}

	if appErr := oaic.sendToWebSocket(ctx, sessionConfigPayload); appErr != nil {
		return appErr
	}

	return nil
}

func (oaic *OpenAIClient) getIsStreaming() bool {
	// Return if OpenAI is streaming
	oaic.mu.RLock()
	defer oaic.mu.RUnlock()
	return oaic.isStreaming
}

func (oaic *OpenAIClient) setIsStreaming(status bool) {
	// Set if OpenAI is streaming
	oaic.mu.Lock()
	defer oaic.mu.Unlock()
	oaic.isStreaming = status
}

func (oaic *OpenAIClient) processMessages(ctx context.Context) {
	// Process messages from WebSocket
	for {
		select {
		case <-ctx.Done():
			return
		case err := <-oaic.wsc.GetErrorChannel():
			oaic.errorChannel <- err
		default:
			msg, ok := <-oaic.wsc.GetMessageChannel()
			if !ok {
				switch {
				case ctx.Done() != nil:
					return
				}
				oaic.errorChannel <- *errorhandler.NewAppError(errorhandler.ErrorLevel, OAIDisconnectedMsg, nil)
				return
			}
			oaic.handleEvent(ctx, msg)
		}
	}
}

func (oaic *OpenAIClient) handleEvent(ctx context.Context, event []byte) {
	// Handle event from OpenAI
	var messageType OAIStreamingEvent
	if err := json.Unmarshal(event, &messageType); err != nil {
		oaic.errorChannel <- *common.NewErrJsonUnmarshalAppError(err)
		return
	}
	
	msgType := messageType.Type
	switch msgType {
	case OAISessionCreatedEventType:
		oaic.handleSessionCreated(event)
	case OAISessionUpdateEventType:
		oaic.handleSessionUpdate(event)
	case OAIResponseCreatedEventType:
		oaic.handleResponseCreated(event)
	case OAIResponseDeltaEventType:
		oaic.handleResponseDelta(event)
	case OAIResponseDeltaDoneEventType, OAIResponseDoneEventType, OAIResponseOutputItemDoneEventType, OAIConversationItemDoneEventType, OAIFunctionCallDoneEventType:
		oaic.handleResponseDone(ctx, msgType, event)
	case OAIResponseFailedEventType, OAIResponseErrorEventType:
		oaic.handleResponseError(event)
	}
}

func (oaic *OpenAIClient) handleSessionCreated(msg []byte) {
	// Handle session created event
	var created OAISessionCreatedEventPayload
	if err := json.Unmarshal(msg, &created); err != nil {
		oaic.errorChannel <- *common.NewErrJsonUnmarshalAppError(err)
		return
	}
	oaic.sessionID = created.Session.ID
	logger.Debug(fmt.Sprintf(OAISessionCreatedWithIDMsg, oaic.sessionID))
}

func (oaic *OpenAIClient) handleSessionUpdate(msg []byte) {
	// Handle session update event
	logger.Debug(OAISessionUpdatedMsg)
}

func (oaic *OpenAIClient) handleResponseCreated(msg []byte) {
	// Handle response created event
	var created OAIResponseCreatedEventPayload
	if err := json.Unmarshal(msg, &created); err != nil {
		oaic.errorChannel <- *common.NewErrJsonUnmarshalAppError(err)
		return
	}
	oaic.responseID = created.Response.ID
	logger.Debug(fmt.Sprintf(OAIResponseCreatedWithIDMsg, oaic.responseID))
}

func (oaic *OpenAIClient) handleResponseDelta(msg []byte) {
	// Handle response delta event, returning deltas to simulate chat streaming
	if !oaic.getIsStreaming() {
		oaic.setIsStreaming(true)
	}
	
	var delta OAIResponseOutPutTextDeltaPayload
	if err := json.Unmarshal(msg, &delta); err != nil {
		oaic.errorChannel <- *common.NewErrJsonUnmarshalAppError(err)
		return
	}
	oaic.messageChannel <- clients.MessageEvent{Type: OAIResponseDeltaEventType, Text: delta.Delta, Done: false}
}

func (oaic *OpenAIClient) handleResponseDone(ctx context.Context, msgType string, msg []byte) {
	// Handle response done event
	switch msgType {
	case OAIResponseDeltaDoneEventType:
		oaic.messageChannel <- clients.MessageEvent{Type: OAIResponseDeltaDoneEventType, Text: "", Done: true}
	case OAIFunctionCallDoneEventType:
		oaic.handleFunctionCallDone(ctx, msg)
	}
	oaic.setIsStreaming(false)
}

func (oaic *OpenAIClient) handleFunctionCallDone(ctx context.Context, msg []byte) {
	// Handle function call done event
	var functionCallDone OAIFunctionCallDonePayload

	if err := json.Unmarshal(msg, &functionCallDone); err != nil {
		oaic.errorChannel <- *common.NewErrJsonUnmarshalAppError(err)
		return
	}

	logger.Debug(fmt.Sprintf(OAIExecutingFunctionWithArgsMsg, functionCallDone.Name, functionCallDone.Arguments))

	// oaic.messageChannel <- clients.MessageEvent{
	// 	Type: OAIFunctionCallDeltaEventType,
	// 	Text: fmt.Sprintf(OAIFunctionCallDeltaMsg, functionCallDone.Name, functionCallDone.Arguments),
	// 	Done: false,
	// }

	result, appErr := oaic.functionHandler.Execute(ctx, functionCallDone.Name, functionCallDone.Arguments)
	if appErr != nil {
		oaic.errorChannel <- *appErr
		return
	}
	
	resultMap, ok := result.(functions.FunctionResponse)
	if !ok {
		oaic.errorChannel <- *errorhandler.NewAppError(errorhandler.WarningLevel, fmt.Sprintf(OAIUnexpectedFunctionResultType, result), nil)
		return
	}

	resultToSend := fmt.Sprintf("%v", resultMap.Result)
	oaic.sendFunctionResult(ctx, functionCallDone, resultToSend)
}

func (oaic *OpenAIClient) sendFunctionResult(ctx context.Context, functionData OAIFunctionCallDonePayload, result interface{}) {
	// Send function result to OpenAI
	functionResultPayload := OAIFunctionCallResultPayload{
		Type: OAIConversationItemCreateEventType,
		Item: OAIFunctionResultItemMetadata{
			Type:   OAIFunctionCallResultText,
			CallID: functionData.CallID,
			Output: result,
		},
	}

	if appErr := oaic.sendToWebSocket(ctx, functionResultPayload); appErr != nil {
		oaic.errorChannel <- *appErr
		return
	}

	continueFunctionCallPayload := OAIFunctionResultContinuePayload{
		Type: OAIResponseCreateEventType,
		Response: FunctionContinueResponseMetadata{
			Instructions: fmt.Sprintf(OAIFunctionCallInstructions, functionData.Name, functionData.Arguments, result),
		},
	}

	if appErr := oaic.sendToWebSocket(ctx, continueFunctionCallPayload); appErr != nil {
		oaic.errorChannel <- *appErr
		return
	}
}

func (oaic *OpenAIClient) handleResponseError(event []byte) {
	// Handle response error event by type of error
	oaic.setIsStreaming(false)
	var errorType OAIStreamingEvent
	if err := json.Unmarshal(event, &errorType); err != nil {
		oaic.errorChannel <- *common.NewErrJsonUnmarshalAppError(err)
		return
	}

	switch errorType.Type {
	case OAIResponseFailedEventType:
		var messageFailure OAIResponseFailedEventPayload
		if err := json.Unmarshal(event, &messageFailure); err != nil {
			oaic.errorChannel <- *common.NewErrJsonUnmarshalAppError(err)
			return
		}
		errorMsg := fmt.Sprintf(OAIFailedResponseErr, messageFailure.Response.Error.Code, messageFailure.Response.Error.Message)
		oaic.errorChannel <- *errorhandler.NewAppError(errorhandler.ErrorLevel, errorMsg, errors.New(OAIFailedResponseErr))
	case OAIResponseErrorEventType:
		var messageError OAIResponseErrorPayload
		if err := json.Unmarshal(event, &messageError); err != nil {
			oaic.errorChannel <- *common.NewErrJsonUnmarshalAppError(err)
			return
		}
		errorMsg := fmt.Sprintf(OAIFailedResponseErr, messageError.Error.Code, messageError.Error.Message)
		oaic.errorChannel <- *errorhandler.NewAppError(errorhandler.ErrorLevel, errorMsg, errors.New(OAIErrorResponseErr))
	}
}
