/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ledger

import (
	"testing"

	"github.com/hyperledger/indy-sdk-go/common/indyerror"
	"github.com/hyperledger/indy-sdk-go/common/role"
	"github.com/hyperledger/indy-sdk-go/pool"
	"github.com/hyperledger/indy-sdk-go/wallet"
)

const (
	did1       = "XsKV7grR1BUE29mG2Fm2k1"
	did2       = "XsKV7grR1BUE29mG2Fm2k2"
	verKey1    = "CnEDk9HrMnmiHXEV1WFgbVCRteYnPqsJwrTdcZaNhFVW"
	poolName   = "pool1"
	poolConfig = `{"genesis_txn": "../test/testdata/docker_pool_transactions_genesis"}`
	walletName = "wallet1"
)

func TestLedger(t *testing.T) {
	pool, err := getPool(t, poolName, poolConfig)
	if err != nil {
		t.Fatalf("Error received from getPool: %s", err)
	}
	defer pool.Close()

	wallet, err := getWallet(walletName, poolName)
	if err != nil {
		t.Fatalf("Error received from getWallet: %s", err)
	}
	defer wallet.Close()

	nymReq, err := BuildNYMRequest(did1, did2, verKey1, nil, role.Steward)
	if err != nil {
		t.Fatalf("Error received from BuildNYMRequest: %s", err)
	}
	t.Logf("Nym request created: %s", nymReq)

	resp, err := SignAndSubmitRequest(pool, wallet, did1, nymReq)
	if err != nil {
		t.Fatalf("Error received from SignAndSubmitRequest: %s", err)
	}
	t.Logf("Response received: %s", resp)
}

func getPool(t *testing.T, poolName, configPath string) (*pool.Pool, error) {
	err := pool.Create(poolName, configPath)
	if err != nil && indyerror.Code(err) != indyerror.PoolLedgerConfigAlreadyExistsError {
		t.Fatalf("Error received from Create: %s", err)
	}
	return pool.Open(poolName, "")
}

func getWallet(walletName, poolName string) (*wallet.Wallet, error) {
	err := wallet.Create(poolName, walletName, "", "", "")
	if err != nil && indyerror.Code(err) != indyerror.WalletAlreadyExistsError {
		return nil, err
	}
	return wallet.Open(walletName, "", "")
}
