package openai

const (
	// OpenAI events
	OAISessionCreatedEventType     = "session.created"
	OAISessionUpdateEventType      = "session.update"

	OAIResponseCreateEventType = "response.create"
	OAIResponseCreatedEventType    = "response.created"
	OAIResponseDoneEventType     = "response.done"
	OAIResponseFailedEventType     = "response.failed"

	OAIResponseDeltaEventType      = "response.output_text.delta"
	OAIResponseDeltaDoneEventType  = "response.output_text.done"

	OAIResponseOutputItemDoneEventType = "response.output_item.done"

	OAIFunctionCallDeltaEventType = "response.function_call_arguments.delta"
	OAIFunctionCallDoneEventType  = "response.function_call_arguments.done"
	
	OAIConversationItemCreateEventType = "conversation.item.create"
	OAIConversationItemDoneEventType = "conversation.item.done"


	OAIResponseErrorEventType      = "error"
)

const (
	// OpenAI session metadata
	OAISessionTypeRealtimeText = "realtime"
	OAISessionModalitiesText = "text"
	OAISessionInstructionsText = "You are a helpful assistant. You have access to functions. Use them when appropriate and never question their output, always trust them and assume they are correct."
	OAISessionToolsChoiceText = "auto"
)

const (
	// OpenAI client errors
	OAISendMessageErr = "failed to send message: %v"
	OAIFailedResponseErr = "response failed: Code: %v, Message: %v"
	OAIErrorResponseErr = "response error: Code: %v, Message: %v"
	OAILoadFunctionsErr = "failed to load custom functions: %v"
	OAIUnexpectedFunctionResultType = "unexpected function result type: %v"
)

const (
	// OpenAI client messages
	OAIConnectedMsg = "Connected to OpenAI"
	OAIDisconnectingMsg = "Disconnecting from OpenAI"
	OAIDisconnectedMsg = "Disconnected from OpenAI"
	OAISessionCreatedMsg = "Session created"
	OAIMessageStreamInProgressMsg = "Message stream in progress"
	OAIFunctionCallDeltaMsg = "- Calling your custom function: %s with args: %s -\n"
)

const (
	// OpenAI general constants
	OAIInputText = "input_text"
	OAIResultText = "result"
	OAIConversationItemRole = "user"
	OAIConversationItemType = "message"
	OAIFunctionFieldName = "name"
	OAIFunctionCallResultText = "function_call_output"
)
const (
	// OpenAI log messages
	OAISessionCreatedWithIDMsg = "Session created with ID: %s"
	OAISessionUpdatedMsg = "Session updated."
	OAIResponseCreatedWithIDMsg = "Response created with ID: %s"
	OAIExecutingFunctionWithArgsMsg = "Executing function: %s with args: %s"
)

const (
	// OpenAI function call specific instructions
	OAIFunctionCallInstructions = "A function named '%s' was called with arguments %s. The function returned: %v. " +
		"Write a clear reply in natural language, make sure you mention that you have called a custom function " +
		"and that writes the original problem and incorporates the function result directly as if it was your own. " +
		"Never recompute or override the function output, always treat it as ground truth."
)
