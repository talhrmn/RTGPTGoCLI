package ui

import (
	"RTGPTGoCLI/internal/config"
	"fmt"
	"strings"
	"time"
)

func stringRepeat(str string, count int) string {
	// repeat string helper
	return strings.Repeat(str, count)
}

func ShowWelcome(welcomeText string, additionalText string) {
	// show welcome message
	line := "╔" + stringRepeat("═", len(welcomeText) + 4) + "╗"
	fmt.Println(UICyanColor + line + UIResetColor)
	fmt.Println(UICyanColor + "║" + UIYellowColor + "  " + welcomeText + "  " + UICyanColor + "║" + UIResetColor)
	fmt.Println(UICyanColor + "╚" + stringRepeat("═", len(welcomeText) + 4) + "╝" + UIResetColor)
	fmt.Println()
	fmt.Println(additionalText)
	fmt.Println()
}

func ShowGoodbye(goodbyeText string) {
	// show goodbye message
	fmt.Println()
	fmt.Println(goodbyeText)
}

func ShowPrompt(promptText string) {
	// show input prompt
	fmt.Print(UIBlueColor + promptText + UIResetColor)
}

func Clear() {
	// clear terminal
	fmt.Print(UIClearScreenCommand)
}

func ShowError(err error) {
	// show error message
	fmt.Println(UIRedColor + UIErrorPrefix + UIResetColor + err.Error())
}

func Show(prefix string, text string) {
	// show message
	fmt.Println(prefix + text + "\n")
}

func ShowDebug(cfg *config.Config) {
	// show debug config
	cfgString, err := cfg.GetConfigInfo()
	if err != nil {
		ShowError(err)
		return
	}
	fmt.Println(cfgString)
}

func EndStreaming() {
	// end streaming message
	fmt.Println()
}

func ShowUserMessage(prefix string, message string) {
	// show user message with a user prefix
	fmt.Printf(UIBlueColor + prefix + UIResetColor + "%s\n", message)
}

func ShowChatPrefix(prefix string) {
	// show chatbot message prefix
	ClearLine()
	fmt.Print(UIGreenColor + prefix + UIResetColor)
}

func ShowChatDelta(delta string) {
	// show chatbot message delta
	fmt.Print(delta)
}

func ClearLine() {
	// Clear the line
	fmt.Print(UIClearLineCommand)
}

func ShowChatProcessing(streamingChannel <-chan struct{}) {
	// show chatbot message processing with changing dots
	i := 0
	for {
		select {
		case <-streamingChannel:
			return
		default:
			fmt.Printf(UIProcessingText, UIGreenColor, UIProcessingDots[i%len(UIProcessingDots)], UIResetColor)
			i++
			time.Sleep(UIProcessSleepTime)
		}
	}
}

func ShowFunctions(prefix string, functions []string) {
	// show available custom functions
	fmt.Println(prefix)
	for _, fn := range functions {
		fmt.Println("- " + fn)
	}
	fmt.Println()
}
