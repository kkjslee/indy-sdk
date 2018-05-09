/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package types

// Handle is used for interchanging entities with the Indy "C" interface
type Handle int

// Alias contains an alias
type Alias struct {
	value string
}

// NewAlias creates a new Alias
func NewAlias(value string) *Alias {
	return &Alias{value: value}
}

// String returns the string representation of the alias
func (a *Alias) String() string {
	return a.value
}
