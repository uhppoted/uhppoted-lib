package log

import (
	"fmt"
	syslog "log"
)

type LogLevel int

const (
	none LogLevel = iota
	debug
	info
	warn
	errors
)

var log = syslog.Default()
var debugging = false
var level = info
var hook func()

func SetDebug(enabled bool) {
	debugging = enabled
}

func SetLevel(l string) {
	switch l {
	case "none":
		level = none
	case "debug":
		level = debug
	case "info":
		level = info
	case "warn":
		level = warn
	case "error":
		level = errors
	}
}

func SetLogger(l *syslog.Logger) {
	log = l
}

func SetFatalHook(f func()) {
	hook = f
}

func Debugf(format string, args ...any) {
	if debugging || level < info {
		log.Printf("%-5v  %v", "DEBUG", fmt.Sprintf(format, args...))
	}
}

func Infof(format string, args ...any) {
	if level < warn {
		log.Printf("%-5v  %v", "INFO", fmt.Sprintf(format, args...))
	}
}

func Warnf(format string, args ...any) {
	if level < errors {
		log.Printf("%-5v  %v", "WARN", fmt.Sprintf(format, args...))
	}
}

func Errorf(format string, args ...any) {
	log.Printf("%-5v  %v", "ERROR", fmt.Sprintf(format, args...))
}

func Fatalf(format string, args ...any) {
	if hook != nil {
		hook()
	}

	log.Fatalf("%-5v  %v", "FATAL", fmt.Sprintf(format, args...))
}
