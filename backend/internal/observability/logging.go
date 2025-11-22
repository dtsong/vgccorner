package observability

import (
	"log"
)

type Logger struct {
	*log.Logger
}

// You can later swap this for zerolog/zap/etc without changing callers.
func NewLogger() *Logger {
	return &Logger{Logger: log.Default()}
}

func (l *Logger) Infof(format string, args ...any) {
	l.Printf("[INFO] "+format, args...)
}

func (l *Logger) Errorf(format string, args ...any) {
	l.Printf("[ERROR] "+format, args...)
}

func (l *Logger) Fatalf(format string, args ...any) {
	l.Logger.Fatalf("[FATAL] "+format, args...)
}
