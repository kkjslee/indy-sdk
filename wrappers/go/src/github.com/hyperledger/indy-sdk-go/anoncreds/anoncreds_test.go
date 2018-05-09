/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package anoncreds

import (
	"testing"
)

const (
	did1             = "VsKV7grR1BUE29mG2Fm2kZ"
	verKey1          = "CnEDk9HrMnmiHXEV1WFgbVCRteYnPqsJwrTdcZaNhFVW"
	schemaAttributes = `["name", "age", "sex", "height"]`
)

func TestIssuerCreateSchema(t *testing.T) {
	schemaID, schemaJSON, err := IssuerCreateSchema(did1, "schema1", "1.0", schemaAttributes)
	if err != nil {
		t.Fatalf("Error received from IssuerCreateSchema: %s", err)
	}
	t.Logf("Created schema - ID [%s], JSON [%s]", schemaID, schemaJSON)
}
