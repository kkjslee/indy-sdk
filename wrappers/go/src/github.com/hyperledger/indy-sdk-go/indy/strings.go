/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package indy

/*
#include <stdlib.h>
*/
import "C"
import "unsafe"

// New creates a C string out of a Go string
func newChar(s string) *C.char {
	return C.CString(s)
}

// Free frees the given C string
func freeChar(cs *C.char) {
	C.free(unsafe.Pointer(cs))
}
