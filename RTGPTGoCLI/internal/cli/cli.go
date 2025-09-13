package cli

import (
	"RTGPTGoCLI/RTGPTGoCLI/internal/cli/ui"
	"RTGPTGoCLI/RTGPTGoCLI/internal/clients/openai"
	"RTGPTGoCLI/RTGPTGoCLI/internal/config"
	"RTGPTGoCLI/RTGPTGoCLI/pkg/errorhandler"
	"RTGPTGoCLI/RTGPTGoCLI/pkg/logger"
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"
)

func New(cfg *config.Config, oaiClient openai.OpenAIClientInterface) *CLI {
	// Create new CLI
	return &CLI{
		config:    cfg,
		scanner:   bufio.NewScanner(os.Stdin),
		oaiClient: oaiClient,
		streamingChannel: make(chan struct{}, cfg.ChannelBuffer),
		errorHandler: errorhandler.NewErrorHandler(cfg.Debug),
	}
}

func (cli *CLI) Run(ctx context.Context, cancel context.CancelFunc) {
	// Run CLI
	go cli.handleChatOutput(ctx)

	if err := cli.waitUntilReady(ctx); err != nil {
		appErr := *errorhandler.NewAppError(errorhandler.ErrorLevel, fmt.Sprintf(CLIFailedToWaitUntilReadyErr, err), err)
		cli.errorHandler.HandleError(appErr)
		return
	}

	welcomeAdditionalText := CLIDescriptionText + "\n" + CLIHelpText + "\n" + CLIFunctionsText
	ui.ShowWelcome(CLIWelcomeText, welcomeAdditionalText)

	for {
		select {
		case <-ctx.Done():
			ui.ShowGoodbye(CLIGoodbyeText)
			return
		default:
			inputPrompt, appErr := cli.getInput(ctx)
			if appErr != nil {
				cli.errorHandler.HandleError(*appErr)
				return
			}

			if inputPrompt == "" {
				continue
			}
			cli.processInput(ctx, cancel, inputPrompt)
		}

	}
}

func (cli *CLI) waitUntilReady(ctx context.Context) error {
	// Wait until ready
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if cli.oaiClient.IsConnected() {
				return nil
			}
			time.Sleep(time.Duration(cli.config.Timeout) * time.Millisecond)
		}
	}
}

func (cli *CLI) getInput(ctx context.Context) (string, *errorhandler.AppError) {
	// Get input
	ui.ShowPrompt(CLIPromptText)

	inputChannel := make(chan string, cli.config.ChannelBuffer)
	errorChannel := make(chan errorhandler.AppError, cli.config.ChannelBuffer)

	go func() {
		if cli.scanner.Scan() {
			inputChannel <- cli.scanner.Text()
		} else {
			if err := cli.scanner.Err(); err != nil {
				errorChannel <- *errorhandler.NewAppError(errorhandler.ErrorLevel, fmt.Sprintf(CLIFailedInputScannerErr, err), err)
			} else {
				inputChannel <- CLIPromptExit
			}
		}
	}()

	select {
	case <-ctx.Done():
		return CLIPromptExit, nil
	case err := <-errorChannel:
		return "", &err
	case inputPrompt := <-inputChannel:
		return strings.TrimSpace(inputPrompt), nil
	}
}

func (cli *CLI) processInput(ctx context.Context, cancel context.CancelFunc, input string) {
	// Process input
	switch {
	case input == "":
	case input == CLIPromptExit || input == CLIPromptQuit || input == CLIPromptQ:
		cancel()
	case input == CLIPromptHelp || input == CLIPromptHPrompt:
		ui.Show("", CLIHelpText)
	case input == CLIPromptClear:
		ui.Clear()
	case input == CLIPromptDebug:
		cfgString, err := cli.config.GetConfigInfo()
		if err != nil {
			ui.ShowError(err)
		}
		ui.Show(CLIDebugConfigText, cfgString)
	case input == CLIPromptFunctionsPrompt || input == CLIPromptFPrompt:
		ui.ShowFunctions(CLIAvailableFunctionsText, cli.oaiClient.GetAvailableFunctions())
	case strings.HasPrefix(input, "/"):
		ui.ShowError(fmt.Errorf(CLIUnknownCommandText, input))
	default:
		cli.handleChatInput(ctx, input)
	}
}


func (cli *CLI) handleChatInput(ctx context.Context, prompt string) {
	// Handle chat input
	ui.ShowUserMessage(CLIUserPrefixText, prompt)

	go ui.ShowChatProcessing(cli.streamingChannel)

	if appErr := cli.oaiClient.SendMessage(ctx, prompt); appErr != nil {
		cli.errorHandler.HandleError(*appErr)
	}
}

func (cli *CLI) handleChatOutput(ctx context.Context) {
	// Handle chat output
	defer func() {
		if r := recover(); r != nil {
			logger.Warning(fmt.Sprintf(CLIResponsePanicText, r))
		}
	}()

	isFirstDelta := true
	for {
		select {
		case <-ctx.Done():
			return
		case appErr := <-cli.oaiClient.GetErrorChannel():
			cli.errorHandler.HandleError(appErr)
		case msg := <-cli.oaiClient.GetMessageChannel():
			if msg.Text == "" && !msg.Done {
				continue
			}

			if isFirstDelta {
				cli.streamingChannel <- StreamSignal
				ui.ShowChatPrefix(CLIChatPrefixText)
				isFirstDelta = false
			}

			if !msg.Done {
				ui.ShowChatDelta(msg.Text)
			} else {
				ui.EndStreaming()
				ui.ShowPrompt(CLIPromptText)
				isFirstDelta = true
			}
		}
	}
}
