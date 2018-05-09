/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package indy

import (
	"github.com/hyperledger/indy-sdk-go/common/callback"
	"github.com/hyperledger/indy-sdk-go/common/indyerror"
	"github.com/hyperledger/indy-sdk-go/common/types"
)

/*
#cgo CFLAGS: -I${SRCDIR}/../../../../../../../../libindy/include
#cgo CFLAGS: -I/home/indy/libindy/include
#cgo LDFLAGS: -lindy

#include <indy_mod.h>
#include <indy_types.h>
#include <indy_did.h>
*/
import "C"

func CreateAndStoreMyDID(walletHandle types.Handle, didJSON string, cb callback.Callback) error {
	csDidJSON := newChar(didJSON)
	defer freeChar(csDidJSON)

	handle := callback.Register(cb)
	errCode := C.indy_create_and_store_my_did((C.indy_handle_t)(handle), (C.indy_handle_t)(walletHandle), csDidJSON, String2())
	return indyerror.New(int32(errCode))
}

func KeyForDID(poolHandle types.Handle, walletHandle types.Handle, did string, cb callback.Callback) error {
	csDID := newChar(did)
	defer freeChar(csDID)

	handle := callback.Register(cb)
	errCode := C.indy_key_for_did((C.indy_handle_t)(handle), (C.indy_handle_t)(poolHandle), (C.indy_handle_t)(walletHandle), csDID, String())
	return indyerror.New(int32(errCode))
}
