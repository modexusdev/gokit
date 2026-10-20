package log

import (
	"fmt"
	"strings"
	"time"
)

type Logger struct {
	name string
}

func New(name string) *Logger {
	return &Logger{
		name: name,
	}
}

func (l *Logger) Info(messages ...string) {
	l.print("INFO", colorGreen, messages...)
}

func (l *Logger) Warn(messages ...string) {
	l.print("WARN", colorYellow, messages...)
}

func (l *Logger) Error(messages ...string) {
	l.print("ERROR", colorRed, messages...)
}

func (l *Logger) Debug(messages ...string) {
	l.print("DEBUG", colorCyan, messages...)
}

func (l *Logger) print(
	level string,
	levelColor string,
	messages ...string,
) {
	now := time.Now().Format("15:04:05")
	message := strings.Join(messages, " ")

	fmt.Printf(
		"%s[%s]%s %s[%s]%s %s%-8s%s",
		colorGray,
		now,
		colorReset,

		colorBrightMagenta,
		l.name,
		colorReset,

		levelColor,
		level,
		colorReset,
	)

	if message != "" {
		fmt.Printf(" %s", message)
	}

	fmt.Println()
}
