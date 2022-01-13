package main

import (
	"github.com/elamre/logger/pkg/logger"
	"log"
)

var mainLogger = logger.NewLogger()

func errorAndRecover() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("recovered from: %s", r)
		}
		log.Println("What")
	}()
	mainLogger.LogError("error panics!")
}

func main() {
	mainLogger.LogDebug("This will not be printed")
	mainLogger.PrintDebug(true)
	mainLogger.LogDebug("This will be printed")
	mainLogger.LogInfo("Log with filename info!")
	mainLogger.LogInfof("printf %d %02x %s", 123, 0xaabb, "Hallo")

	mainLogger.LogWarning("ANSII colours")
	//errorAndRecover()
	logger.GetSettings().SetPanicOnError(false)
	mainLogger.LogError("Now no panic on error")
	logger.GetSettings().FormatAnsiiColours(false)
	mainLogger.LogInfo("No")
	mainLogger.LogDebug("Colours")
	mainLogger.LogWarning("Enabled")
	mainLogger.LogError("Now")

	logger.GetSettings().FormatAnsiiColours(true)

	logger.GetSettings().LogAll(true)
	newLogger := logger.NewNamedLogger("tester")
	newLogger.LogDebug("Logging all enabled, debug!")
	logger.GetSettings().GetLogger("tester").PrintDebug(false)
	newLogger.LogDebug("No printing...")

	logger.GetSettings().SetHook(logger.GetDefaultFileHook("output.txt"))
	newLogger.LogInfo("Now we continue in a text file")
	mainLogger.LogInfo("I also still work fine")
}
