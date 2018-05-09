/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package role

var (
	// Trustee is the Trustee role
	Trustee = NewRole("TRUSTEE")

	// Steward is the Steward role
	Steward = NewRole("STEWARD")

	// TrustAnchor is the Trust Anchor role
	TrustAnchor = NewRole("TRUST_ANCHOR")

	// Reset resets the role
	Reset = NewRole("")
)

// Role contains a role
type Role struct {
	value string
}

// NewRole creates a new Role
func NewRole(value string) *Role {
	return &Role{value: value}
}

// String returns the string representation of the role
func (r *Role) String() string {
	return r.value
}
