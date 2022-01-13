package main

import (
	"github.com/elamre/logger/pkg/logger"
)

var mainLogger = logger.NewLogger()

func main() {
	logger.GetSettings().SetPanicOnError(false)
	logger.GetSettings().LogAll(true)
	logger.GetSettings().FormatAnsiiColours(true)
	count := 0

	for count < 1000000 {
		mainLogger.LogDebugf("Debug %d", count)
		count++
	}
}
