package logger

import (
	"bytes"
	"log"
	"testing"
)

func TestLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0) // Установка флага в 0 убирает временную метку

	tests := []struct {
		name     string
		level    string
		enabled  bool
		logFunc  func(*Logger)
		expected string
	}{
		{
			name:     "Error logging enabled",
			level:    "error",
			enabled:  true,
			logFunc:  func(l *Logger) { l.Error("This is an error message") },
			expected: "[ERROR] This is an error message\n",
		},
		{
			name:     "Warn logging enabled",
			level:    "warn",
			enabled:  true,
			logFunc:  func(l *Logger) { l.Warn("This is a warning message") },
			expected: "[WARN] This is a warning message\n",
		},
		{
			name:     "Info logging enabled",
			level:    "info",
			enabled:  true,
			logFunc:  func(l *Logger) { l.Info("This is an info message") },
			expected: "[INFO] This is an info message\n",
		},
		{
			name:     "Debug logging enabled",
			level:    "debug",
			enabled:  true,
			logFunc:  func(l *Logger) { l.Debug("This is a debug message") },
			expected: "[DEBUG] This is a debug message\n",
		},
		{
			name:     "Error logging disabled",
			level:    "error",
			enabled:  false,
			logFunc:  func(l *Logger) { l.Error("This should not be logged") },
			expected: "",
		},
		{
			name:     "Debug logging not logged at warn level",
			level:    "warn",
			enabled:  true,
			logFunc:  func(l *Logger) { l.Debug("This should not be logged") },
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Logger{
				logger:    logger,
				logLevel:  parseLogLevel(tt.level),
				isEnabled: tt.enabled,
			}
			buf.Reset() // Сброс буфера перед каждым тестом
			tt.logFunc(l)

			if got := buf.String(); got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}
