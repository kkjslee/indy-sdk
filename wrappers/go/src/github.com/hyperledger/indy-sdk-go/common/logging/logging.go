/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package logging

// Logger allows the aplication to output logs at various log levels
type Logger interface {
	Debugf(msg string, args ...interface{})
	Infof(msg string, args ...interface{})
	Warnf(msg string, args ...interface{})
	Errorf(msg string, args ...interface{})
}

// Level is the log level
type Level int

const (
	// DEBUG debug log level
	DEBUG Level = iota
	// INFO info log level
	INFO
	// WARN warning log level
	WARN
	// ERROR error log level
	ERROR
)

var level = INFO

// SetLevel sets the log level
func SetLevel(l Level) {
	level = l
}

// MustGetLogger returns the logger
func MustGetLogger(module string) Logger {
	return &logger{
		module: module,
	}
}
