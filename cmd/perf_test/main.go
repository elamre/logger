package main

import (
	"fmt"
	"github.com/elamre/logger/internal"
	"github.com/elamre/logger/pkg/logger"
	"log"
	"sync"
	"time"
)

const serverPort = 12321
const resultCheck = true

func TestServerSpeed(ansii bool, loops int, testMessage string) time.Duration {
	logger.GetSettings().Reset()
	logger.GetSettings().LogAll(true)
	logger.GetSettings().FormatAnsiiColours(ansii)
	internal.CheckError(logger.GetSettings().EnableUdpBroadcast(serverPort))
	wait := new(sync.Mutex)
	go func() {
		wait.Lock()
		defer wait.Unlock()
		internal.CheckError(logger.GetSettings().WaitForUdpBroadcast())
	}()
	waitFunc := startTerminal(serverPort, loops, "", testMessage)
	wait.Lock()

	testLogger := logger.NewLogger()
	count := 0

	start := time.Now()
	for count < loops {
		testLogger.LogDebugf("%s %d", testMessage, count)
		count++
	}
	dur := time.Now().Sub(start)
	waitFunc()
	logger.GetSettings().DisableUdpBroadcast()
	return dur
}

func TestLocalSpeed(ansii bool, loops int, testMessage string) time.Duration {
	logger.GetSettings().Reset()
	logger.GetSettings().LogAll(true)
	logger.GetSettings().FormatAnsiiColours(ansii)

	testLogger := logger.NewLogger()
	count := 0

	start := time.Now()
	for count < loops {
		testLogger.LogDebugf("%s %d", testMessage, count)
		count++
	}
	dur := time.Now().Sub(start)
	return dur
}

func TestToFile(ansii bool, loops int, testMessage string) time.Duration {
	logger.GetSettings().Reset()
	logger.GetSettings().LogAll(true)
	logger.GetSettings().FormatAnsiiColours(ansii)
	logger.GetSettings().SetHook(logger.GetDefaultFileHook(fmt.Sprintf("output%v.tmp", ansii)))
	testLogger := logger.NewLogger()
	count := 0

	start := time.Now()
	for count < loops {
		testLogger.LogDebugf("%s %d", testMessage, count)
		count++
	}
	dur := time.Now().Sub(start)
	return dur
}

func main() {
	resultsMap := make(map[string]time.Duration)
	resultsMap["server_ansii"] = TestServerSpeed(true, 10000, "ansii server")
	resultsMap["server"] = TestServerSpeed(false, 10000, "ansii server")
	resultsMap["local_ansii"] = TestLocalSpeed(true, 10000, "ansii server")
	resultsMap["local"] = TestLocalSpeed(false, 10000, "ansii server")
	resultsMap["file_ansii"] = TestToFile(true, 10000, "ansii server")
	resultsMap["file"] = TestToFile(false, 10000, "ansii server")
	log.Printf("%+v", resultsMap)
}
