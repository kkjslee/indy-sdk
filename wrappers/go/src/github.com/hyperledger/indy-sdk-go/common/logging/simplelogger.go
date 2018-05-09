/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package logging

import "fmt"

// logger is a temporary logger that will be replaced with a suitable implementation
type logger struct {
	module string
}

func (l *logger) Debugf(msg string, args ...interface{}) {
	if level <= DEBUG {
		fmt.Printf("DEBUG: "+msg+"\n", args...)
	}
}

func (l *logger) Infof(msg string, args ...interface{}) {
	if level <= INFO {
		fmt.Printf("INFO: "+msg+"\n", args...)
	}
}

func (l *logger) Warnf(msg string, args ...interface{}) {
	if level <= WARN {
		fmt.Printf("WARN: "+msg+"\n", args...)
	}
}

func (l *logger) Errorf(msg string, args ...interface{}) {
	if level <= ERROR {
		fmt.Printf("ERROR: "+msg+"\n", args...)
	}
}
