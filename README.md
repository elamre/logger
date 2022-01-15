# Simple logger
A simple logger that can log to network, file and terminal output. Ansii colouring enabled. The loggers are all connected to a global settings instance. Enabling debug print of individual loggers. 

## Usage

Loggers can be made within a function or static. Either assigned a name automatically based on the filename, or a given name.

```golang
package main

import "github.com/elamre/logger/pkg/logger"

var fileNameLogger = logger.NewLogger()
var namedLogger = logger.NewNamedLogger("Named")

func main() {
	internalLogger := logger.NewLogger()
	internalLogger.PrintDebug(true)
	// Now we can print to debug
	internalLogger.LogDebugf("Debug %d", 213)
	internalLogger.LogInfo("Info")
	internalLogger.LogWarning("Warning")
	fileNameLogger.LogInfo("Different kind of logger")
	// Let no loggers write debug statements
	logger.GetSettings().SetSuppressDebug(true)
	namedLogger.LogInfo("now writing something to debug")
	namedLogger.LogDebug("This wont be printed")
	// Now lets write to a file instead
	logger.GetSettings().SetHook(logger.GetDefaultFileHook("output.txt"))
	namedLogger.LogInfo("Now we continue in a text file")
	internalLogger.LogInfo("I also still work fine")
}
```
The example above shows how to use most of the functionality, it also works when used across different files. It does not work with multiple threads.
For performance usages the [perf_test](cmd/perf_test/main.go) can be examined. 

Network functionality makes use of Gob, it's far from the most efficient, but it was easy to implement, and it's good enough for my goals.
These examples can be found in [server](cmd/logger_server_test/main.go) and [terminal](cmd/logger_terminal/main.go).

## Performance
Measured on my laptop, writing a total of 40,000 lines, all in milliseconds:

11th Gen Intel(R) Core(TM) i7-1165G7

|          | Without | With Ansii |
|----------|---------|------------|
| File     | 41.291  | 43.636     |
| Terminal | 584.65  | 632.114    |
| Server   | 128.951 | 144.776    |

