/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package anoncreds

import (
	"fmt"

	"github.com/hyperledger/indy-sdk-go/common/callback"
	"github.com/hyperledger/indy-sdk-go/common/logging"
	"github.com/hyperledger/indy-sdk-go/common/types"
	"github.com/hyperledger/indy-sdk-go/indy"
	"github.com/hyperledger/indy-sdk-go/wallet"
)

var logger = logging.MustGetLogger("indy-sdk")

// IssuerCreateSchema creates a credential schema entity that describes credential attributes list and allows credentials
// interoperability.
// Schema is public and intended to be shared with all anoncreds workflow actors usually by publishing SCHEMA transaction
// to Indy distributed ledger.
// It is IMPORTANT for current version POST Schema in Ledger and after that GET it from Ledger
// with correct seq_no to save compatibility with Ledger.
// After that can call indy_issuer_create_and_store_credential_def to build corresponding Credential Definition.
//
// issuerDid The DID of the issuer.
// name      Human-readable name of schema.
// version   Version of schema.
// attrs:    List of schema attributes descriptions
func IssuerCreateSchema(issuerDid, name, version, attrs string) (schemaID, schemaJSON string, err error) {
	respChan, errChan := issuerCreateSchema(issuerDid, name, version, attrs)
	select {
	case resp := <-respChan:
		schemaID = resp.schemaID
		schemaJSON = resp.schemaJSON
	case err = <-errChan:
	}
	return
}

// IssuerCreateAndStoreCredentialDef creates credential definition entity that encapsulates credentials issuer DID,
// credential schema, secrets used for signing credentials and secrets used for credentials revocation.
// Credential definition entity contains private and public parts. Private part will be stored in the wallet. Public part
// will be returned as json intended to be shared with all anoncreds workflow actors usually by publishing CRED_DEF transaction
// to Indy distributed ledger.
//
// It is IMPORTANT for current version GET Schema from Ledger with correct seq_no to save compatibility with Ledger.
//
// wallet     The wallet.
// issuerDid  DID of the issuer signing cred_def transaction to the Ledger
// schemaJson Сredential schema as a json
// tag        Allows to distinct between credential definitions for the same issuer and schema
// signature_type       Credential definition signature_type (optional, 'CL' by default) that defines credentials signature and revocation math.
//                   Supported types are:
//                   - 'CL': Camenisch-Lysyanskaya credential signature signature_type
// configJson Type-specific configuration of credential definition as json:
//                   - 'CL':
//                   - revocationSupport: whether to request non-revocation credential (optional, default false)
func IssuerCreateAndStoreCredentialDef(wallet *wallet.Wallet, issuerDID, schemaJSON, tag, signatureType, configJSON string) (credentialDefID, credentialDefJSON string, err error) {
	respChan, errChan := issuerCreateAndStoreCredentialDef(wallet, issuerDID, schemaJSON, tag, signatureType, configJSON)
	select {
	case resp := <-respChan:
		credentialDefID = resp.credentialDefID
		credentialDefJSON = resp.credentialDefJSON
	case err = <-errChan:
	}
	return
}

// IssuerCreateCredentialOffer creates credential offer that will be used by Prover for
// credential request creation. Offer includes nonce and key correctness proof
// for authentication between protocol steps and integrity checking.
//
// wallet    The wallet.
// credDefID ID of stored in ledger credential definition.
// return A JSON string containing the credential offer.
//     {
//         "schema_id": string,
//         "cred_def_id": string,
//         // Fields below can depend on Cred Def type
//         "nonce": string,
//         "key_correctness_proof" : <key_correctness_proof>
//     }
func IssuerCreateCredentialOffer(wallet *wallet.Wallet, credDefID string) (credentialOfferJSON string, err error) {
	respChan, errChan := issuerCreateCredentialOffer(wallet, credDefID)
	select {
	case credentialOfferJSON = <-respChan:
	case err = <-errChan:
	}
	return
}

