package ui

import "time"

const (
	// ui colors
	UIBlueColor = "\033[34m"
	UIResetColor = "\033[0m"
	UIGreenColor = "\033[32m"
	UIRedColor = "\033[31m"
	UICyanColor = "\033[36m"
	UIYellowColor = "\033[33m"
)

const (
	// ui os commands
	UIClearLineCommand = "\r\033[K"
	UIClearScreenCommand = "\033[H\033[2J"
)

const (
	// ui constants
	UIProcessSleepTime = 400 * time.Millisecond
	UIProcessingText = "\r%sProcessing%s%s"
	UIErrorPrefix = "Error: "
)

// Dots for streaming
var UIProcessingDots = [...]string{".  ", ".. ", "..."}
