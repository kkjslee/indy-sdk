/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ledger

import (
	"fmt"

	"github.com/hyperledger/indy-sdk-go/common/role"
	"github.com/hyperledger/indy-sdk-go/common/types"
	"github.com/hyperledger/indy-sdk-go/pool"
	"github.com/hyperledger/indy-sdk-go/wallet"

	"github.com/hyperledger/indy-sdk-go/common/callback"
	"github.com/hyperledger/indy-sdk-go/common/logging"
	"github.com/hyperledger/indy-sdk-go/indy"
)

var logger = logging.MustGetLogger("indy-sdk")

// BuildNYMRequest builds a NYM request. Request to create a new NYM record for a specific user.
//
// submitterDid DID of the submitter stored in secured Wallet.
// targetDid    Target DID as base58-encoded string for 16 or 32 bit DID value.
// verkey       Target identity verification key as base58-encoded string.
// alias        NYM's alias.
// role         Role of a user NYM record: nil (common USER), Trustee, Steward, TrustAnchor, Reset (to reset the role)
func BuildNYMRequest(submitterDID, targetDID, verkey string, alias *types.Alias, role *role.Role) (nymReq string, err error) {
	nymReqChan, errChan := buildNYMRequest(submitterDID, targetDID, verkey, alias, role)
	select {
	case nymReq = <-nymReqChan:
	case err = <-errChan:
	}
	return
}

// SignAndSubmitRequest signs and submits request message to validator pool.
//
// Adds submitter information to passed request json, signs it with submitter
// sign key (see wallet_sign), and sends signed request message
// to validator pool (see write_request).
//
// pool         A Pool.
// wallet       A Wallet.
// submitterDid Id of Identity stored in secured Wallet.
// requestJson  Request data json.
func SignAndSubmitRequest(pool *pool.Pool, wallet *wallet.Wallet, submitterDID, requestJSON string) (responseJSON string, err error) {
	respChan, errChan := signAndSubmitRequest(pool, wallet, submitterDID, requestJSON)
	select {
	case responseJSON = <-respChan:
	case err = <-errChan:
	}
	return
}

// SubmitRequest publishes request message to validator pool (no signing, unlike sign_and_submit_request).
// The request is sent to the validator pool as is. It's assumed that it's already prepared.
//
// pool        The Pool to publish to.
// requestJson Request data json.
func SubmitRequest(pool *pool.Pool, requestJSON string) (responseJSON string, err error) {
	respChan, errChan := submitRequest(pool, requestJSON)
	select {
	case responseJSON = <-respChan:
	case err = <-errChan:
	}
	return
}

// BuildSchemaRequest builds a SCHEMA request. Request to add Credential's schema.
//
// submitterDid DID of the submitter stored in secured Wallet.
// Credential schema:
// {
// 	id: identifier of schema
// 	attrNames: array of attribute name strings
// 	name: Schema's name string
// 	version: Schema's version string,
// 	ver: Version of the Schema json
// }
func BuildSchemaRequest(submitterDID, data string) (request string, err error) {
	requestChan, errChan := buildSchemaRequest(submitterDID, data)
	select {
	case request = <-requestChan:
	case err = <-errChan:
	}
	return
}

// BuildGetSchemaRequest builds a GET_SCHEMA request. Request to get Credential's Schema.
//
// submitterDid DID of read request sender.
// id           Schema ID in ledger
func BuildGetSchemaRequest(submitterDID, id string) (request string, err error) {
	requestChan, errChan := buildGetSchemaRequest(submitterDID, id)
	select {
	case request = <-requestChan:
	case err = <-errChan:
	}
	return
}

// ParseGetSchemaResponse parses a GET_SCHEMA response to get Schema in the format compatible with Anoncreds API
//
// getSchemaResponse the response of GET_SCHEMA request.
// returns:
// {
//     id: identifier of schema
//     attrNames: array of attribute name strings
//     name: Schema's name string
//     version: Schema's version string
//     ver: Version of the Schema json
// }
func ParseGetSchemaResponse(response string) (id, json string, err error) {
	resultChan, errChan := parseGetSchemaResponse(response)
	select {
	case result := <-resultChan:
		id = result.id
		json = result.json
	case err = <-errChan:
	}
	return
}

// BuildCredDefRequest builds a CRED_DEF request. Request to add a credential definition (in particular, public key),
// that Issuer creates for a particular Credential Schema.
//
// @param submitterDid DID of the submitter stored in secured Wallet.
// @param data         Credential definition json
// {
//     id: string - identifier of credential definition
//     schemaId: string - identifier of stored in ledger schema
//     type: string - type of the credential definition. CL is the only supported type now.
//     tag: string - allows to distinct between credential definitions for the same issuer and schema
//     value: Dictionary with Credential Definition's data: {
//         primary: primary credential public key,
//         Optional<revocation>: revocation credential public key
//     },
//     ver: Version of the CredDef json
// }
func BuildCredDefRequest(submitterDID, data string) (request string, err error) {
	requestChan, errChan := buildCredDefRequest(submitterDID, data)
	select {
	case request = <-requestChan:
	case err = <-errChan:
	}
	return
}

