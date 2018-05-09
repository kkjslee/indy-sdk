/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package indyerror

import "fmt"

const (
	// Success operation succeeded
	Success = 0

	/* --- Common errors --- */

	// CommonInvalidParam1 Caller passed invalid value as param 1 (null, invalid json and etc..)
	CommonInvalidParam1 = 100

	// CommonInvalidParam2 Caller passed invalid value as param 2 (null, invalid json and etc..)
	CommonInvalidParam2 = 101

	// CommonInvalidParam3 Caller passed invalid value as param 3 (null, invalid json and etc..)
	CommonInvalidParam3 = 102

	// CommonInvalidParam4 Caller passed invalid value as param 4 (null, invalid json and etc..)
	CommonInvalidParam4 = 103

	// CommonInvalidParam5 Caller passed invalid value as param 5 (null, invalid json and etc..)
	CommonInvalidParam5 = 104

	// CommonInvalidParam6 Caller passed invalid value as param 6 (null, invalid json and etc..)
	CommonInvalidParam6 = 105

	// CommonInvalidParam7 Caller passed invalid value as param 7 (null, invalid json and etc..)
	CommonInvalidParam7 = 106

	// CommonInvalidParam8 Caller passed invalid value as param 8 (null, invalid json and etc..)
	CommonInvalidParam8 = 107

	// CommonInvalidParam9 Caller passed invalid value as param 9 (null, invalid json and etc..)
	CommonInvalidParam9 = 108

	// CommonInvalidParam10 Caller passed invalid value as param 10 (null, invalid json and etc..)
	CommonInvalidParam10 = 109

	// CommonInvalidParam11 Caller passed invalid value as param 11 (null, invalid json and etc..)
	CommonInvalidParam11 = 110

	// CommonInvalidParam12 Caller passed invalid value as param 12 (null, invalid json and etc..)
	CommonInvalidParam12 = 111

	// CommonInvalidState Invalid library state was detected in runtime. It signals library bug
	CommonInvalidState = 112

	// CommonInvalidStructure Object (json, config, key, credential and etc...) passed by library caller has invalid structure
	CommonInvalidStructure = 113

	// CommonIOError IO Error
	CommonIOError = 114

	/* --- Wallet errors --- */

	// WalletInvalidHandle Caller passed invalid wallet handle
	WalletInvalidHandle = 200

	// WalletUnknownTypeError Unknown type of wallet was passed on create_wallet
	WalletUnknownTypeError = 201

	// WalletTypeAlreadyRegisteredError Attempt to register already existing wallet type
	WalletTypeAlreadyRegisteredError = 202

	// WalletAlreadyExistsError Attempt to create wallet with name used for another exists wallet
	WalletAlreadyExistsError = 203

	// WalletNotFoundError Requested entity id isn't present in wallet
	WalletNotFoundError = 204

	// WalletIncompatiblePoolError Trying to use wallet with pool that has different name
	WalletIncompatiblePoolError = 205

	// WalletAlreadyOpenedError Trying to open wallet that was opened already
	WalletAlreadyOpenedError = 206

	// WalletAccessFailed Attempt to open encrypted wallet with invalid credentials
	WalletAccessFailed = 207

	/* --- Ledger errors --- */

	// PoolLedgerNotCreatedError Trying to open pool ledger that wasn't created before
	PoolLedgerNotCreatedError = 300

	// PoolLedgerInvalidPoolHandle Caller passed invalid pool ledger handle
	PoolLedgerInvalidPoolHandle = 301

	// PoolLedgerTerminated Pool ledger terminated
	PoolLedgerTerminated = 302

	// LedgerNoConsensusError No concensus during ledger operation
	LedgerNoConsensusError = 303

	// LedgerInvalidTransaction Attempt to parse invalid transaction response
	LedgerInvalidTransaction = 304

	// LedgerSecurityError Attempt to send transaction without the necessary privileges
	LedgerSecurityError = 305

	// PoolLedgerConfigAlreadyExistsError Attempt to create pool ledger config with name used for another existing pool
	PoolLedgerConfigAlreadyExistsError = 306

	// PoolLedgerTimeout Timeout for action
	PoolLedgerTimeout = 307

	// AnoncredsRevocationRegistryFullError Revocation registry is full and creation of new registry is necessary
	AnoncredsRevocationRegistryFullError = 400

	// AnoncredsInvalidUserRevocID Invalid user revocation ID
	AnoncredsInvalidUserRevocID = 401

	// AnoncredsMasterSecretDuplicateNameError Attempt to generate master secret with dupplicated name
	AnoncredsMasterSecretDuplicateNameError = 404

	// AnoncredsProofRejected Proof was rejected
	AnoncredsProofRejected = 405

	// AnoncredsCredentialRevoked Credentials were revoked
	AnoncredsCredentialRevoked = 406

	// AnoncredsCredDefAlreadyExistsError Attempt to create credential definition with duplicated did schema pair
	AnoncredsCredDefAlreadyExistsError = 407

	/* --- Crypto errors --- */

	// UnknownCryptoTypeError Unknown format of DID entity keys
	UnknownCryptoTypeError = 500

	// DidAlreadyExistsError Attempt to create duplicate did
	DidAlreadyExistsError = 600

	// Undefined indicates that the error is not an Indy error
	Undefined = -1
)