// IssuerCreateCredential checks Cred Request for the given Cred Offer and issue Credential for the given Cred Request.
//
// Cred Request must match Cred Offer. The credential definition and revocation registry definition
// referenced in Cred Offer and Cred Request must be already created and stored into the wallet.
//
// Information for this credential revocation will be store in the wallet as part of revocation registry under
// generated cred_revoc_id local for this wallet.
//
// This call returns revoc registry delta as json file intended to be shared as REVOC_REG_ENTRY transaction.
// Note that it is possible to accumulate deltas to reduce ledger load.
//
// wallet                  The wallet.
// credOfferJSON           Cred offer created by issuerCreateCredentialOffer
// credReqJSON             Credential request created by proverCreateCredentialReq
// credValuesJSON          Credential containing attribute values for each of requested attribute names.
//                         Example:
//                         {
//                          "attr1" : {"raw": "value1", "encoded": "value1_as_int" },
//                          "attr2" : {"raw": "value1", "encoded": "value1_as_int" }
//                         }
// revRegId                (Optional) id of stored in ledger revocation registry definition
// blobStorageReaderHandle Pre-configured blob storage reader instance handle that will allow to read revocation tails
// return An IssuerCreateCredentialResult containing:
// credentialJSON: Credential json containing signed credential values
//  {
// 	  "schema_id": string,
// 	  "cred_def_id": string,
// 	  "rev_reg_def_id", Optional<string>,
// 	  "values": <see credValuesJson above>,
// 	  // Fields below can depend on Cred Def type
// 	  "signature": <signature>,
// 	  "signature_correctness_proof": <signature_correctness_proof>
//  }
// credRevocID: local id for revocation info (Can be used for revocation of this cred)
// revocRegDeltaJSON: Revocation registry delta json with a newly issued credential
func IssuerCreateCredential(wallet *wallet.Wallet, credOfferJSON, credReqJSON, credValuesJSON, revRegID string, blobStorageReaderHandle types.Handle) (credJSON, credRevID, revocRegDeltaJSON string, err error) {
	respChan, errChan := issuerCreateCredential(wallet, credOfferJSON, credReqJSON, credValuesJSON, revRegID, blobStorageReaderHandle)
	select {
	case resp := <-respChan:
		credJSON = resp.CredJSON
		credRevID = resp.CredRevID
		revocRegDeltaJSON = resp.RevocRegDeltaJSON
	case err = <-errChan:
	}
	return
}

// ProverCreateMasterSecret creates a master secret with a given name and stores it in the wallet.
//
// wallet         A wallet.
// masterSecretId (Optional, if not present random one will be generated) New master id
func ProverCreateMasterSecret(wallet *wallet.Wallet, secretID string) (masterSecretID string, err error) {
	respChan, errChan := proverCreateMasterSecret(wallet, secretID)
	select {
	case masterSecretID = <-respChan:
	case err = <-errChan:
	}
	return
}

// ProverCreateCredentialReq creates a claim request for the given credential offer.
//
// The method creates a blinded master secret for a master secret identified by a provided name.
// The master secret identified by the name must be already stored in the secure wallet (see proverCreateMasterSecret)
// The blinded master secret is a part of the credential request.
//
// wallet              A wallet.
// proverDID           The DID of the prover.
// credentialOfferJSON Credential offer as a json containing information about the issuer and a credential
// credentialDefJSON   Credential definition json
// masterSecretID      The ID of the master secret stored in the wallet
func ProverCreateCredentialReq(wallet *wallet.Wallet, proverDID, credentialOfferJSON, credentialDefJSON, masterSecretID string) (requestJSON string, requestMetadataJSON string, err error) {
	respChan, errChan := proverCreateCredentialReq(wallet, proverDID, credentialOfferJSON, credentialDefJSON, masterSecretID)
	select {
	case resp := <-respChan:
		requestJSON = resp.requestJSON
		requestMetadataJSON = resp.requestMetadataJSON
	case err = <-errChan:
	}
	return
}

