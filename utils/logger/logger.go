package logger

import (
	"errors"
	"fmt"
	"time"
)

const (
	// NONE - not logging level
	NONE = 0
	// ERROR - error logging level
	ERROR = 1
	// WARN - warning logging level
	WARN = 2
	// INFO - info logging level
	INFO = 3
	// DEBUG - debug logging level
	DEBUG = 4
)

func logTimeDecorator(log func(string), text string) {
	t := time.Now()
	timestring := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	str := fmt.Sprintf("[%s]: %s", timestring, text)
	log(str)
}

func logLevelDecorator(lvl int, log func(string), text string) {
	switch lvl {
	case NONE:
	case ERROR:
		str := fmt.Sprintf("[ERROR]: %s", text)
		logTimeDecorator(log, str)
	case WARN:
		str := fmt.Sprintf("[WARNING]: %s", text)
		log(str)
	case INFO:
		str := fmt.Sprintf("[INFO]: %s", text)
		log(str)
	case DEBUG:
		str := fmt.Sprintf("[DEBUG]: %s", text)
		logTimeDecorator(log, str)
	default:
		panic(errors.New("Unknown logger level"))
	}
}

func log(text string) {
	fmt.Println(text)
}

// Error - log error
func Error(text string) {
	logLevelDecorator(ERROR, log, text)
}

// Warning - log warning
func Warning(text string) {
	logLevelDecorator(WARN, log, text)
}

// Info - log info
func Info(text string) {
	logLevelDecorator(INFO, log, text)
}

// Debug - log debug info
func Debug(text string) {
	logLevelDecorator(DEBUG, log, text)
}

// None - log nothing
func None(text string) {
	logLevelDecorator(NONE, log, text)
}
