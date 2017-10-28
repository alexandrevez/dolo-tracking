package logger

import (
	"bufio"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
)

var log *logrus.Logger

func initLogger(level logrus.Level) {
	log = logrus.New()
	log.Level = level
	log.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
	}
}

func decorateRuntimeContext(msg string) string {
	if log.Level.String() != logrus.DebugLevel.String() {
		return msg
	}

	if pc, file, line, ok := runtime.Caller(2); ok {
		funcName := runtime.FuncForPC(pc).Name()
		return msg + "\n\t" + funcName + ": " + file + ":" + strconv.Itoa(line)
	}
	return msg
}
func decorateStack(msg string) string {
	msg += "\n"

	stackStr := string(debug.Stack())
	stackScanner := bufio.NewScanner(strings.NewReader(stackStr))
	i := 0
	for stackScanner.Scan() {
		i++
		if i < 8 {
			continue
		}
		msg += "\t" + stackScanner.Text() + "\n"
	}
	return msg
}

// Debug prints a debug message to logger
func Debug(msg string) {
	if log == nil {
		initLogger(logrus.DebugLevel)
	}
	log.Debugln(decorateRuntimeContext(msg))
}

// Warn prints a warning message to logger
func Warn(msg string) {
	if log == nil {
		initLogger(logrus.DebugLevel)
	}

	log.Warnln(decorateStack(msg))
}

// Error prints an error message to logger
func Error(msg string) {
	if log == nil {
		initLogger(logrus.DebugLevel)
	}
	log.Errorln(decorateStack(msg))
}