// ProverStoreCredential checks credential provided by Issuer for the given credential request,
// updates the credential by a master secret and stores in a secure wallet.
//
// wallet              A Wallet.
// credId              (optional, default is a random one) Identifier by which credential will be stored in the wallet
// credReqMetadataJson Credential request metadata created by proverCreateCredentialReq
// credJson            Credential json received from issuer
// credDefJson         Credential definition json
// revRegDefJson       (optional) Revocation registry definition json
func ProverStoreCredential(wallet *wallet.Wallet, credID, credReqMetadataJSON, credJSON, credDefJSON, revRegDefJSON string) (responseJSON string, err error) {
	respChan, errChan := proverStoreCredential(wallet, credID, credReqMetadataJSON, credJSON, credDefJSON, revRegDefJSON)
	select {
	case responseJSON = <-respChan:
	case err = <-errChan:
	}
	return
}

// ProverGetCredentialsForProofReq gets human readable credentials matching the given proof request.
//
// wallet       A wallet.
// proofRequest proof request json
//     {
//         "name": string,
//         "version": string,
//         "nonce": string,
//         "requested_attributes": { // set of requested attributes
//              "<attr_referent>": <attr_info>, // see below
//              ...,
//         },
//         "requested_predicates": { // set of requested predicates
//              "<predicate_referent>": <predicate_info>, // see below
//              ...,
//          },
//         "non_revoked": Optional<<non_revoc_interval>>, // see below,
//                        // If specified prover must proof non-revocation
//                        // for date in this interval for each attribute
//                        // (can be overridden on attribute level)
//     }
//     where
//     attr_referent: Describes requested attribute
//     {
//         "name": string, // attribute name, (case insensitive and ignore spaces)
//         "restrictions": Optional<[<attr_filter>]> // see below,
//                          // if specified, credential must satisfy to one of the given restriction.
//         "non_revoked": Optional<<non_revoc_interval>>, // see below,
//                        // If specified prover must proof non-revocation
//                        // for date in this interval this attribute
//                        // (overrides proof level interval)
//     }
//     predicate_referent: Describes requested attribute predicate
//     {
//         "name": attribute name, (case insensitive and ignore spaces)
//         "p_type": predicate type (Currently >= only)
//         "p_value": predicate value
//         "restrictions": Optional<[<attr_filter>]> // see below,
//                         // if specified, credential must satisfy to one of the given restriction.
//         "non_revoked": Optional<<non_revoc_interval>>, // see below,
//                        // If specified prover must proof non-revocation
//                        // for date in this interval this attribute
//                        // (overrides proof level interval)
//     }
//     non_revoc_interval: Defines non-revocation interval
//     {
//         "from": Optional<int>, // timestamp of interval beginning
//         "to": Optional<int>, // timestamp of interval ending
//     }
//     filter: see filter above
//
// return A json with credentials for the given pool request.
//     {
//         "requested_attrs": {
//             "<attr_referent>": [{ cred_info: <credential_info>, interval: Optional<non_revoc_interval> }],
//             ...,
//         },
//         "requested_predicates": {
//             "requested_predicates": [{ cred_info: <credential_info>, timestamp: Optional<integer> }, { cred_info: <credential_2_info>, timestamp: Optional<integer> }],
//             "requested_predicate_2_referent": [{ cred_info: <credential_2_info>, timestamp: Optional<integer> }]
//         }
//     }, where credential is
//     {
//         "referent": <string>,
//         "attrs": [{"attr_name" : "attr_raw_value"}],
//         "schema_id": string,
//         "cred_def_id": string,
//         "rev_reg_id": Optional<int>,
//         "cred_rev_id": Optional<int>,
//     }
func ProverGetCredentialsForProofReq(wallet *wallet.Wallet, proofRequest string) (responseJSON string, err error) {
	respChan, errChan := proverGetCredentialsForProofReq(wallet, proofRequest)
	select {
	case responseJSON = <-respChan:
	case err = <-errChan:
	}
	return
}

