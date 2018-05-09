/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package crypto

import (
	"testing"
)

const (
	verKey1 = "CnEDk9HrMnmiHXEV1WFgbVCRteYnPqsJwrTdcZaNhFVW"
)

func TestAnonCrypt(t *testing.T) {
	message := []byte("some message")

	encryptedMsg, err := AnonCrypt(verKey1, message)
	if err != nil {
		t.Fatalf("Error received from AnonCrypt: %s", err)
	}
	t.Logf("Got encrypted message [%#x]", encryptedMsg)
}
