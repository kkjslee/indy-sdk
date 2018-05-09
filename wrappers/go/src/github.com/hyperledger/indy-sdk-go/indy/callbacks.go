/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package indy

import (
	"fmt"
	"unsafe"

	"github.com/hyperledger/indy-sdk-go/common/callback"
	"github.com/hyperledger/indy-sdk-go/common/indyerror"
	"github.com/hyperledger/indy-sdk-go/common/logging"
	"github.com/hyperledger/indy-sdk-go/common/types"
)

/*
#cgo CFLAGS: -I${SRCDIR}/../../../../../../../../libindy/include
#cgo CFLAGS: -I/home/indy/libindy/include
#cgo LDFLAGS: -lindy

#include <stdlib.h>
#include <stdint.h>

extern void def_callback(int32_t xcommand_handle, int32_t err);
typedef void (*callback_fcn)(int32_t, int32_t);

extern void handle_callback(int32_t xcommand_handle, int32_t err, int32_t pool_handle);
typedef void (*handle_callback_fcn)(int32_t, int32_t, int32_t);

extern void string_callback(int32_t xcommand_handle, int32_t err, char *const);
typedef void (*string_callback_fcn)(int32_t, int32_t, const char *);

extern void string2_callback(int32_t xcommand_handle, int32_t err, char *, char *);
typedef void (*string2_callback_fcn)(int32_t, int32_t, const char *, const char *);

extern void string3_callback(int32_t xcommand_handle, int32_t err, char *, char *, char *);
typedef void (*string3_callback_fcn)(int32_t, int32_t, const char *, const char *, const char *);

extern void bytes_callback(int32_t, int32_t, uint8_t *, int32_t);
typedef void (*bytes_callback_fcn)(int32_t, int32_t, const uint8_t *, int32_t);

extern void string_bytes_callback(int32_t, int32_t, char *, uint8_t *, int32_t);
typedef void (*string_bytes_callback_fcn)(int32_t, int32_t, char *, const uint8_t *, int32_t);

extern void bool_callback(int32_t xcommand_handle, int32_t err, unsigned int);
typedef void (*bool_callback_fcn)(int32_t, int32_t, unsigned int);
*/
import "C"

var logger = logging.MustGetLogger("indy-sdk")

func Default() C.callback_fcn {
	return (C.callback_fcn)(unsafe.Pointer(C.def_callback))
}

func SingleHandle() C.handle_callback_fcn {
	return (C.handle_callback_fcn)(unsafe.Pointer(C.handle_callback))
}

func String() C.string_callback_fcn {
	return (C.string_callback_fcn)(unsafe.Pointer(C.string_callback))
}

func String2() C.string2_callback_fcn {
	return (C.string2_callback_fcn)(unsafe.Pointer(C.string2_callback))
}

func String3() C.string3_callback_fcn {
	return (C.string3_callback_fcn)(unsafe.Pointer(C.string3_callback))
}

func Bytes() C.bytes_callback_fcn {
	return (C.bytes_callback_fcn)(unsafe.Pointer(C.bytes_callback))
}

func StringAndBytes() C.string_bytes_callback_fcn {
	return (C.string_bytes_callback_fcn)(unsafe.Pointer(C.string_bytes_callback))
}

func Bool() C.bool_callback_fcn {
	return (C.bool_callback_fcn)(unsafe.Pointer(C.bool_callback))
}

//export def_callback
func def_callback(handle int32, errCode int32) {
	cb, ok := callback.Remove(types.Handle(handle))
	if !ok {
		cb(fmt.Errorf("unable to find callback for handle: %d", handle), nil)
	} else {
		cb(indyerror.New(errCode), nil)
	}
}

//export handle_callback
func handle_callback(handle int32, errCode int32, poolHandle int32) {
	cb, ok := callback.Remove(types.Handle(handle))
	if !ok {
		cb(fmt.Errorf("unable to find callback for handle: %d", handle), nil)
	} else {
		cb(indyerror.New(errCode), types.Handle(poolHandle))
	}
}

//export string_callback
func string_callback(handle int32, errCode int32, s *C.char) {
	cb, ok := callback.Remove(types.Handle(handle))
	if !ok {
		cb(fmt.Errorf("unable to find callback for handle: %d", handle), nil)
	} else {
		cb(indyerror.New(errCode), C.GoString(s))
	}
}

//export string2_callback
func string2_callback(handle int32, errCode int32, s1 *C.char, s2 *C.char) {
	cb, ok := callback.Remove(types.Handle(handle))
	if !ok {
		cb(fmt.Errorf("unable to find callback for handle: %d", handle), nil)
	} else {
		cb(indyerror.New(errCode), []string{C.GoString(s1), C.GoString(s2)})
	}
}

//export string3_callback
func string3_callback(handle int32, errCode int32, s1 *C.char, s2 *C.char, s3 *C.char) {
	cb, ok := callback.Remove(types.Handle(handle))
	if !ok {
		cb(fmt.Errorf("unable to find callback for handle: %d", handle), nil)
	} else {
		cb(indyerror.New(errCode), []string{C.GoString(s1), C.GoString(s2), C.GoString(s3)})
	}
}

//export bytes_callback
func bytes_callback(handle int32, errCode int32, b *C.uchar, blength int32) {
	cb, ok := callback.Remove(types.Handle(handle))
	if !ok {
		cb(fmt.Errorf("unable to find callback for handle: %d", handle), nil)
	} else {
		cb(indyerror.New(errCode), C.GoBytes(unsafe.Pointer(b), C.int(blength)))
	}
}

//export string_bytes_callback
func string_bytes_callback(handle int32, errCode int32, s *C.char, b *C.uchar, blength int32) {
	cb, ok := callback.Remove(types.Handle(handle))
	if !ok {
		cb(fmt.Errorf("unable to find callback for handle: %d", handle), nil)
	} else {
		cb(indyerror.New(errCode), []interface{}{C.GoString(s), C.GoBytes(unsafe.Pointer(b), C.int(blength))})
	}
}

//export bool_callback
func bool_callback(handle int32, errCode int32, b C.uint) {
	cb, ok := callback.Remove(types.Handle(handle))
	if !ok {
		cb(fmt.Errorf("unable to find callback for handle: %d", handle), nil)
	} else {
		cb(indyerror.New(errCode), b == 1)
	}
}
