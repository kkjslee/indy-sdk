/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package pool

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/indy-sdk-go/common/callback"
	"github.com/hyperledger/indy-sdk-go/common/logging"
	"github.com/hyperledger/indy-sdk-go/common/types"
	"github.com/hyperledger/indy-sdk-go/indy"
)

var logger = logging.MustGetLogger("indy-sdk")

// Pool is the Pool Ledger
type Pool struct {
	Name   string `json:"pool"`
	handle types.Handle
}

// Handle returns the Indy handle of the pool
func (p *Pool) Handle() types.Handle {
	return p.handle
}

// Create creates a new local pool ledger configuration that can be used later to connect pool nodes.
//
// configName Name of the pool ledger configuration.
// config Pool configuration json. if NULL, then default config will be used.
func Create(name string, configPath string) error {
	return <-create(name, configPath)
}

// Delete deletes created pool ledger configuration.
//
// configName Name of the pool ledger configuration to delete.
func Delete(name string) error {
	return <-delete(name)
}

// Open opens pool ledger and performs connecting to pool nodes.
//
// configName Name of the pool ledger configuration.
// config Runtime pool configuration json. If empty, then default config will be used.
func Open(name string, config string) (*Pool, error) {
	poolChan, errChan := open(name, config)
	select {
	case pool := <-poolChan:
		return pool, nil
	case err := <-errChan:
		return nil, err
	}
}

// List returns a list of names of local pools.
func List() ([]string, error) {
	poolNameChan, errChan := list()
	select {
	case poolName := <-poolNameChan:
		return poolName, nil
	case err := <-errChan:
		return nil, err
	}
}

// Refresh refreshes a local copy of a pool ledger and updates pool nodes connections.
func (p *Pool) Refresh() error {
	return <-p.refresh()
}

// Close closes opened pool ledger, opened nodes connections and frees allocated resources.
func (p *Pool) Close() error {
	return <-p.close()
}

func create(name string, configPath string) chan error {
	logger.Debugf("Creating pool ledger: %s - Config Path: %s", name, configPath)

	errChan := make(chan error, 1)

	if name == "" {
		errChan <- fmt.Errorf("pool name must be specified")
		return errChan
	}
	if configPath == "" {
		errChan <- fmt.Errorf("config path must be specified")
		return errChan
	}

	err := indy.CreatePoolLedgerConfig(name, configPath, callback.New(errChan))
	if err != nil {
		// Send the error immediately
		errChan <- err
	}

	return errChan
}

func delete(name string) chan error {
	logger.Debugf("Deleting pool ledger: %s", name)

	errChan := make(chan error, 1)

	if name == "" {
		errChan <- fmt.Errorf("pool name must be specified")
		return errChan
	}

	err := indy.DeletePoolLedgerConfig(name, callback.New(errChan))
	if err != nil {
		// Send the error immediately
		errChan <- err
	}

	return errChan
}

func open(name string, config string) (chan *Pool, chan error) {
	logger.Debugf("Opening pool ledger [%s]", name)

	poolChan := make(chan *Pool)
	errChan := make(chan error, 1)

	if name == "" {
		errChan <- fmt.Errorf("pool name must be specified")
		return poolChan, errChan
	}

	cb := func(err error, data callback.Data) {
		if err != nil {
			errChan <- err
		} else {
			handle := data.(types.Handle)
			poolChan <- &Pool{
				Name:   name,
				handle: handle,
			}
		}
	}

	err := indy.OpenPoolLedger(name, config, cb)
	if err != nil {
		// Send the error immediately
		errChan <- err
	}

	return poolChan, errChan
}

func list() (chan []string, chan error) {
	logger.Debugf("Listing pools...")

	poolsChan := make(chan []string)
	errChan := make(chan error, 1)

	cb := func(err error, data callback.Data) {
		if err != nil {
			errChan <- err
		} else {
			json := data.(string)
			logger.Debugf("Pool ledger list: %s", json)
			pools, err := asPools(json)
			if err != nil {
				// FIXME: Compose better error
				errChan <- err
			} else {
				poolsChan <- pools
			}
		}
	}

	err := indy.ListPools(cb)
	if err != nil {
		// Send the error immediately
		errChan <- err
	}

	return poolsChan, errChan
}

func (p *Pool) refresh() chan error {
	logger.Debugf("Refreshing pool [%s]...", p.Name)

	errChan := make(chan error, 1)
	err := indy.RefreshPoolLedger(p.handle, callback.New(errChan))
	if err != nil {
		// Send the error immediately
		errChan <- err
	}

	return errChan
}

func (p *Pool) close() chan error {
	logger.Debugf("Closing pool [%s]...", p.Name)

	errChan := make(chan error, 1)
	err := indy.ClosePoolLedger(p.handle, callback.New(errChan))
	if err != nil {
		// Send the error immediately
		errChan <- err
	}

	return errChan
}

func asPools(poolsJSON string) ([]string, error) {
	var pools []Pool
	if err := json.Unmarshal([]byte(poolsJSON), &pools); err != nil {
		return nil, err
	}
	var poolNames []string
	for _, pool := range pools {
		poolNames = append(poolNames, pool.Name)
	}
	return poolNames, nil
}
