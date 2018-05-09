/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package callback

import (
	"math/rand"
	"sync"

	"github.com/hyperledger/indy-sdk-go/common/types"
)

// Data is the data associated with the callback
type Data interface {
}

// Callback is the function invoked by the Indy SDK
type Callback func(error, Data)

// New creates a new callback
func New(errChan chan error) Callback {
	return func(err error, _ Data) {
		errChan <- err
	}
}

// Register registers the callback and returns a Handle
// which may be used to retrieve the callback
func Register(cb Callback) types.Handle {
	return getRegistry().register(cb)
}

// Remove removes and returns the callback associated with the handle
func Remove(handle types.Handle) (Callback, bool) {
	return getRegistry().remove(handle)
}

type registry struct {
	registry map[types.Handle]Callback
	mutex    sync.Mutex
}

func (r *registry) register(cb Callback) types.Handle {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for {
		handle := types.Handle(rand.Int31())
		if _, ok := r.registry[handle]; !ok {
			r.registry[handle] = cb
			return handle
		}
	}
}

func (r *registry) remove(handle types.Handle) (Callback, bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	cb, ok := r.registry[handle]
	if ok {
		delete(r.registry, handle)
	}
	return cb, ok
}

var initCBRegistry sync.Once
var cbRegistry *registry

func getRegistry() *registry {
	initCBRegistry.Do(func() {
		cbRegistry = &registry{
			registry: make(map[types.Handle]Callback),
		}
	})
	return cbRegistry
}
