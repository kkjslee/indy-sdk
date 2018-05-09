/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package did

import (
	"fmt"

	"github.com/hyperledger/indy-sdk-go/pool"

	"github.com/hyperledger/indy-sdk-go/common/callback"
	"github.com/hyperledger/indy-sdk-go/common/logging"
	"github.com/hyperledger/indy-sdk-go/indy"
	"github.com/hyperledger/indy-sdk-go/wallet"
)

var logger = logging.MustGetLogger("indy-sdk")

// Info is the Decentralized ID
type Info struct {
	DID    string `json:"did"`
	VerKey string `json:"verkey"`
}

// CreateAndStoreMyDID creates keys (signing and encryption keys) for a new
// DID (owned by the caller of the library).
// Identity's DID must be either explicitly provided, or taken as the first 16 bit of verkey.
// Saves the Identity DID with keys in a secured Wallet, so that it can be used to sign
// and encrypt transactions.
//
// wallet  The wallet.
// didJson Identity information as json.
func CreateAndStoreMyDID(wallet *wallet.Wallet, didJSON string) (didInfo *Info, err error) {
	infoChan, errChan := createAndStoreMyDID(wallet, didJSON)
	select {
	case didInfo = <-infoChan:
	case err = <-errChan:
	}
	return
}

// KeyForDID returns ver key (key id) for the given DID.
//
// "keyForDid" call follow the idea that we resolve information about their DID from
// the ledger with cache in the local wallet. The "openWallet" call has freshness parameter
// that is used for checking the freshness of cached pool value.
//
// Note if you don't want to resolve their DID info from the ledger you can use
// "keyForLocalDid" call instead that will look only to local wallet and skip
// freshness checking.
//
// Note that "createAndStoreMyDid" makes similar wallet record as "createKey".
// As result we can use returned ver key in all generic crypto and messaging functions.
//
// pool   The pool.
// wallet The wallet.
// did    The DID to resolve key.
func KeyForDID(pool *pool.Pool, wallet *wallet.Wallet, did string) (key string, err error) {
	keyChan, errChan := keyForDID(pool, wallet, did)
	select {
	case key = <-keyChan:
	case err = <-errChan:
	}
	return
}

func createAndStoreMyDID(wallet *wallet.Wallet, didJSON string) (chan *Info, chan error) {
	logger.Debugf("Creating and storing DID - Wallet [%s] - Data: %s", wallet.Name, didJSON)

	didChan := make(chan *Info)
	errChan := make(chan error, 1)

	if didJSON == "" {
		errChan <- fmt.Errorf("JSON for DID must be specified")
		return didChan, errChan
	}

	cb := func(err error, data callback.Data) {
		if err != nil {
			errChan <- err
		} else {
			sd := data.([]string)
			didChan <- &Info{
				DID:    sd[0],
				VerKey: sd[1],
			}
		}
	}

	err := indy.CreateAndStoreMyDID(wallet.Handle(), didJSON, cb)
	if err != nil {
		// Send the error immediately
		errChan <- err
	}

	return didChan, errChan
}

func keyForDID(pool *pool.Pool, wallet *wallet.Wallet, did string) (chan string, chan error) {
	logger.Debugf("Getting key for DID [%s] - Pool [%s], Wallet [%s]", did, pool.Name, wallet.Name)

	keyChan := make(chan string)
	errChan := make(chan error, 1)

	if did == "" {
		errChan <- fmt.Errorf("DID must be specified")
		return keyChan, errChan
	}

	cb := func(err error, data callback.Data) {
		if err != nil {
			errChan <- err
		} else {
			key := data.(string)
			keyChan <- key
		}
	}

	err := indy.KeyForDID(pool.Handle(), wallet.Handle(), did, cb)
	if err != nil {
		// Send the error immediately
		errChan <- err
	}

	return keyChan, errChan
}
