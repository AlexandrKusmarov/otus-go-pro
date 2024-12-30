package logger

import (
	"log"
	"os"
	"strings"
)

type Log interface {
	Error(msg string, a ...any)
	Warn(msg string, a ...any)
	Info(msg string, a ...any)
	Debug(msg string, a ...any)
}

type Logger struct {
	logger    *log.Logger
	logLevel  int
	isEnabled bool
}

const (
	LevelError = iota
	LevelWarn
	LevelInfo
	LevelDebug
)

// New создает новый логгер с заданным уровнем логирования и настройками.

func New(enabled bool, level string, prefix string) *Logger {
	logLevel := parseLogLevel(level)

	return &Logger{
		logger:    log.New(os.Stdout, prefix, log.LstdFlags),
		logLevel:  logLevel,
		isEnabled: enabled,
	}
}

func parseLogLevel(level string) int {
	switch strings.ToLower(level) {
	case "error":
		return LevelError
	case "warn":
		return LevelWarn
	case "info":
		return LevelInfo
	case "debug":
		return LevelDebug
	default:
		log.Printf("Неизвестный уровень логирования: %s, используется 'info'", level)
		return LevelInfo
	}
}

// Error логирует сообщение уровня ERROR.
func (l *Logger) Error(msg string, _ ...any) {
	if l.isEnabled && l.logLevel >= LevelError {
		l.logger.SetPrefix("[ERROR] ")
		l.logger.Println(msg)
	}
}

// Warn логирует сообщение уровня WARN.
func (l *Logger) Warn(msg string, _ ...any) {
	if l.isEnabled && l.logLevel >= LevelWarn {
		l.logger.SetPrefix("[WARN] ")
		l.logger.Println(msg)
	}
}

// Info логирует сообщение уровня INFO.
func (l *Logger) Info(msg string, _ ...any) {
	if l.isEnabled && l.logLevel >= LevelInfo {
		l.logger.SetPrefix("[INFO] ")
		l.logger.Println(msg)
	}
}

// Debug логирует сообщение уровня DEBUG.
func (l *Logger) Debug(msg string, _ ...any) {
	if l.isEnabled && l.logLevel >= LevelDebug {
		l.logger.SetPrefix("[DEBUG] ")
		l.logger.Println(msg)
	}
}
