/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package wallet

import (
	"fmt"

	"github.com/hyperledger/indy-sdk-go/common/callback"
	"github.com/hyperledger/indy-sdk-go/common/logging"
	"github.com/hyperledger/indy-sdk-go/common/types"
	"github.com/hyperledger/indy-sdk-go/indy"
)

var logger = logging.MustGetLogger("indy-sdk")

// Wallet is the wallet
type Wallet struct {
	Name   string `json:"pool"`
	handle types.Handle
}

// Type is a custom wallet type
type Type interface {
}

// Create creates a new secure wallet with the given unique name.
//
// poolName Name of the pool that corresponds to this wallet.
// name Name of the wallet.
// walletType Type of the wallet. Defaults to 'default'.
// config Wallet configuration json. List of supported keys are defined by wallet type.
// credentials Wallet credentials json. List of supported keys are defined by wallet type.
func Create(poolName, name, walletType, config, credentials string) error {
	return <-create(poolName, name, walletType, config, credentials)
}

// Delete deletes the given wallet
// name Name of the wallet to delete.
// credentials Wallet credentials json. List of supported keys are defined by wallet type.
func Delete(name, credentials string) error {
	return <-delete(name, credentials)
}

// Open opens the wallet with specific name.
//
// name Name of the wallet.
// runtimeConfig Runtime wallet configuration json. if NULL, then default runtime_config will be used.
// credentials Wallet credentials json. List of supported keys are defined by wallet type.
// return A future that resolves no value.
func Open(name, config, credentials string) (wallet *Wallet, err error) {
	wChan, errChan := open(name, config, credentials)
	select {
	case wallet = <-wChan:
	case err = <-errChan:
	}
	return
}

// RegisterType registers a custome wallet implementation.
func RegisterType(typeName string, walletType Type) error {
	return <-registerType(typeName, walletType)
}

// Handle returns the Indy handle to the wallet
func (w *Wallet) Handle() types.Handle {
	return w.handle
}

// Close closes the wallet
func (w *Wallet) Close() error {
	return <-w.close()
}

func (w *Wallet) close() chan error {
	logger.Debugf("Closing wallet [%s]...", w.Name)

	errChan := make(chan error, 1)
	err := indy.CloseWallet(w.handle, callback.New(errChan))
	if err != nil {
		// Send the error immediately
		errChan <- err
	}

	return errChan
}

func create(poolName, name, walletType, config, credentials string) chan error {
	logger.Debugf("Creating wallet: %s, Pool [%s], Type [%s], Config [%s], Credentials [%s]", name, poolName, walletType, config, credentials)

	errChan := make(chan error, 1)

	if poolName == "" {
		errChan <- fmt.Errorf("pool name must be specified")
		return errChan
	}
	if name == "" {
		errChan <- fmt.Errorf("wallet name must be specified")
		return errChan
	}

	err := indy.CreateWallet(poolName, name, walletType, config, credentials, callback.New(errChan))
	if err != nil {
		// Send the error immediately
		errChan <- err
	}

	return errChan
}

func delete(name, credentials string) chan error {
	logger.Debugf("Deleting wallet [%s] - Credentials [%s]", name, credentials)

	errChan := make(chan error, 1)

	if name == "" {
		errChan <- fmt.Errorf("wallet must be specified")
		return errChan
	}

	err := indy.DeleteWallet(name, credentials, callback.New(errChan))
	if err != nil {
		// Send the error immediately
		errChan <- err
	}

	return errChan
}

func open(name, config, credentials string) (chan *Wallet, chan error) {
	logger.Debugf("Opening wallet: %s, Config [%s], Credentials [%s]", name, config, credentials)

	errChan := make(chan error, 1)
	walletChan := make(chan *Wallet)

	if name == "" {
		errChan <- fmt.Errorf("wallet name must be specified")
		return walletChan, errChan
	}

	cb := func(err error, data callback.Data) {
		if err != nil {
			logger.Debugf("Error opening wallet [%s]: %s", name, err)
			errChan <- err
		} else {
			logger.Debugf("Successfully opened wallet ledger [%s]", name)
			handle := data.(types.Handle)
			walletChan <- &Wallet{
				Name:   name,
				handle: handle,
			}
		}
	}

	err := indy.OpenWallet(name, config, credentials, cb)
	if err != nil {
		// Send the error immediately
		errChan <- err
	}

	return walletChan, errChan
}

func registerType(typeName string, walletType Type) chan error {
	logger.Debugf("Registering wallet type [%s]", typeName)

	errChan := make(chan error, 1)

	if typeName == "" {
		errChan <- fmt.Errorf("wallet type name must be specified")
		return errChan
	}

	// handle := callback.Register(func(err error, _ callback.Data) {
	// 	if err != nil {
	// 		logger.Debugf("Error registering wallet type [%s]: %s", typeName, err)
	// 	} else {
	// 		logger.Debugf("Successfully registered wallet type [%s]", typeName)
	// 	}
	// 	errChan <- err
	// })

	// if errCode := C.indy_register_wallet_type(
	// 	(C.indy_handle_t)(handle), csPoolName, csName, csType, csConfig, csCredentials,
	// 	(C.callback_fcn)(unsafe.Pointer(C.def_callback))); errCode != indyerror.Success {
	// 	errChan <- indyerror.New(int32(errCode))
	// }

	// return errChan
	panic("not implemented")
}
