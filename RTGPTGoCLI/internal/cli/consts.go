package cli

const (
	// Errors
	CLIFailedToWaitUntilReadyErr = "failed to wait until ready: %v"
	CLIFailedToGetInputErr = "failed to get input: %v"
	CLIFailedInputScannerErr = "input scanner error: %v"
)

const (
	// Cli commands
	CLIPromptExit  string = "exit"
	CLIPromptQuit  string = "quit"
	CLIPromptQ     string = "/q"
	CLIPromptClear string = "clear"
	CLIPromptDebug string = "/debug"
	CLIPromptHelp  string = "/help"
	CLIPromptHPrompt     string = "/h"
	CLIPromptFunctionsPrompt string = "/functions"
	CLIPromptFPrompt     string = "/f"
)

const (
	// Log messages
	CLIResponsePanicText = "Response handler panic: %v\n"
)

const (
	// Display messages
	CLIHelpText = 
	`Available commands:
		/help, /h			Show this help message
		/debug				Toggle debug mode
		/functions, /f		Show available functions
		clear				Clear the screen
		exit, quit, /q		Exit the application
	`
	
	CLIWelcomeText = "Chat with Me!"
	CLIDescriptionText = "Ask me anything, am a chatbot base on the gpr-4o-mini-preview model"
	CLIFunctionsText = "Additional current custom functions:\n - multiply"
	CLIGoodbyeText = "Goodbye!"
	CLIPromptText = "> "
	CLIUnknownCommandText = "Unknown command: %s. Type /help or /h for a list of commands.\n"
	CLIUserPrefixText = "You: "
	CLIChatPrefixText = "Chat: "
	CLIAvailableFunctionsText = "Available functions:"
	CLIDebugConfigText = "Debug config: %v"
)

// Signal token for streaming
var StreamSignal = struct{}{}
