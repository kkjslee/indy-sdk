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
#include <indy_mod.h>
#include <indy_types.h>
#include <indy_anoncreds.h>
*/
import "C"

func IssuerCreateSchema(issuerDID, name, version, attrs string, cb callback.Callback) error {
	csIssuerDID := newChar(issuerDID)
	defer freeChar(csIssuerDID)

	csName := newChar(name)
	defer freeChar(csName)

	csVersion := newChar(version)
	defer freeChar(csVersion)

	csAttrs := newChar(attrs)
	defer freeChar(csAttrs)

	handle := callback.Register(cb)
	errCode := C.indy_issuer_create_schema((C.indy_handle_t)(handle), csIssuerDID, csName, csVersion, csAttrs, String2())
	return indyerror.New(int32(errCode))
}

func IssuerCreateAndStoreCredentialDef(walletHandle types.Handle, issuerDID, schemaJSON, tag, signatureType, configJSON string, cb callback.Callback) error {
	csIssuerDID := newChar(issuerDID)
	defer freeChar(csIssuerDID)

	csSchemaJSON := newChar(schemaJSON)
	defer freeChar(csSchemaJSON)

	csTag := newChar(tag)
	defer freeChar(csTag)

	csSignatureType := newChar(signatureType)
	defer freeChar(csSignatureType)

	csConfigJSON := newChar(configJSON)
	defer freeChar(csConfigJSON)

	handle := callback.Register(cb)
	errCode := C.indy_issuer_create_and_store_credential_def((C.indy_handle_t)(handle), (C.indy_handle_t)(walletHandle), csIssuerDID, csSchemaJSON, csTag, csSignatureType, csConfigJSON, String2())
	return indyerror.New(int32(errCode))
}

func IssuerCreateCredentialOffer(walletHandle types.Handle, credDefID string, cb callback.Callback) error {
	csCredDefID := newChar(credDefID)
	defer freeChar(csCredDefID)

	handle := callback.Register(cb)
	errCode := C.indy_issuer_create_credential_offer((C.indy_handle_t)(handle), (C.indy_handle_t)(walletHandle), csCredDefID, String())
	return indyerror.New(int32(errCode))
}

func IssuerCreateCredential(walletHandle types.Handle, credOfferJSON, credReqJSON, credValuesJSON, revRegID string, blobStorageReaderHandle types.Handle, cb callback.Callback) error {
	csCredOfferJSON := newChar(credOfferJSON)
	defer freeChar(csCredOfferJSON)

	csCredReqJSON := newChar(credReqJSON)
	defer freeChar(csCredReqJSON)

	csCredValuesJSON := newChar(credValuesJSON)
	defer freeChar(csCredValuesJSON)

	var csRevRegID *C.char
	if revRegID != "" {
		csRevRegID = newChar(revRegID)
		defer freeChar(csRevRegID)
	}

	handle := callback.Register(cb)
	errCode := C.indy_issuer_create_credential((C.indy_handle_t)(handle), (C.indy_handle_t)(walletHandle), csCredOfferJSON, csCredReqJSON, csCredValuesJSON, csRevRegID, (C.indy_i32_t)(blobStorageReaderHandle), String3())
	return indyerror.New(int32(errCode))
}

func ProverCreateMasterSecret(walletHandle types.Handle, masterSecretID string, cb callback.Callback) error {
	var csMasterSecretID *C.char
	if masterSecretID != "" {
		csMasterSecretID = newChar(masterSecretID)
		defer freeChar(csMasterSecretID)
	}

	handle := callback.Register(cb)
	errCode := C.indy_prover_create_master_secret((C.indy_handle_t)(handle), (C.indy_handle_t)(walletHandle), csMasterSecretID, String())
	return indyerror.New(int32(errCode))
}

