package tendermint

import (
	log "github.com/sirupsen/logrus"
	tmlog "github.com/tendermint/tendermint/libs/log"
)

// Logger is the implementation of logger
type Logger struct {
	logger tmlog.Logger
}

// NewLogger is the constructor of Logger
func NewLogger() *Logger {
	logger := tmlog.NewTMLogger(log.StandardLogger().Out)
	return &Logger{
		logger: logger,
	}
}

// Debug is the implementation of interface
func (l *Logger) Debug(msg string, keyvals ...interface{}) {
	// do nothing
}

// Info is the implementation of interface
func (l *Logger) Info(msg string, keyvals ...interface{}) {
	// l.logger.Info(msg, keyvals)
}

// Error is the implementation of interface
func (l *Logger) Error(msg string, keyvals ...interface{}) {
	l.logger.Error(msg, keyvals)
}

// With is the implementation of interface
func (l *Logger) With(keyvals ...interface{}) tmlog.Logger {
	return l
}