// BuildGetCredDefRequest builds a GET_CRED_DEF request. Request to get a credential definition (in particular, public key),
// that Issuer creates for a particular Credential Schema.
//
// submitterDid DID of read request sender.
// id           Credential Definition ID in ledger.
func BuildGetCredDefRequest(submitterDID, id string) (request string, err error) {
	requestChan, errChan := buildGetCredDefRequest(submitterDID, id)
	select {
	case request = <-requestChan:
	case err = <-errChan:
	}
	return
}

// ParseGetCredDefResponse parses a GET_CRED_DEF response to get Credential Definition in the format compatible with Anoncreds API.
//
// getCredDefResponse response of GET_CRED_DEF request.
// return A Credential Definition ID and Credential Definition JSON.
// {
//     id: string - identifier of credential definition
//     schemaId: string - identifier of stored in ledger schema
//     type: string - type of the credential definition. CL is the only supported type now.
//     tag: string - allows to distinct between credential definitions for the same issuer and schema
//     value: Dictionary with Credential Definition's data: {
//         primary: primary credential public key,
//         Optional<revocation>: revocation credential public key
//     },
//     ver: Version of the Credential Definition json
// }
func ParseGetCredDefResponse(response string) (id, json string, err error) {
	resultChan, errChan := parseGetCredDefResponse(response)
	select {
	case result := <-resultChan:
		id = result.id
		json = result.json
	case err = <-errChan:
	}
	return
}

func buildNYMRequest(submitterDID, targetDID, verkey string, alias *types.Alias, role *role.Role) (chan string, chan error) {
	logger.Debugf("Building NYM request - SubmitterDID [%s], TargetDID [%s], VerKey [%s], Alias [%s], Role [%v]", submitterDID, targetDID, verkey, alias, role)

	reqChan := make(chan string)
	errChan := make(chan error, 1)

	if submitterDID == "" {
		errChan <- fmt.Errorf("Submitter DID must be specified")
		return reqChan, errChan
	}
	if targetDID == "" {
		errChan <- fmt.Errorf("Target DID must be specified")
		return reqChan, errChan
	}

	cb := func(err error, data callback.Data) {
		if err != nil {
			errChan <- err
		} else {
			reqChan <- data.(string)
		}
	}

	err := indy.BuildNYMRequest(submitterDID, targetDID, verkey, alias, role, cb)
	if err != nil {
		// Send the error immediately
		errChan <- err
	}

	return reqChan, errChan
}

func signAndSubmitRequest(pool *pool.Pool, wallet *wallet.Wallet, submitterDID, requestJSON string) (chan string, chan error) {
	logger.Debugf("Signing and submitting request - Pool [%s], Wallet [%s], SubmitterDID [%s], JSON [%s]", pool.Name, wallet.Name, submitterDID, requestJSON)

	respChan := make(chan string)
	errChan := make(chan error, 1)

	if submitterDID == "" {
		errChan <- fmt.Errorf("Submitter DID must be specified")
		return respChan, errChan
	}
	if requestJSON == "" {
		errChan <- fmt.Errorf("Request JSON must be specified")
		return respChan, errChan
	}

	cb := func(err error, data callback.Data) {
		if err != nil {
			errChan <- err
		} else {
			respChan <- data.(string)
		}
	}

	err := indy.SignAndSubmitRequest(pool.Handle(), wallet.Handle(), submitterDID, requestJSON, cb)
	if err != nil {
		// Send the error immediately
		errChan <- err
	}

	return respChan, errChan
}

func submitRequest(pool *pool.Pool, requestJSON string) (chan string, chan error) {
	logger.Debugf("Submitting request - Pool [%s], JSON [%s]", pool.Name, requestJSON)

	respChan := make(chan string)
	errChan := make(chan error, 1)

	if requestJSON == "" {
		errChan <- fmt.Errorf("Request JSON must be specified")
		return respChan, errChan
	}

	cb := func(err error, data callback.Data) {
		if err != nil {
			errChan <- err
		} else {
			respChan <- data.(string)
		}
	}

	err := indy.SubmitRequest(pool.Handle(), requestJSON, cb)
	if err != nil {
		// Send the error immediately
		errChan <- err
	}

	return respChan, errChan
}