func ProverCreateCredentialReq(walletHandle types.Handle, proverDID, credentialOfferJSON, credentialDefJSON, masterSecretID string, cb callback.Callback) error {
	csProverDID := newChar(proverDID)
	defer freeChar(csProverDID)

	csCredentialOfferJSON := newChar(credentialOfferJSON)
	defer freeChar(csCredentialOfferJSON)

	csCredentialDefJSON := newChar(credentialDefJSON)
	defer freeChar(csCredentialDefJSON)

	csMasterSecretID := newChar(masterSecretID)
	defer freeChar(csMasterSecretID)

	handle := callback.Register(cb)
	errCode := C.indy_prover_create_credential_req((C.indy_handle_t)(handle), (C.indy_handle_t)(walletHandle), csProverDID, csCredentialOfferJSON, csCredentialDefJSON, csMasterSecretID, String2())
	return indyerror.New(int32(errCode))
}

func ProverStoreCredential(walletHandle types.Handle, credID, credReqMetadataJSON, credJSON, credDefJSON, revRegDefJSON string, cb callback.Callback) error {
	var csCredID *C.char
	if credID != "" {
		csCredID = newChar(credID)
		defer freeChar(csCredID)
	}

	csCredReqMetadataJSON := newChar(credReqMetadataJSON)
	defer freeChar(csCredReqMetadataJSON)

	csCredJSON := newChar(credJSON)
	defer freeChar(csCredJSON)

	csCredDefJSON := newChar(credDefJSON)
	defer freeChar(csCredDefJSON)

	var csRevRegDefJSON *C.char
	if revRegDefJSON != "" {
		csRevRegDefJSON = newChar(revRegDefJSON)
		defer freeChar(csRevRegDefJSON)
	}

	handle := callback.Register(cb)
	errCode := C.indy_prover_store_credential((C.indy_handle_t)(handle), (C.indy_handle_t)(walletHandle), csCredID, csCredReqMetadataJSON, csCredJSON, csCredDefJSON, csRevRegDefJSON, String())
	return indyerror.New(int32(errCode))
}

func ProverGetCredentialsForProofReq(walletHandle types.Handle, proofRequest string, cb callback.Callback) error {
	csProofRequest := newChar(proofRequest)
	defer freeChar(csProofRequest)

	handle := callback.Register(cb)
	errCode := C.indy_prover_get_credentials_for_proof_req((C.indy_handle_t)(handle), (C.indy_handle_t)(walletHandle), csProofRequest, String())
	return indyerror.New(int32(errCode))
}

func ProverCreateProof(walletHandle types.Handle, proofRequest, requestedCredentials, masterSecret, schemas, credentialDefs, revStates string, cb callback.Callback) error {
	csProofRequest := newChar(proofRequest)
	defer freeChar(csProofRequest)

	csRequestedCredentials := newChar(requestedCredentials)
	defer freeChar(csRequestedCredentials)

	csMasterSecret := newChar(masterSecret)
	defer freeChar(csMasterSecret)

	csSchemas := newChar(schemas)
	defer freeChar(csSchemas)

	csCredentialDefs := newChar(credentialDefs)
	defer freeChar(csCredentialDefs)

	csRevStates := newChar(revStates)
	defer freeChar(csRevStates)

	handle := callback.Register(cb)
	errCode := C.indy_prover_create_proof((C.indy_handle_t)(handle), (C.indy_handle_t)(walletHandle), csProofRequest, csRequestedCredentials, csMasterSecret, csSchemas, csCredentialDefs, csRevStates, String())
	return indyerror.New(int32(errCode))
}

func VerifierVerifyProof(proofRequest, proof, schemas, credentialDefs, revocRegDefs, revocRegs string, cb callback.Callback) error {
	csProofRequest := newChar(proofRequest)
	defer freeChar(csProofRequest)

	csProof := newChar(proof)
	defer freeChar(csProof)

	csSchemas := newChar(schemas)
	defer freeChar(csSchemas)

	csCredentialDefs := newChar(credentialDefs)
	defer freeChar(csCredentialDefs)

	csRevocRegDefs := newChar(revocRegDefs)
	defer freeChar(csRevocRegDefs)

	csRevocRegs := newChar(revocRegs)
	defer freeChar(csRevocRegs)

	handle := callback.Register(cb)
	errCode := C.indy_verifier_verify_proof((C.indy_handle_t)(handle), csProofRequest, csProof, csSchemas, csCredentialDefs, csRevocRegDefs, csRevocRegs, Bool())
	return indyerror.New(int32(errCode))
}