// ProverCreateProof creates a proof according to the given proof request.
//
// wallet               A wallet.
// proofRequest proof request json
//     {
//         "name": string,
//         "version": string,
//         "nonce": string,
//         "requested_attributes": { // set of requested attributes
//              "<attr_referent>": <attr_info>, // see below
//              ...,
//         },
//         "requested_predicates": { // set of requested predicates
//              "<predicate_referent>": <predicate_info>, // see below
//              ...,
//          },
//         "non_revoked": Optional<<non_revoc_interval>>, // see below,
//                        // If specified prover must proof non-revocation
//                        // for date in this interval for each attribute
//                        // (can be overridden on attribute level)
//     }
// requestedCredentials either a credential or self-attested attribute for each requested attribute
//     {
//         "self_attested_attributes": {
//             "self_attested_attribute_referent": string
//         },
//         "requested_attributes": {
//             "requested_attribute_referent_1": {"cred_id": string, "timestamp": Optional<number>, revealed: <bool> }},
//             "requested_attribute_referent_2": {"cred_id": string, "timestamp": Optional<number>, revealed: <bool> }}
//         },
//         "requested_predicates": {
//             "requested_predicates_referent_1": {"cred_id": string, "timestamp": Optional<number> }},
//         }
//     }
// masterSecret         Id of the master secret stored in the wallet
// schemas              All schemas json participating in the proof request
//     {
//         <schema1_id>: <schema1_json>,
//         <schema2_id>: <schema2_json>,
//         <schema3_id>: <schema3_json>,
//     }
// credentialDefs       All credential definitions json participating in the proof request
//     {
//         "cred_def1_id": <credential_def1_json>,
//         "cred_def2_id": <credential_def2_json>,
//         "cred_def3_id": <credential_def3_json>,
//     }
// revStates            All revocation states json participating in the proof request
//     {
//         "rev_reg_def1_id": {
//             "timestamp1": <rev_state1>,
//             "timestamp2": <rev_state2>,
//         },
//         "rev_reg_def2_id": {
//             "timestamp3": <rev_state3>
//         },
//         "rev_reg_def3_id": {
//             "timestamp4": <rev_state4>
//         },
//     }
//
// return A a Proof json
// For each requested attribute either a proof (with optionally revealed attribute value) or
// self-attested attribute value is provided.
// Each proof is associated with a credential and corresponding schema_id, cred_def_id, rev_reg_id and timestamp.
// There is also aggregated proof part common for all credential proofs.
//     {
//         "requested": {
//             "revealed_attrs": {
//                 "requested_attr1_id": {sub_proof_index: number, raw: string, encoded: string},
//                 "requested_attr4_id": {sub_proof_index: number: string, encoded: string},
//             },
//             "unrevealed_attrs": {
//                 "requested_attr3_id": {sub_proof_index: number}
//             },
//             "self_attested_attrs": {
//                 "requested_attr2_id": self_attested_value,
//             },
//             "requested_predicates": {
//                 "requested_predicate_1_referent": {sub_proof_index: int},
//                 "requested_predicate_2_referent": {sub_proof_index: int},
//             }
//         }
//         "proof": {
//             "proofs": [ <credential_proof>, <credential_proof>, <credential_proof> ],
//             "aggregated_proof": <aggregated_proof>
//         }
//         "identifiers": [{schema_id, cred_def_id, Optional<rev_reg_id>, Optional<timestamp>}]
//     }
func ProverCreateProof(wallet *wallet.Wallet, proofRequest, requestedCredentials, masterSecret, schemas, credentialDefs, revStates string) (proofJSON string, err error) {
	respChan, errChan := proverCreateProof(wallet, proofRequest, requestedCredentials, masterSecret, schemas, credentialDefs, revStates)
	select {
	case proofJSON = <-respChan:
	case err = <-errChan:
	}
	return
}