func buildSchemaRequest(submitterDID, data string) (chan string, chan error) {
	logger.Debugf("Building schema request - SubmitterDID [%s], Data [%s]", submitterDID, data)

	reqChan := make(chan string)
	errChan := make(chan error, 1)

	if submitterDID == "" {
		errChan <- fmt.Errorf("submitter DID must be specified")
		return reqChan, errChan
	}
	if data == "" {
		errChan <- fmt.Errorf("data must be specified")
		return reqChan, errChan
	}

	cb := func(err error, data callback.Data) {
		if err != nil {
			errChan <- err
		} else {
			reqChan <- data.(string)
		}
	}

	err := indy.BuildSchemaRequest(submitterDID, data, cb)
	if err != nil {
		// Send the error immediately
		errChan <- err
	}

	return reqChan, errChan
}

func buildGetSchemaRequest(submitterDID, id string) (chan string, chan error) {
	logger.Debugf("Building get-schema request - SubmitterDID [%s], ID [%s]", submitterDID, id)

	reqChan := make(chan string)
	errChan := make(chan error, 1)

	if submitterDID == "" {
		errChan <- fmt.Errorf("submitter DID must be specified")
		return reqChan, errChan
	}
	if id == "" {
		errChan <- fmt.Errorf("ID must be specified")
		return reqChan, errChan
	}

	cb := func(err error, data callback.Data) {
		if err != nil {
			errChan <- err
		} else {
			reqChan <- data.(string)
		}
	}

	err := indy.BuildGetSchemaRequest(submitterDID, id, cb)
	if err != nil {
		// Send the error immediately
		errChan <- err
	}

	return reqChan, errChan
}

type parseResponse struct {
	id   string
	json string
}

func parseGetSchemaResponse(response string) (chan *parseResponse, chan error) {
	logger.Debugf("Parsing get-schema response - Response [%s]", response)

	resultChan := make(chan *parseResponse)
	errChan := make(chan error, 1)

	if response == "" {
		errChan <- fmt.Errorf("response must be specified")
		return resultChan, errChan
	}

	cb := func(err error, data callback.Data) {
		if err != nil {
			errChan <- err
		} else {
			sd := data.([]string)
			resultChan <- &parseResponse{
				id:   sd[0],
				json: sd[1],
			}
		}
	}

	err := indy.ParseGetSchemaResponse(response, cb)
	if err != nil {
		// Send the error immediately
		errChan <- err
	}

	return resultChan, errChan
}

func buildCredDefRequest(submitterDID, data string) (chan string, chan error) {
	logger.Debugf("Building cred def request - SubmitterDID [%s], Data [%s]", submitterDID, data)

	reqChan := make(chan string)
	errChan := make(chan error, 1)

	if submitterDID == "" {
		errChan <- fmt.Errorf("submitter DID must be specified")
		return reqChan, errChan
	}
	if data == "" {
		errChan <- fmt.Errorf("data must be specified")
		return reqChan, errChan
	}

	cb := func(err error, data callback.Data) {
		if err != nil {
			errChan <- err
		} else {
			reqChan <- data.(string)
		}
	}

	err := indy.BuildCredDefRequest(submitterDID, data, cb)
	if err != nil {
		// Send the error immediately
		errChan <- err
	}

	return reqChan, errChan
}

func buildGetCredDefRequest(submitterDID, id string) (chan string, chan error) {
	logger.Debugf("Building get cred def request - SubmitterDID [%s], ID [%s]", submitterDID, id)

	reqChan := make(chan string)
	errChan := make(chan error, 1)

	if submitterDID == "" {
		errChan <- fmt.Errorf("submitter DID must be specified")
		return reqChan, errChan
	}
	if id == "" {
		errChan <- fmt.Errorf("ID must be specified")
		return reqChan, errChan
	}

	cb := func(err error, data callback.Data) {
		if err != nil {
			errChan <- err
		} else {
			reqChan <- data.(string)
		}
	}

	err := indy.BuildGetCredDefRequest(submitterDID, id, cb)
	if err != nil {
		// Send the error immediately
		errChan <- err
	}

	return reqChan, errChan
}

func parseGetCredDefResponse(response string) (chan *parseResponse, chan error) {
	logger.Debugf("Parsing get-cred-def response - Response [%s]", response)

	resultChan := make(chan *parseResponse)
	errChan := make(chan error, 1)

	if response == "" {
		errChan <- fmt.Errorf("response must be specified")
		return resultChan, errChan
	}

	cb := func(err error, data callback.Data) {
		if err != nil {
			errChan <- err
		} else {
			sd := data.([]string)
			resultChan <- &parseResponse{
				id:   sd[0],
				json: sd[1],
			}
		}
	}

	err := indy.ParseGetCredDefResponse(response, cb)
	if err != nil {
		// Send the error immediately
		errChan <- err
	}

	return resultChan, errChan
}
