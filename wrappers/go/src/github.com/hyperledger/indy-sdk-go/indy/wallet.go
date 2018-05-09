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
#include <indy_wallet.h>
*/
import "C"

func CreateWallet(poolName, name, xtype, config, credentials string, cb callback.Callback) error {
	csPoolName := newChar(poolName)
	defer freeChar(csPoolName)

	csName := newChar(name)
	defer freeChar(csName)

	var csType *C.char
	if xtype != "" {
		csType = newChar(xtype)
		defer freeChar(csType)
	}

	var csConfig *C.char
	if config != "" {
		csConfig = newChar(config)
		defer freeChar(csConfig)
	}

	var csCredentials *C.char
	if credentials != "" {
		csCredentials = newChar(credentials)
		defer freeChar(csCredentials)
	}

	handle := callback.Register(cb)
	errCode := C.indy_create_wallet((C.indy_handle_t)(handle), csPoolName, csName, csType, csConfig, csCredentials, Default())
	return indyerror.New(int32(errCode))
}

func DeleteWallet(name, credentials string, cb callback.Callback) error {
	csName := newChar(name)
	defer freeChar(csName)

	var csCredentials *C.char
	if credentials != "" {
		csCredentials = C.CString(credentials)
		defer freeChar(csCredentials)
	}

	handle := callback.Register(cb)
	errCode := C.indy_delete_wallet((C.indy_handle_t)(handle), csName, csCredentials, Default())
	return indyerror.New(int32(errCode))
}

func OpenWallet(name, config, credentials string, cb callback.Callback) error {
	csName := C.CString(name)
	defer freeChar(csName)

	var csConfig *C.char
	if config != "" {
		csConfig = C.CString(config)
		defer freeChar(csConfig)
	}

	var csCredentials *C.char
	if credentials != "" {
		csCredentials = C.CString(credentials)
		defer freeChar(csCredentials)
	}

	handle := callback.Register(cb)
	errCode := C.indy_open_wallet((C.indy_handle_t)(handle), csName, csConfig, csCredentials, SingleHandle())
	return indyerror.New(int32(errCode))
}

func CloseWallet(walletHandle types.Handle, cb callback.Callback) error {
	handle := callback.Register(cb)
	errCode := C.indy_close_wallet((C.indy_handle_t)(handle), (C.indy_handle_t)(walletHandle), Default())
	return indyerror.New(int32(errCode))
}