// VerifierVerifyProof verifies a proof (of multiple credential).
// All required schemas, public keys and revocation registries must be provided.
//
// proofRequest   proof request json
//     {
//         "name": string,
//         "version": string,
//         "nonce": string,
//         "requested_attributes": { // set of requested attributes
//              "<attr_referent>": <attr_info>, // see below
//              ...,
//         },
//         "requested_predicates": { // set of requested predicates
//              "<predicate_referent>": <predicate_info>, // see below
//              ...,
//          },
//         "non_revoked": Optional<<non_revoc_interval>>, // see below,
//                        // If specified prover must proof non-revocation
//                        // for date in this interval for each attribute
//                        // (can be overridden on attribute level)
//     }
//
// proof          Proof json
//     {
//         "requested": {
//             "revealed_attrs": {
//                 "requested_attr1_id": {sub_proof_index: number, raw: string, encoded: string},
//                 "requested_attr4_id": {sub_proof_index: number: string, encoded: string},
//             },
//             "unrevealed_attrs": {
//                 "requested_attr3_id": {sub_proof_index: number}
//             },
//             "self_attested_attrs": {
//                 "requested_attr2_id": self_attested_value,
//             },
//             "requested_predicates": {
//                 "requested_predicate_1_referent": {sub_proof_index: int},
//                 "requested_predicate_2_referent": {sub_proof_index: int},
//             }
//         }
//         "proof": {
//             "proofs": [ <credential_proof>, <credential_proof>, <credential_proof> ],
//             "aggregated_proof": <aggregated_proof>
//         }
//         "identifiers": [{schema_id, cred_def_id, Optional<rev_reg_id>, Optional<timestamp>}]
//     }
//
// schemas        All schemas json participating in the proof request
//     {
//         <schema1_id>: <schema1_json>,
//         <schema2_id>: <schema2_json>,
//         <schema3_id>: <schema3_json>,
//     }
//
// credentialDefs  All credential definitions json participating in the proof request
//     {
//         "cred_def1_id": <credential_def1_json>,
//         "cred_def2_id": <credential_def2_json>,
//         "cred_def3_id": <credential_def3_json>,
//     }
//
// revocRegDefs   All revocation registry definitions json participating in the proof
//     {
//         "rev_reg_def1_id": <rev_reg_def1_json>,
//         "rev_reg_def2_id": <rev_reg_def2_json>,
//         "rev_reg_def3_id": <rev_reg_def3_json>,
//     }
//
// revocRegs      all revocation registries json participating in the proof
//     {
//         "rev_reg_def1_id": {
//             "timestamp1": <rev_reg1>,
//             "timestamp2": <rev_reg2>,
//         },
//         "rev_reg_def2_id": {
//             "timestamp3": <rev_reg3>
//         },
//         "rev_reg_def3_id": {
//             "timestamp4": <rev_reg4>
//         },
//     }
//
// return true if signature is valid, otherwise false
func VerifierVerifyProof(proofRequest, proof, schemas, credentialDefs, revocRegDefs, revocRegs string) (valid bool, err error) {
	respChan, errChan := verifierVerifyProof(proofRequest, proof, schemas, credentialDefs, revocRegDefs, revocRegs)
	select {
	case valid = <-respChan:
	case err = <-errChan:
	}
	return
}

type issuerCreateSchemaResponse struct {
	schemaID   string
	schemaJSON string
}

func issuerCreateSchema(issuerDID, name, version, attrs string) (chan *issuerCreateSchemaResponse, chan error) {
	logger.Debugf("Creating issuer schema - IssuerDID: [%s], Name: [%s], Version: [%s], Attrs: [%s]", issuerDID, name, version, attrs)

	respChan := make(chan *issuerCreateSchemaResponse)
	errChan := make(chan error, 1)

	if issuerDID == "" {
		errChan <- fmt.Errorf("issuer DID must be specified")
		return respChan, errChan
	}
	if name == "" {
		errChan <- fmt.Errorf("name must be specified")
		return respChan, errChan
	}
	if version == "" {
		errChan <- fmt.Errorf("version must be specified")
		return respChan, errChan
	}
	if attrs == "" {
		errChan <- fmt.Errorf("attrs must be specified")
		return respChan, errChan
	}

	cb := func(err error, data callback.Data) {
		if err != nil {
			errChan <- err
		} else {
			resp := data.([]string)
			respChan <- &issuerCreateSchemaResponse{
				schemaID:   resp[0],
				schemaJSON: resp[1],
			}
		}
	}

	err := indy.IssuerCreateSchema(issuerDID, name, version, attrs, cb)
	if err != nil {
		// Send the error immediately
		errChan <- err
	}

	return respChan, errChan
}

