package logger

import (
	"fmt"
	"runtime"
	"strings"
)

type Logger struct {
	debug      bool
	settings   *LoggerSettings
	LoggerName string
}

func NewNamedLogger(name string) *Logger {
	return GetSettings().GetLogger(name)
}

func NewLogger() *Logger {
	_, file, _, _ := runtime.Caller(1)
	sfile := strings.Split(file, "/")
	name := strings.Split(sfile[len(sfile)-1], ".")[0]
	return GetSettings().GetLogger(name)
}

func (i *Logger) PrintDebug(debug bool) *Logger {
	i.debug = debug
	return i
}

func (i *Logger) DebugSet() bool {
	return i.debug
}

func (i *Logger) LogDebug(format string) {
	if !i.debug {
		return
	}
	i.settings.LogDebug(i.LoggerName, format)
}

func (i *Logger) LogInfo(format string) {
	i.settings.LogInfo(i.LoggerName, format)
}

func (i *Logger) LogWarning(format string) {
	i.settings.LogWarning(i.LoggerName, format)
}

func (i *Logger) LogError(format string) {
	i.settings.LogError(i.LoggerName, format)
}

func (i *Logger) LogDebugf(format string, a ...interface{}) {
	if !i.debug {
		return
	}
	i.settings.LogDebug(i.LoggerName, fmt.Sprintf(format, a...))
}

func (i *Logger) LogInfof(format string, a ...interface{}) {
	i.settings.LogInfo(i.LoggerName, fmt.Sprintf(format, a...))
}

func (i *Logger) LogWarningf(format string, a ...interface{}) {
	i.settings.LogWarning(i.LoggerName, fmt.Sprintf(format, a...))
}

func (i *Logger) LogErrorf(format string, a ...interface{}) {
	i.settings.LogError(i.LoggerName, fmt.Sprintf(format, a...))
}
