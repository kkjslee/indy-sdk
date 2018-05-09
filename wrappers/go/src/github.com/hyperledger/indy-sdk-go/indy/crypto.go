/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package indy

import (
	"unsafe"

	"github.com/hyperledger/indy-sdk-go/common/callback"
	"github.com/hyperledger/indy-sdk-go/common/indyerror"
	"github.com/hyperledger/indy-sdk-go/common/types"
)

/*
#cgo CFLAGS: -I${SRCDIR}/../../../../../../../../libindy/include
#cgo CFLAGS: -I/home/indy/libindy/include
#cgo LDFLAGS: -lindy

#include <stdlib.h>
#include <indy_mod.h>
#include <indy_types.h>
#include <indy_crypto.h>
*/
import "C"

func AnonCrypt(recipientVK string, message []byte, cb callback.Callback) error {
	csRecipientVK := newChar(recipientVK)
	defer freeChar(csRecipientVK)

	cbMessage := C.CBytes(message)
	defer C.free(unsafe.Pointer(cbMessage))

	handle := callback.Register(cb)
	errCode := C.indy_crypto_anon_crypt((C.indy_handle_t)(handle), csRecipientVK, (*C.indy_u8_t)(cbMessage), (C.indy_u32_t)(len(message)), Bytes())
	return indyerror.New(int32(errCode))
}

func AnonDecrypt(walletHandle types.Handle, recipientVK string, message []byte, cb callback.Callback) error {
	csRecipientVK := newChar(recipientVK)
	defer freeChar(csRecipientVK)

	cbMessage := C.CBytes(message)
	defer C.free(unsafe.Pointer(cbMessage))

	handle := callback.Register(cb)
	errCode := C.indy_crypto_anon_decrypt((C.indy_handle_t)(handle), (C.indy_handle_t)(walletHandle), csRecipientVK, (*C.indy_u8_t)(cbMessage), (C.indy_u32_t)(len(message)), Bytes())
	return indyerror.New(int32(errCode))
}

func AuthCrypt(walletHandle types.Handle, senderVK, recipientVK string, message []byte, cb callback.Callback) error {
	csRecipientVK := newChar(recipientVK)
	defer freeChar(csRecipientVK)

	csSenderVK := newChar(senderVK)
	defer freeChar(csSenderVK)

	cbMessage := C.CBytes(message)
	defer C.free(unsafe.Pointer(cbMessage))

	handle := callback.Register(cb)
	errCode := C.indy_crypto_auth_crypt((C.indy_handle_t)(handle), (C.indy_handle_t)(walletHandle), csSenderVK, csRecipientVK, (*C.indy_u8_t)(cbMessage), (C.indy_u32_t)(len(message)), Bytes())
	return indyerror.New(int32(errCode))
}

func AuthDecrypt(walletHandle types.Handle, recipientVK string, message []byte, cb callback.Callback) error {
	csRecipientVK := newChar(recipientVK)
	defer freeChar(csRecipientVK)

	cbMessage := C.CBytes(message)
	defer C.free(unsafe.Pointer(cbMessage))

	handle := callback.Register(cb)
	errCode := C.indy_crypto_auth_decrypt((C.indy_handle_t)(handle), (C.indy_handle_t)(walletHandle), csRecipientVK, (*C.indy_u8_t)(cbMessage), (C.indy_u32_t)(len(message)), StringAndBytes())
	return indyerror.New(int32(errCode))
}
