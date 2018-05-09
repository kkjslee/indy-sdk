/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package indy

import (
	"github.com/hyperledger/indy-sdk-go/common/callback"
	"github.com/hyperledger/indy-sdk-go/common/indyerror"
	"github.com/hyperledger/indy-sdk-go/common/role"
	"github.com/hyperledger/indy-sdk-go/common/types"
)

/*
#cgo CFLAGS: -I${SRCDIR}/../../../../../../../../libindy/include
#cgo CFLAGS: -I/home/indy/libindy/include
#cgo LDFLAGS: -lindy

#include <indy_mod.h>
#include <indy_types.h>
#include <indy_ledger.h>
*/
import "C"

func BuildNYMRequest(submitterDID, targetDID, verkey string, alias *types.Alias, role *role.Role, cb callback.Callback) error {
	csSubmitterDID := newChar(submitterDID)
	defer freeChar(csSubmitterDID)

	csTargetDID := newChar(targetDID)
	defer freeChar(csTargetDID)

	csVerKey := newChar(verkey)
	defer freeChar(csVerKey)

	var csAlias *C.char
	if alias != nil {
		csAlias = newChar(alias.String())
		defer freeChar(csAlias)
	}

	var csRole *C.char
	if role != nil {
		csRole = newChar(role.String())
		defer freeChar(csRole)
	}

	handle := callback.Register(cb)
	errCode := C.indy_build_nym_request((C.indy_handle_t)(handle), csSubmitterDID, csTargetDID, csVerKey, csAlias, csRole, String())
	return indyerror.New(int32(errCode))
}

func SignAndSubmitRequest(poolHandle types.Handle, walletHandle types.Handle, submitterDID, requestJSON string, cb callback.Callback) error {
	csSubmitterDID := newChar(submitterDID)
	defer freeChar(csSubmitterDID)

	csRequestJSON := newChar(requestJSON)
	defer freeChar(csRequestJSON)

	handle := callback.Register(cb)
	errCode := C.indy_sign_and_submit_request((C.indy_handle_t)(handle), (C.indy_handle_t)(poolHandle), (C.indy_handle_t)(walletHandle), csSubmitterDID, csRequestJSON, String())
	return indyerror.New(int32(errCode))
}

func SubmitRequest(poolHandle types.Handle, requestJSON string, cb callback.Callback) error {
	csRequestJSON := newChar(requestJSON)
	defer freeChar(csRequestJSON)

	handle := callback.Register(cb)
	errCode := C.indy_submit_request((C.indy_handle_t)(handle), (C.indy_handle_t)(poolHandle), csRequestJSON, String())
	return indyerror.New(int32(errCode))
}

func BuildSchemaRequest(submitterDID, data string, cb callback.Callback) error {
	csSubmitterDID := newChar(submitterDID)
	defer freeChar(csSubmitterDID)

	csData := newChar(data)
	defer freeChar(csData)

	handle := callback.Register(cb)
	errCode := C.indy_build_schema_request((C.indy_handle_t)(handle), csSubmitterDID, csData, String())
	return indyerror.New(int32(errCode))
}

func BuildGetSchemaRequest(submitterDID, id string, cb callback.Callback) error {
	csSubmitterDID := newChar(submitterDID)
	defer freeChar(csSubmitterDID)

	csID := newChar(id)
	defer freeChar(csID)

	handle := callback.Register(cb)
	errCode := C.indy_build_get_schema_request((C.indy_handle_t)(handle), csSubmitterDID, csID, String())
	return indyerror.New(int32(errCode))
}

func ParseGetSchemaResponse(response string, cb callback.Callback) error {
	csResponse := newChar(response)
	defer freeChar(csResponse)

	handle := callback.Register(cb)
	errCode := C.indy_parse_get_schema_response((C.indy_handle_t)(handle), csResponse, String2())
	return indyerror.New(int32(errCode))
}

func BuildCredDefRequest(submitterDID, data string, cb callback.Callback) error {
	csSubmitterDID := newChar(submitterDID)
	defer freeChar(csSubmitterDID)

	csData := newChar(data)
	defer freeChar(csData)

	handle := callback.Register(cb)
	errCode := C.indy_build_cred_def_request((C.indy_handle_t)(handle), csSubmitterDID, csData, String())
	return indyerror.New(int32(errCode))
}

func BuildGetCredDefRequest(submitterDID, id string, cb callback.Callback) error {
	csSubmitterDID := newChar(submitterDID)
	defer freeChar(csSubmitterDID)

	csID := newChar(id)
	defer freeChar(csID)

	handle := callback.Register(cb)
	errCode := C.indy_build_get_cred_def_request((C.indy_handle_t)(handle), csSubmitterDID, csID, String())
	return indyerror.New(int32(errCode))
}

func ParseGetCredDefResponse(response string, cb callback.Callback) error {
	csResponse := newChar(response)
	defer freeChar(csResponse)

	handle := callback.Register(cb)
	errCode := C.indy_parse_get_cred_def_response((C.indy_handle_t)(handle), csResponse, String2())
	return indyerror.New(int32(errCode))
}
