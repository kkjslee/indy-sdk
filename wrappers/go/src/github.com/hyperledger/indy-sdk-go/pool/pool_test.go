/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package pool

import (
	"testing"
)

const (
	config = `{"genesis_txn": "../test/testdata/docker_pool_transactions_genesis"}`
)

func TestPool(t *testing.T) {
	poolName := "pool3"

	// Delete any pool that was created from a previous test
	Delete(poolName)

	err := Create(poolName, config)
	if err != nil {
		t.Fatalf("Error received from Create: %s", err)
	}

	pool, err := Open(poolName, "")
	if err != nil {
		t.Fatalf("Error received from Open: %s", err)
	}
	t.Logf("Pool opened: %s", pool.Name)

	err = pool.Refresh()
	if err != nil {
		t.Fatalf("Error received from Refresh: %s", err)
	}

	err = pool.Close()
	if err != nil {
		t.Fatalf("Error received from Close: %s", err)
	}

	pools, err := List()
	if err != nil {
		t.Fatalf("Error received from ListPools: %s", err)
	}
	t.Logf("Pools: %v", pools)

	err = Delete(poolName)
	if err != nil {
		t.Fatalf("Error received from Delete: %s", err)
	} else {
		t.Logf("Success received from Delete")
	}
}
