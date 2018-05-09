/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package assert

import (
	"fmt"
	"reflect"
	"testing"
)

// NoError asserts that err is nil
func NoError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
}

// NoErrorf asserts that err is nil
func NoErrorf(t *testing.T, err error, msg string, args ...interface{}) {
	if err != nil {
		t.Fatalf(fmt.Sprintf(msg, args...)+" - %s", err)
	}
}

// Equal asserts that the actual value is equal to the expected value
func Equal(t *testing.T, expect, actual interface{}) {
	if !reflect.DeepEqual(expect, actual) {
		t.Fatalf("Expecting [%s] but got [%s]", expect, actual)
	}
}

// Equalf asserts that the actual value is equal to the expected value
func Equalf(t *testing.T, expect, actual interface{}, msg string, args ...interface{}) {
	if !reflect.DeepEqual(expect, actual) {
		t.Fatalf(fmt.Sprintf(msg, args...)+": Expecting [%s] but got [%s]", expect, actual)
	}
}
