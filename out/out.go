package out

import (
	"fmt"
	"os"
	"time"
)

func logMessage(level, format string, v ...interface{}) {
	var color string

	switch level {
	case "ERROR":
		color = "\033[1;31m"
	case "WARN":
		color = "\033[1;33m"
	case "INFO":
		color = "\033[1;32m"
	case "FATAL":
		color = "\033[1;35m"
	case "DEBUG":
		color = "\033[1;36m"
	default:
		color = "\033[0m"
	}
	now := time.Now()
	fmt.Printf("%02d:%02d %s%s\033[0m %s\n", now.Minute(), now.Second(), color, level, fmt.Sprintf(format, v...))
}

func Print(s string, v ...interface{}) {
	logMessage("INFO", s, v...)
}

func Debug(s string, v ...interface{}) {
	logMessage("DEBUG", s, v...)
}

func Error(s string, v ...interface{}) {
	logMessage("ERROR", s, v...)
}

func Warn(s string, v ...interface{}) {
	logMessage("WARN", s, v...)
}

func Fatal(s string, v ...interface{}) {
	logMessage("FATAL", s, v...)
	os.Exit(1)
}
