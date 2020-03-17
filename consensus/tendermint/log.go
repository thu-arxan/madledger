// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package tendermint

import (
	logrus "github.com/sirupsen/logrus"
	tmlog "github.com/tendermint/tendermint/libs/log"
)

// Logger is the implementation of logger
type Logger struct {
	logger tmlog.Logger
}

// NewLogger is the constructor of Logger
func NewLogger() *Logger {
	logger := tmlog.NewTMLogger(logrus.StandardLogger().Out)
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
	// l.logger.Info(msg, "info", keyvals)
}

// Error is the implementation of interface
func (l *Logger) Error(msg string, keyvals ...interface{}) {
	// l.logger.Error(msg, "error", keyvals)
}

// With is the implementation of interface
func (l *Logger) With(keyvals ...interface{}) tmlog.Logger {
	return l
}
