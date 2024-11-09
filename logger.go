package blnkgo

import "log"

// Logger interface for custom loggers
type Logger interface {
	Info(msg string)
	Error(msg string)
}

type DefaultLogger struct {
	logger *log.Logger
}

func (l *DefaultLogger) Info(msg string) {
	l.logger.Println("[INFO]", msg)
}

func (l *DefaultLogger) Error(msg string) {
	l.logger.Println("[ERROR]", msg)
}

// NewDefaultLogger
func NewDefaultLogger() *DefaultLogger {
	return &DefaultLogger{
		logger: log.Default(),
	}
}
