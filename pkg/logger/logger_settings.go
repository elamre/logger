package logger

import (
	"bufio"
	"fmt"
	"github.com/elamre/logger/internal"
	"log"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

var onceSettings sync.Once
var settings LoggerSettings

var utilLog internal.LoggerInterface
var serverLog internal.LoggerInterface

type HookFunc func(time time.Time, loggerName string, fileName string, format string)

type LoggerSettings struct {
	file           *os.File
	writer         *bufio.Writer
	panicOnError   bool
	warningAsError bool
	suppressDebug  bool
	defaultLog     bool
	showColor      bool
	loggers        map[string]*Logger
	loggerHook     HookFunc

	server internal.LoggerServer
}

func newLoggerSettings() LoggerSettings {
	setting := LoggerSettings{
		file:          nil,
		panicOnError:  true,
		suppressDebug: false,
		showColor:     true,
		writer:        bufio.NewWriter(os.Stderr),
		loggers:       make(map[string]*Logger),
	}
	return setting
}

func (settings *LoggerSettings) LogAll(log bool) {
	for _, l := range settings.loggers {
		l.debug = log
	}
	settings.defaultLog = log
}

func (settings *LoggerSettings) SetPanicOnError(newPanicOnError bool) {
	settings.panicOnError = newPanicOnError
	if !settings.panicOnError {
		utilLog.LogWarning("No panic on error. This might give unwanted results")
	}
}

func (settings *LoggerSettings) SetWarningAsError(newWarningAsError bool) {
	settings.warningAsError = newWarningAsError
	if newWarningAsError {
		utilLog.LogInfo("All warnings are now errors")
	}
}

func (settings *LoggerSettings) SetSuppressDebug(newSuppressDebug bool) {
	settings.suppressDebug = newSuppressDebug
	if settings.suppressDebug {
		utilLog.LogInfo("Suppressing all debug logging")
	}
}
func init() {
	settings = newLoggerSettings()
	utilLog = settings.createLogger("util")
	serverLog = settings.createLogger("log_server")
	settings.SetHook(GetDefaultHook())
	internal.SetHelperLogger(utilLog)
}

func GetSettings() *LoggerSettings {
	return &settings
}

func (s *LoggerSettings) Reset() *LoggerSettings {
	s.file = nil
	s.panicOnError = true
	s.suppressDebug = false
	s.showColor = true
	s.writer = bufio.NewWriter(os.Stderr)
	s.loggers = make(map[string]*Logger)

	settings = newLoggerSettings()
	utilLog = settings.createLogger("util")
	serverLog = settings.createLogger("log_server")
	settings.SetHook(GetDefaultHook())
	internal.SetHelperLogger(utilLog)

	return GetSettings()
}

func (settings *LoggerSettings) EnableUdpBroadcast(port int) error {
	settings.server = internal.NewLoggerServer(port)
	settings.server.SetLogger(serverLog)

	if err := settings.server.Run(); err != nil {
		utilLog.LogWarningf("unable to start server %s", err.Error())
		return err
	} else {
		settings.SetHook(GetDefaultServerHook())
		utilLog.LogDebugf("started udp broadcast on: %d", port)
	}
	return nil
}

func (settings *LoggerSettings) DisableUdpBroadcast() {
	settings.SetHook(GetDefaultServerHook())
	settings.server.Disable()
	utilLog.LogDebugf("disabled udp broadcast")
}

func (settings *LoggerSettings) WaitForUdpBroadcast() error {
	if !settings.server.IsEnabled() {
		return fmt.Errorf("cannot wait for connection, server not enabled. maybe you forgot to call \"EnableUdpBroadcast(port)\"")
	}
	for {
		if settings.server.IsConnected() {
			break
		}
	}
	return nil
}

func (settings LoggerSettings) GetLoggerNames() []string {
	retVal := make([]string, 0)
	for k, _ := range settings.loggers {
		retVal = append(retVal, k)
	}
	return retVal
}

func (settings *LoggerSettings) createLogger(loggerName string) *Logger {
	logger, _ := settings.loggers[loggerName]

	settings.loggers[loggerName] = new(Logger)
	settings.loggers[loggerName].settings = settings
	settings.loggers[loggerName].LoggerName = loggerName
	settings.loggers[loggerName].debug = settings.defaultLog
	logger = settings.loggers[loggerName]
	settings.loggers[loggerName] = logger

	return logger
}

func (settings *LoggerSettings) GetLogger(loggerName string) *Logger {
	logger, ok := settings.loggers[loggerName]
	if !ok {
		settings.loggers[loggerName] = new(Logger)
		settings.loggers[loggerName].settings = GetSettings()
		settings.loggers[loggerName].LoggerName = loggerName
		settings.loggers[loggerName].debug = settings.defaultLog
		logger = settings.loggers[loggerName]
		GetSettings().loggers[loggerName] = logger
	}
	return logger
}

func (settings *LoggerSettings) SetOutputFile(file *os.File) {
	settings.writer = bufio.NewWriter(file)
}

func (settings *LoggerSettings) FormatAnsiiColours(show bool) {
	settings.showColor = show
}

func (settings *LoggerSettings) WriteToFile(fileName string) {
	var err error
	settings.file, err = os.Create(fileName)
	if err != nil {
		panic(err)
	}
	settings.writer = bufio.NewWriter(settings.file)
	log.SetOutput(settings.writer)
}

func (settings *LoggerSettings) output(loggerName string, format string) {
	now := time.Now()

	pc, _, line, _ := runtime.Caller(3)
	caller := strings.Split(runtime.FuncForPC(pc).Name(), "/")

	settings.loggerHook(now, loggerName, fmt.Sprintf("[%s:%d] ", caller[len(caller)-1], line), format)
	//settings.loggerHook(printString)
}

func (settings *LoggerSettings) Close() {
	if settings.showColor && settings.file != nil {
		_, _ = settings.writer.WriteString("\x1b[0m")
	}
	internal.CheckError(settings.writer.Flush())
	if settings.file != nil {
		internal.CheckError(settings.file.Close())
	}
}

func (settings *LoggerSettings) LogDebug(name string, format string) {
	if settings.suppressDebug {
		return
	}
	if settings.showColor && settings.file == nil {
		format = "\033[1;30m" + format + "\033[0m"
	}
	settings.output("[Debug]["+name+"]", format)
}

func (settings *LoggerSettings) LogWarning(name string, format string) {
	if settings.warningAsError {
		settings.LogError("[ Warn]"+name, format)
	} else {
		if settings.showColor && settings.file == nil {
			format = "\033[33;1m" + format + "\033[0m"
		}
		settings.output("[ Warn]["+name+"]", format)
	}
}

func (settings *LoggerSettings) LogInfo(name string, format string) {
	settings.output("[ Info]["+name+"]", format)
}

func (settings *LoggerSettings) LogError(name string, format string) {
	if settings.panicOnError {
		defer func() {
			if r := recover(); r != nil {
				debug.Stack()
				settings.output("[Error]["+name+"]", fmt.Sprintf("error: %s", r))
				settings.output("[Error]["+name+"]", "stacktrace from panic: \n"+string(internal.StackWithGoroutine()))
				os.Exit(1)
			}
		}()
		panic(format)
	}
	if settings.showColor && settings.file == nil {
		format = "\x1b[0;31m" + format + "\x1b[0m"
	}
	settings.output("[Error]["+name+"]", format)
}

func (settings *LoggerSettings) SetHook(hook HookFunc) {
	settings.loggerHook = hook
}