type issuerCAndSCredDefResp struct {
	credentialDefID   string
	credentialDefJSON string
}

func issuerCreateAndStoreCredentialDef(wallet *wallet.Wallet, issuerDID, schemaJSON, tag, signatureType, configJSON string) (chan *issuerCAndSCredDefResp, chan error) {
	logger.Debugf("Creating and storing credential def - Wallet: [%s], IssuerDid: [%s], schemaJSON: %s, tag: [%s], signatureType: [%s], configJSON: %s", wallet.Name, issuerDID, schemaJSON, tag, signatureType, configJSON)

	respChan := make(chan *issuerCAndSCredDefResp)
	errChan := make(chan error, 1)

	if issuerDID == "" {
		errChan <- fmt.Errorf("issuer DID must be specified")
		return respChan, errChan
	}
	if schemaJSON == "" {
		errChan <- fmt.Errorf("schema must be specified")
		return respChan, errChan
	}
	if tag == "" {
		errChan <- fmt.Errorf("tag must be specified")
		return respChan, errChan
	}
	if configJSON == "" {
		errChan <- fmt.Errorf("config must be specified")
		return respChan, errChan
	}

	cb := func(err error, data callback.Data) {
		if err != nil {
			errChan <- err
		} else {
			resp := data.([]string)
			respChan <- &issuerCAndSCredDefResp{
				credentialDefID:   resp[0],
				credentialDefJSON: resp[1],
			}
		}
	}

	err := indy.IssuerCreateAndStoreCredentialDef(wallet.Handle(), issuerDID, schemaJSON, tag, signatureType, configJSON, cb)
	if err != nil {
		// Send the error immediately
		errChan <- err
	}

	return respChan, errChan
}

func issuerCreateCredentialOffer(wallet *wallet.Wallet, credDefID string) (chan string, chan error) {
	logger.Debugf("Creating credential offer - Wallet: [%s], credDefID: [%s]", wallet.Name, credDefID)

	respChan := make(chan string)
	errChan := make(chan error, 1)

	if credDefID == "" {
		errChan <- fmt.Errorf("credential def ID must be specified")
		return respChan, errChan
	}

	cb := func(err error, data callback.Data) {
		if err != nil {
			errChan <- err
		} else {
			respChan <- data.(string)
		}
	}

	err := indy.IssuerCreateCredentialOffer(wallet.Handle(), credDefID, cb)
	if err != nil {
		// Send the error immediately
		errChan <- err
	}

	return respChan, errChan
}

type issuerCreateCredentialResp struct {
	CredJSON          string
	CredRevID         string
	RevocRegDeltaJSON string
}

func issuerCreateCredential(wallet *wallet.Wallet, credOfferJSON, credReqJSON, credValuesJSON, revRegID string, blobStorageReaderHandle types.Handle) (chan *issuerCreateCredentialResp, chan error) {
	logger.Debugf("Creating credential - Wallet: [%s], credOfferJSON: [%s], credReqJSON: %s, credValuesJSON: %s, revRegID: [%s]", wallet.Name, credOfferJSON, credReqJSON, credValuesJSON, revRegID)

	respChan := make(chan *issuerCreateCredentialResp)
	errChan := make(chan error, 1)

	if credOfferJSON == "" {
		errChan <- fmt.Errorf("cred offer JSON must be specified")
		return respChan, errChan
	}
	if credReqJSON == "" {
		errChan <- fmt.Errorf("cred request JSON must be specified")
		return respChan, errChan
	}
	if credValuesJSON == "" {
		errChan <- fmt.Errorf("cred values JSON must be specified")
		return respChan, errChan
	}

	cb := func(err error, data callback.Data) {
		if err != nil {
			errChan <- err
		} else {
			sd := data.([]string)
			respChan <- &issuerCreateCredentialResp{
				CredJSON:          sd[0],
				CredRevID:         sd[1],
				RevocRegDeltaJSON: sd[2],
			}
		}
	}

	err := indy.IssuerCreateCredential(wallet.Handle(), credOfferJSON, credReqJSON, credValuesJSON, revRegID, blobStorageReaderHandle, cb)
	if err != nil {
		// Send the error immediately
		errChan <- err
	}

	return respChan, errChan
}

