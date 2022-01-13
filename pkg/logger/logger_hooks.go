package logger

import (
	"bufio"
	"github.com/elamre/logger/internal"
	"os"
	"time"
)

func GetDefaultServerHook() HookFunc {
	return GetSettings().serverHook
}

func (settings *LoggerSettings) serverHook(now time.Time, loggerName string, fileName string, format string) {
	if settings.server.IsConnected() {
		if err := settings.server.SendMessageGob(now, loggerName, fileName, format); err != nil {
			//if err := settings.server.SendMessage(internal.FormatString(now, loggerName, fileName, format)); err != nil {
			settings.SetHook(GetDefaultHook())
			utilLog.LogWarningf("something went wrong with writing to the server: %s\n\tdisabling server", err.Error()) // TODO
			settings.defaultHook(now, loggerName, fileName, format)
		}
	} else {
		settings.defaultHook(now, loggerName, fileName, format)
	}
}

func (settings *LoggerSettings) defaultHook(now time.Time, loggerName string, fileName string, format string) {
	s := internal.FormatString(now, loggerName, fileName, format)
	if _, err := settings.writer.WriteString(s); err != nil {
		panic(err)
	}
	if err := settings.writer.Flush(); err != nil {
		panic(err)
	}
}

func GetDefaultHook() HookFunc {
	GetSettings().writer = bufio.NewWriter(os.Stderr)
	return settings.defaultHook
}

func (settings *LoggerSettings) fileHook(now time.Time, loggerName string, fileName string, format string) {
	s := internal.FormatString(now, loggerName, fileName, format)
	if _, err := settings.writer.WriteString(s); err != nil {
		panic(err)
	}
	if err := settings.writer.Flush(); err != nil {
		panic(err)
	}
}

func GetDefaultFileHook(fileName string) HookFunc {
	var err error
	GetSettings().file, err = os.Create(fileName)
	if err != nil {
		panic(err)
	}
	GetSettings().writer = bufio.NewWriter(GetSettings().file)
	return GetSettings().fileHook
}
