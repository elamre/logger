package internal

import (
	"runtime"
	"time"
)

type LoggerMasterInterface interface {
	LogError(string, string)
	LogDebug(string, string)
	LogInfo(string, string)
	LogWarning(string, string)
	GetLogger(string) *LoggerInterface
}

type LoggerInterface interface {
	LogError(string)
	LogErrorf(string, ...interface{})
	LogDebug(string)
	LogDebugf(string, ...interface{})
	LogInfo(string)
	LogInfof(string, ...interface{})
	LogWarning(string)
	LogWarningf(string, ...interface{})
}

var helperLogger LoggerInterface

func SetHelperLogger(logger LoggerInterface) {
	helperLogger = logger
}

func CheckError(err error) {
	if err != nil {
		if helperLogger != nil {
			panic(err)
		} else {
			helperLogger.LogError(err.Error())
		}
	}
}

func StackWithGoroutine() []byte {
	buf := make([]byte, 2048)
	for {
		n := runtime.Stack(buf, true)
		if n < len(buf) {
			return buf[:n]
		}
		buf = make([]byte, 2*len(buf))
	}
}

// Cheap integer to fixed-width decimal ASCII. Give a negative width to avoid zero-padding.
func itoa(buf []byte, i int, wid int) (newpos int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	for i := bp; i < len(b); i++ {
		buf[i-bp] = b[i]
	}
	return len(b) - bp
}

var timeBuf [20]byte

func FormatString(now time.Time, loggerName string, fileName string, format string) string {
	hour, min, sec := now.Clock()
	year, month, day := now.Date()
	pos := itoa(timeBuf[0:], year, 4)
	timeBuf[pos] = '/'
	pos++
	pos += itoa(timeBuf[pos:], int(month), 2)
	timeBuf[pos] = '/'
	pos++
	pos += itoa(timeBuf[pos:], day, 2)
	timeBuf[pos] = ' '
	pos++
	pos += itoa(timeBuf[pos:], hour, 2)
	timeBuf[pos] = ':'
	pos++
	pos += itoa(timeBuf[pos:], min, 2)
	timeBuf[pos] = ':'
	pos++
	pos += itoa(timeBuf[pos:], sec, 2)
	timeBuf[pos] = ' '

	printString := string(timeBuf[:]) + loggerName + fileName + format
	if len(format) == 0 || format[len(format)-1] != '\n' {
		printString += "\n"
	}
	return printString
}
