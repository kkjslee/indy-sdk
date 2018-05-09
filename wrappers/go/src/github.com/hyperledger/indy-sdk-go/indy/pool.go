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

#include <stdlib.h>
#include <indy_pool.h>
*/
import "C"

func CreatePoolLedgerConfig(name, configPath string, cb callback.Callback) error {
	csName := newChar(name)
	defer freeChar(csName)

	csConfigPath := newChar(configPath)
	defer freeChar(csConfigPath)

	handle := callback.Register(cb)
	errCode := C.indy_create_pool_ledger_config((C.indy_handle_t)(handle), csName, csConfigPath, Default())
	return indyerror.New(int32(errCode))
}

func DeletePoolLedgerConfig(name string, cb callback.Callback) error {
	csName := newChar(name)
	defer freeChar(csName)

	handle := callback.Register(cb)
	errCode := C.indy_delete_pool_ledger_config((C.indy_handle_t)(handle), csName, Default())
	return indyerror.New(int32(errCode))
}

func OpenPoolLedger(name, config string, cb callback.Callback) error {
	csName := newChar(name)
	defer freeChar(csName)

	var csConfig *C.char
	if config != "" {
		csConfig = newChar(config)
		defer freeChar(csConfig)
	}

	handle := callback.Register(cb)
	errCode := C.indy_open_pool_ledger((C.indy_handle_t)(handle), csName, csConfig, SingleHandle())
	return indyerror.New(int32(errCode))
}

func ListPools(cb callback.Callback) error {
	handle := callback.Register(cb)
	errCode := C.indy_list_pools((C.indy_handle_t)(handle), String())
	return indyerror.New(int32(errCode))
}

func RefreshPoolLedger(poolHandle types.Handle, cb callback.Callback) error {
	handle := callback.Register(cb)
	errCode := C.indy_refresh_pool_ledger((C.indy_handle_t)(handle), (C.indy_handle_t)(poolHandle), Default())
	return indyerror.New(int32(errCode))
}

func ClosePoolLedger(poolHandle types.Handle, cb callback.Callback) error {
	handle := callback.Register(cb)
	errCode := C.indy_close_pool_ledger((C.indy_handle_t)(handle), (C.indy_handle_t)(poolHandle), Default())
	return indyerror.New(int32(errCode))
}