// IndyError extends error and adds a code
type IndyError interface {
	error
	Code() int32
}

var errorStringMap = map[int32]string{
	Success:                                 "Success",
	CommonInvalidParam1:                     "Caller passed invalid value as param 1 (null, invalid json and etc..)",
	CommonInvalidParam2:                     "Caller passed invalid value as param 2 (null, invalid json and etc..)",
	CommonInvalidParam3:                     "Caller passed invalid value as param 3 (null, invalid json and etc..)",
	CommonInvalidParam4:                     "Caller passed invalid value as param 4 (null, invalid json and etc..)",
	CommonInvalidParam5:                     "Caller passed invalid value as param 5 (null, invalid json and etc..)",
	CommonInvalidParam6:                     "Caller passed invalid value as param 6 (null, invalid json and etc..)",
	CommonInvalidParam7:                     "Caller passed invalid value as param 7 (null, invalid json and etc..)",
	CommonInvalidParam8:                     "Caller passed invalid value as param 8 (null, invalid json and etc..)",
	CommonInvalidParam9:                     "Caller passed invalid value as param 9 (null, invalid json and etc..)",
	CommonInvalidParam10:                    "Caller passed invalid value as param 10 (null, invalid json and etc..)",
	CommonInvalidParam11:                    "Caller passed invalid value as param 11 (null, invalid json and etc..)",
	CommonInvalidParam12:                    "Caller passed invalid value as param 12 (null, invalid json and etc..)",
	CommonInvalidState:                      "Invalid library state was detected in runtime. It signals library bug",
	CommonInvalidStructure:                  "Object (json, config, key, credential and etc...) passed by library caller has invalid structure",
	CommonIOError:                           "IO Error",
	WalletInvalidHandle:                     "Caller passed invalid wallet handle",
	WalletUnknownTypeError:                  "Unknown type of wallet was passed on create_wallet",
	WalletTypeAlreadyRegisteredError:        "Attempt to register already existing wallet type",
	WalletAlreadyExistsError:                "Attempt to create wallet with name used for another exists wallet",
	WalletNotFoundError:                     "Requested entity id isn't present in wallet",
	WalletIncompatiblePoolError:             "Trying to use wallet with pool that has different name",
	WalletAlreadyOpenedError:                "Trying to open wallet that was opened already",
	WalletAccessFailed:                      "Attempt to open encrypted wallet with invalid credentials",
	PoolLedgerNotCreatedError:               "Trying to open pool ledger that wasn't created before",
	PoolLedgerInvalidPoolHandle:             "Caller passed invalid pool ledger handle",
	PoolLedgerTerminated:                    "Pool ledger terminated",
	LedgerNoConsensusError:                  "No concensus during ledger operation",
	LedgerInvalidTransaction:                "Attempt to parse invalid transaction response",
	LedgerSecurityError:                     "Attempt to send transaction without the necessary privileges",
	PoolLedgerConfigAlreadyExistsError:      "Attempt to create pool ledger config with name used for another existing pool",
	PoolLedgerTimeout:                       "Timeout for action",
	AnoncredsRevocationRegistryFullError:    "Revocation registry is full and creation of new registry is necessary",
	AnoncredsInvalidUserRevocID:             "AnoncredsInvalidUserRevocId",
	AnoncredsMasterSecretDuplicateNameError: "Attempt to generate master secret with dupplicated name",
	AnoncredsProofRejected:                  "AnoncredsProofRejected",
	AnoncredsCredentialRevoked:              "AnoncredsCredentialRevoked",
	AnoncredsCredDefAlreadyExistsError:      "Attempt to create credential definition with duplicated did schema pair",
	UnknownCryptoTypeError:                  "Unknown format of DID entity keys",
	DidAlreadyExistsError:                   "Attempt to create duplicate did",
}

// New returns a new IndError
func New(code int32) IndyError {
	if code == Success {
		return nil
	}

	var err error
	if msg, ok := errorStringMap[code]; ok {
		err = fmt.Errorf(msg)
	} else {
		err = fmt.Errorf("unknown error: %d", code)
	}

	return &indyError{
		error: err,
		code:  code,
	}
}

// Code returns the error code from the given error.
// If the error doesn't have a code then Undefined (-1) is returned.
func Code(err error) int32 {
	indyErr, ok := err.(IndyError)
	if !ok {
		return Undefined
	}
	return indyErr.Code()
}

type indyError struct {
	error
	code int32
}

// Code returns the Indy error code
func (e *indyError) Code() int32 {
	return e.code
}
