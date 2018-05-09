/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package wallet

import (
	"testing"
)

func TestWallet(t *testing.T) {
	walletName := "wallet1"

	// Delete any wallet that was created from a previous test
	Delete(walletName, "")

	err := Create("poolx", walletName, "", "", "")
	if err != nil {
		t.Fatalf("Error received from Create: %s", err)
	} else {
		t.Log("Success received from Create")
	}

	w, err := Open(walletName, "", "")
	if err != nil {
		t.Fatalf("Error received from Open: %s", err)
	}

	t.Logf("Wallet opened: %s", w.Name)

	err = w.Close()
	if err != nil {
		t.Fatalf("Error received from Close: %s", err)
	} else {
		t.Log("Success received from Close")
	}

	err = Delete(walletName, "")
	if err != nil {
		t.Fatalf("Error received from Delete: %s", err)
	} else {
		t.Log("Success received from Delete")
	}
}

type mockWalletType struct {
}

func TestRegisterWalletType(t *testing.T) {
	t.SkipNow()
	walletTypeName := "walletType1"

	walletType := &mockWalletType{}
	err := RegisterType(walletTypeName, walletType)
	if err != nil {
		t.Fatalf("Error received from RegisterType: %s", err)
	}
	t.Log("Success received from RegisterType")
}
