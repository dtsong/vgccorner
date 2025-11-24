package observability

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func TestNewLogger(t *testing.T) {
	logger := NewLogger()
	if logger == nil {
		t.Fatal("expected logger to be non-nil")
	}
}

func TestLoggerInfof(t *testing.T) {
	// Capture log output
	var buf bytes.Buffer
	logger := &Logger{
		Logger: log.New(&buf, "", 0),
	}

	logger.Infof("test message: %s", "hello")

	output := buf.String()
	if !strings.Contains(output, "[INFO] test message: hello") {
		t.Errorf("expected log to contain '[INFO] test message: hello', got: %s", output)
	}
}

func TestLoggerErrorf(t *testing.T) {
	var buf bytes.Buffer
	logger := &Logger{
		Logger: log.New(&buf, "", 0),
	}

	logger.Errorf("error occurred: %s", "failed")

	output := buf.String()
	if !strings.Contains(output, "[ERROR] error occurred: failed") {
		t.Errorf("expected log to contain '[ERROR] error occurred: failed', got: %s", output)
	}
}

func TestLoggerFatalf(t *testing.T) {
	// We can't actually test Fatalf since it calls os.Exit
	// But we can test that the logger is set up correctly
	var buf bytes.Buffer
	logger := &Logger{
		Logger: log.New(&buf, "", 0),
	}

	// Just verify the logger is functional - we can't test Fatalf without it exiting
	if logger.Logger == nil {
		t.Error("expected logger to be initialized")
	}
}

func TestLoggerFormatting(t *testing.T) {
	var buf bytes.Buffer
	logger := &Logger{
		Logger: log.New(&buf, "", 0),
	}

	tests := []struct {
		name     string
		logFunc  func(string, ...interface{})
		format   string
		args     []interface{}
		expected string
	}{
		{
			name:     "info with string",
			logFunc:  logger.Infof,
			format:   "user %s logged in",
			args:     []interface{}{"alice"},
			expected: "[INFO] user alice logged in",
		},
		{
			name:     "info with multiple args",
			logFunc:  logger.Infof,
			format:   "count: %d, name: %s",
			args:     []interface{}{42, "test"},
			expected: "[INFO] count: 42, name: test",
		},
		{
			name:     "error with format",
			logFunc:  logger.Errorf,
			format:   "failed to connect: %v",
			args:     []interface{}{"connection refused"},
			expected: "[ERROR] failed to connect: connection refused",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			tt.logFunc(tt.format, tt.args...)
			output := buf.String()
			if !strings.Contains(output, tt.expected) {
				t.Errorf("expected log to contain %q, got: %s", tt.expected, output)
			}
		})
	}
}
