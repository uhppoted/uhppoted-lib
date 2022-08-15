package log

import (
	"fmt"
	syslog "log"
	"sync"
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
var hooks = struct {
	hooks []func()
	sync.Mutex
}{
	hooks: []func(){},
}

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

func AddFatalHook(f func()) {
	hooks.Lock()
	defer hooks.Unlock()

	hooks.hooks = append(hooks.hooks, f)
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

// Executes fatal hooks in reverse order before invoking the system Log.Fatalf to exit
func Fatalf(format string, args ...any) {
	N := len(hooks.hooks)
	list := make([]func(), N)

	hooks.Lock()
	copy(list, hooks.hooks)
	hooks.Unlock()

	for i := 1; i <= N; i++ {
		list[N-i]()
	}

	log.Fatalf("%-5v  %v", "FATAL", fmt.Sprintf(format, args...))
}