func proverCreateMasterSecret(wallet *wallet.Wallet, masterSecretID string) (chan string, chan error) {
	logger.Debugf("Creating master secret - Wallet: [%s], MasterSecretID: [%s]", wallet.Name, masterSecretID)

	respChan := make(chan string)
	errChan := make(chan error, 1)

	cb := func(err error, data callback.Data) {
		if err != nil {
			errChan <- err
		} else {
			respChan <- data.(string)
		}
	}

	err := indy.ProverCreateMasterSecret(wallet.Handle(), masterSecretID, cb)
	if err != nil {
		// Send the error immediately
		errChan <- err
	}

	return respChan, errChan
}

type proverCreateCredentialResp struct {
	requestJSON         string
	requestMetadataJSON string
}

func proverCreateCredentialReq(wallet *wallet.Wallet, proverDID, credentialOfferJSON, credentialDefJSON, masterSecretID string) (chan *proverCreateCredentialResp, chan error) {
	logger.Debugf("Creating credential request - Wallet: [%s], proverDID: [%s], credentialOfferJSON: %s, credentialDefJSON: %s, MasterSecretID: [%s]", wallet.Name, proverDID, credentialOfferJSON, credentialDefJSON, masterSecretID)

	respChan := make(chan *proverCreateCredentialResp)
	errChan := make(chan error, 1)

	if proverDID == "" {
		errChan <- fmt.Errorf("prover DID must be specified")
		return respChan, errChan
	}
	if credentialOfferJSON == "" {
		errChan <- fmt.Errorf("credential offer JSON must be specified")
		return respChan, errChan
	}
	if credentialDefJSON == "" {
		errChan <- fmt.Errorf("credential def JSON must be specified")
		return respChan, errChan
	}
	if masterSecretID == "" {
		errChan <- fmt.Errorf("master secret ID must be specified")
		return respChan, errChan
	}

	cb := func(err error, data callback.Data) {
		if err != nil {
			errChan <- err
		} else {
			sd := data.([]string)
			respChan <- &proverCreateCredentialResp{
				requestJSON:         sd[0],
				requestMetadataJSON: sd[1],
			}
		}
	}

	err := indy.ProverCreateCredentialReq(wallet.Handle(), proverDID, credentialOfferJSON, credentialDefJSON, masterSecretID, cb)
	if err != nil {
		// Send the error immediately
		errChan <- err
	}

	return respChan, errChan
}

func proverStoreCredential(wallet *wallet.Wallet, credID, credReqMetadataJSON, credJSON, credDefJSON, revRegDefJSON string) (chan string, chan error) {
	logger.Debugf("Storing credential - Wallet: [%s], credID: [%s], credReqMetadataJSON: %s, credJSON: %s, credDefJSON: %s, revRegDefJSON: %s", wallet.Name, credID, credReqMetadataJSON, credJSON, credDefJSON, revRegDefJSON)

	respChan := make(chan string)
	errChan := make(chan error, 1)

	if credReqMetadataJSON == "" {
		errChan <- fmt.Errorf("cred request metadata JSON must be specified")
		return respChan, errChan
	}
	if credJSON == "" {
		errChan <- fmt.Errorf("cred JSON must be specified")
		return respChan, errChan
	}
	if credDefJSON == "" {
		errChan <- fmt.Errorf("cred def JSON must be specified")
		return respChan, errChan
	}

	cb := func(err error, data callback.Data) {
		if err != nil {
			errChan <- err
		} else {
			respChan <- data.(string)
		}
	}

	err := indy.ProverStoreCredential(wallet.Handle(), credID, credReqMetadataJSON, credJSON, credDefJSON, revRegDefJSON, cb)
	if err != nil {
		// Send the error immediately
		errChan <- err
	}

	return respChan, errChan
}

