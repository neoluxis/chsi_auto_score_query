package logger

import (
	"fmt"
	"log"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

var logLevel LogLevel

func Init(level string) {
	switch level {
	case "debug":
		logLevel = DEBUG
	case "warn":
		logLevel = WARN
	case "error":
		logLevel = ERROR
	default:
		logLevel = INFO
	}
}

func formatLog(level string, msg string) string {
	return fmt.Sprintf("[%s] %s: %s", time.Now().Format("2006-01-02 15:04:05"), level, msg)
}

func Debug(msg string, args ...interface{}) {
	if logLevel <= DEBUG {
		log.Println(formatLog("DEBUG", fmt.Sprintf(msg, args...)))
	}
}

func Info(msg string, args ...interface{}) {
	if logLevel <= INFO {
		log.Println(formatLog("INFO", fmt.Sprintf(msg, args...)))
	}
}

func Warn(msg string, args ...interface{}) {
	if logLevel <= WARN {
		log.Println(formatLog("WARN", fmt.Sprintf(msg, args...)))
	}
}

func Error(msg string, args ...interface{}) {
	if logLevel <= ERROR {
		log.Println(formatLog("ERROR", fmt.Sprintf(msg, args...)))
	}
}
