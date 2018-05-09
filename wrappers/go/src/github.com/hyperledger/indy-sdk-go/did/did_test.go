/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package did

import (
	"fmt"
	"testing"

	"github.com/hyperledger/indy-sdk-go/common/indyerror"
	"github.com/hyperledger/indy-sdk-go/wallet"
)

const (
	did1  = "VsKV7grR1BUE29mG2Fm2kZ"
	seed1 = "00000000000000000000000000000My1"
	seed2 = "00000000000000000000000000000My2"
)

func TestDID(t *testing.T) {
	// didJSON := getJSON("", seed1, "", nil)
	didJSON := getJSON(did1, "", "", nil)

	w, err := getWallet("wallet1", "pool1")
	if err != nil {
		t.Fatalf("error getting wallet: %s", err)
	}
	defer w.Close()

	didInfo, err := CreateAndStoreMyDID(w, didJSON)
	if err != nil {
		t.Fatalf("Error received from CreateAndStoreMyDID: %s", err)
	}
	t.Logf("DID created and stored successfully - DID [%s], VerKey [%s]", didInfo.DID, didInfo.VerKey)

	if didInfo.DID != did1 {
		t.Fatalf("Expecting DID to be [%s] but got [%s]", did1, didInfo.DID)
	}

	// Try again with a duplicate
	didInfo, err = CreateAndStoreMyDID(w, didJSON)
	if err == nil {
		t.Fatalf("Expecting error for duplicate DID but got success")
	}
	errCode := indyerror.Code(err)
	if errCode != indyerror.DidAlreadyExistsError {
		t.Fatalf("Expecting error [%s] but got [%s]", indyerror.New(indyerror.DidAlreadyExistsError), err)
	}

	// keyChan, errChan := KeyForDID(p)
}

func getWallet(walletName, poolName string) (*wallet.Wallet, error) {
	err := wallet.Create(poolName, walletName, "", "", "")
	if err != nil && indyerror.Code(err) != indyerror.WalletAlreadyExistsError {
		return nil, err
	}
	return wallet.Open(walletName, "", "")
}

func getJSON(did, seed, cryptoType string, cid *bool) string {
	json := "{"
	if did != "" {
		json += fmt.Sprintf(`"did":"%s"`, did)
	}
	if seed != "" {
		json += fmt.Sprintf(`"seed":"%s"`, seed)
	}
	if cryptoType != "" {
		json += fmt.Sprintf(`"crypto_type":"%s"`, cryptoType)
	}
	if cid != nil {
		json += fmt.Sprintf(`"cid":"%t"`, *cid)
	}
	json += "}"
	return json
}