func proverGetCredentialsForProofReq(wallet *wallet.Wallet, proofRequest string) (chan string, chan error) {
	logger.Debugf("Storing credential - Wallet: [%s], proofRequest: [%s]", wallet.Name, proofRequest)

	respChan := make(chan string)
	errChan := make(chan error, 1)

	if proofRequest == "" {
		errChan <- fmt.Errorf("proof request must be specified")
		return respChan, errChan
	}

	cb := func(err error, data callback.Data) {
		if err != nil {
			errChan <- err
		} else {
			respChan <- data.(string)
		}
	}

	err := indy.ProverGetCredentialsForProofReq(wallet.Handle(), proofRequest, cb)
	if err != nil {
		// Send the error immediately
		errChan <- err
	}

	return respChan, errChan
}

func proverCreateProof(wallet *wallet.Wallet, proofRequest, requestedCredentials, masterSecret, schemas, credentialDefs, revStates string) (chan string, chan error) {
	logger.Debugf("Storing credential - Wallet: [%s], proofRequest: [%s], requestedCredentials: [%s], masterSecret: [%s], schemas: [%s], credentialDefs: [%s], revStates: [%s]", wallet.Name, proofRequest, requestedCredentials, masterSecret, schemas, credentialDefs, revStates)

	respChan := make(chan string)
	errChan := make(chan error, 1)

	if proofRequest == "" {
		errChan <- fmt.Errorf("proof request must be specified")
		return respChan, errChan
	}
	if requestedCredentials == "" {
		errChan <- fmt.Errorf("requested credentials must be specified")
		return respChan, errChan
	}
	if masterSecret == "" {
		errChan <- fmt.Errorf("master secret must be specified")
		return respChan, errChan
	}
	if schemas == "" {
		errChan <- fmt.Errorf("schemas must be specified")
		return respChan, errChan
	}
	if credentialDefs == "" {
		errChan <- fmt.Errorf("credential defs must be specified")
		return respChan, errChan
	}
	if revStates == "" {
		errChan <- fmt.Errorf("rev states must be specified")
		return respChan, errChan
	}

	cb := func(err error, data callback.Data) {
		if err != nil {
			errChan <- err
		} else {
			respChan <- data.(string)
		}
	}

	err := indy.ProverCreateProof(wallet.Handle(), proofRequest, requestedCredentials, masterSecret, schemas, credentialDefs, revStates, cb)
	if err != nil {
		// Send the error immediately
		errChan <- err
	}

	return respChan, errChan
}

func verifierVerifyProof(proofRequest, proof, schemas, credentialDefs, revocRegDefs, revocRegs string) (chan bool, chan error) {
	logger.Debugf("Storing credential - proofRequest: [%s], proof: [%s], schemas: [%s], credentialDefs: [%s], revocRegDefs: [%s], revocRegs: [%s]", proofRequest, proof, schemas, credentialDefs, revocRegDefs, revocRegs)

	respChan := make(chan bool)
	errChan := make(chan error, 1)

	if proofRequest == "" {
		errChan <- fmt.Errorf("proof request must be specified")
		return respChan, errChan
	}
	if proof == "" {
		errChan <- fmt.Errorf("proof must be specified")
		return respChan, errChan
	}
	if schemas == "" {
		errChan <- fmt.Errorf("schemas must be specified")
		return respChan, errChan
	}
	if credentialDefs == "" {
		errChan <- fmt.Errorf("credential defs must be specified")
		return respChan, errChan
	}
	if revocRegDefs == "" {
		errChan <- fmt.Errorf("revoc reg defs must be specified")
		return respChan, errChan
	}
	if revocRegs == "" {
		errChan <- fmt.Errorf("revoc regs must be specified")
		return respChan, errChan
	}

	cb := func(err error, data callback.Data) {
		if err != nil {
			errChan <- err
		} else {
			respChan <- data.(bool)
		}
	}

	err := indy.VerifierVerifyProof(proofRequest, proof, schemas, credentialDefs, revocRegDefs, revocRegs, cb)
	if err != nil {
		// Send the error immediately
		errChan <- err
	}

	return respChan, errChan
}
