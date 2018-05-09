/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package crypto

import (
	"fmt"

	"github.com/hyperledger/indy-sdk-go/common/callback"
	"github.com/hyperledger/indy-sdk-go/common/logging"
	"github.com/hyperledger/indy-sdk-go/indy"
	"github.com/hyperledger/indy-sdk-go/wallet"
)

var logger = logging.MustGetLogger("indy-sdk")

// AnonCrypt encrypts a message by anonymous-encryption scheme.
//
// Sealed boxes are designed to anonymously send messages to a Recipient given its public key.
// Only the Recipient can decrypt these messages, using its private key.
// While the Recipient can verify the integrity of the message, it cannot verify the identity of the Sender.
// Note to use DID keys with this function you can call keyForDid to get key id (verkey)
// for specific DID.
//
// recipientVK verkey of message recipient
// message a message to be signed
func AnonCrypt(recipientVK string, message []byte) (encryptedMsg []byte, err error) {
	respChan, errChan := anonCrypt(recipientVK, message)
	select {
	case encryptedMsg = <-respChan:
	case err = <-errChan:
	}
	return
}

// AnonDecrypt decrypts a message by anonymous-encryption scheme.
//
// Sealed boxes are designed to anonymously send messages to a Recipient given its public key.
// Only the Recipient can decrypt these messages, using its private key.
// While the Recipient can verify the integrity of the message, it cannot verify the identity of the Sender.
// Note to use DID keys with this function you can call indy_key_for_did to get key id (verkey)
// for specific DID.
//
// wallet       The wallet.
// recipientVk  Id (verkey) of my key. The key must be created by calling createKey or createAndStoreMyDid
// encryptedMsg encrypted message
func AnonDecrypt(wallet *wallet.Wallet, recipientVK string, encryptedMsg []byte) (decryptedMsg []byte, err error) {
	respChan, errChan := anonDecrypt(wallet, recipientVK, encryptedMsg)
	select {
	case decryptedMsg = <-respChan:
	case err = <-errChan:
	}
	return
}

// AuthCrypt encrypts a message by authenticated-encryption scheme.
// Sender can encrypt a confidential message specifically for Recipient, using Sender's public key.
// Using Recipient's public key, Sender can compute a shared secret key.
// Using Sender's public key and his secret key, Recipient can compute the exact same shared secret key.
// That shared secret key can be used to verify that the encrypted message was not tampered with,
// before eventually decrypting it.
// Recipient only needs Sender's public key, the nonce and the ciphertext to peform decryption.
// The nonce doesn't have to be confidential.
// Note to use DID keys with this function you can call indy_key_for_did to get key id (verkey)
// for specific DID.
//
// wallet  The wallet.
// senderVK    id (verkey) of my key. The key must be created by calling indy_create_key or indy_create_and_store_my_did
// recipientVK id (verkey) of their key
// message a message to be signed
func AuthCrypt(wallet *wallet.Wallet, senderVK, recipientVK string, message []byte) (encryptedMsg []byte, err error) {
	respChan, errChan := authCrypt(wallet, senderVK, recipientVK, message)
	select {
	case encryptedMsg = <-respChan:
	case err = <-errChan:
	}
	return
}

// AuthDecrypt decrypts a message by authenticated-encryption scheme.
// Sender can encrypt a confidential message specifically for Recipient, using Sender's public key.
// Using Recipient's public key, Sender can compute a shared secret key.
// Using Sender's public key and his secret key, Recipient can compute the exact same shared secret key.
// That shared secret key can be used to verify that the encrypted message was not tampered with,
// before eventually decrypting it.
// Recipient only needs Sender's public key, the nonce and the ciphertext to peform decryption.
// The nonce doesn't have to be confidential.
// Note to use DID keys with this function you can call indy_key_for_did to get key id (verkey)
// for specific DID.
//
// wallet       The wallet.
// recipientVk  Id (verkey) of my key. The key must be created by calling createKey or createAndStoreMyDid
// encryptedMsg Encrypted message
func AuthDecrypt(wallet *wallet.Wallet, recipientVK string, message []byte) (sender string, decryptedMsg []byte, err error) {
	respChan, errChan := authDecrypt(wallet, recipientVK, message)
	select {
	case resp := <-respChan:
		sender = resp.sender
		decryptedMsg = resp.message
	case err = <-errChan:
	}
	return
}

func anonCrypt(recipientVK string, message []byte) (chan []byte, chan error) {
	logger.Debugf("Anonymously encrypting message - RecipientVK [%s] - Message: [%s]", recipientVK, message)

	respChan := make(chan []byte)
	errChan := make(chan error, 1)

	if recipientVK == "" {
		errChan <- fmt.Errorf("recipient verification key must be specified")
		return respChan, errChan
	}

	cb := func(err error, data callback.Data) {
		if err != nil {
			errChan <- err
		} else {
			respChan <- data.([]byte)
		}
	}

	err := indy.AnonCrypt(recipientVK, message, cb)
	if err != nil {
		// Send the error immediately
		errChan <- err
	}

	return respChan, errChan
}

func anonDecrypt(wallet *wallet.Wallet, recipientVK string, encryptedMsg []byte) (chan []byte, chan error) {
	logger.Debugf("Anonymously decrypting message - RecipientVK [%s] - Message: [%#x]", recipientVK, encryptedMsg)

	respChan := make(chan []byte)
	errChan := make(chan error, 1)

	if recipientVK == "" {
		errChan <- fmt.Errorf("recipient verification key must be specified")
		return respChan, errChan
	}
	if len(encryptedMsg) == 0 {
		errChan <- fmt.Errorf("encrypted message must be specified")
		return respChan, errChan
	}

	cb := func(err error, data callback.Data) {
		if err != nil {
			errChan <- err
		} else {
			respChan <- data.([]byte)
		}
	}

	err := indy.AnonDecrypt(wallet.Handle(), recipientVK, encryptedMsg, cb)
	if err != nil {
		// Send the error immediately
		errChan <- err
	}

	return respChan, errChan
}

func authCrypt(wallet *wallet.Wallet, senderVK, recipientVK string, message []byte) (chan []byte, chan error) {
	logger.Debugf("Auth encrypting message - Wallet [%s], SenderVK [%s], RecipientVK [%s] - Message: [%s]", wallet.Name, senderVK, recipientVK, message)

	respChan := make(chan []byte)
	errChan := make(chan error, 1)

	if senderVK == "" {
		errChan <- fmt.Errorf("sender verification key must be specified")
		return respChan, errChan
	}
	if recipientVK == "" {
		errChan <- fmt.Errorf("recipient verification key must be specified")
		return respChan, errChan
	}

	cb := func(err error, data callback.Data) {
		if err != nil {
			errChan <- err
		} else {
			respChan <- data.([]byte)
		}
	}

	err := indy.AuthCrypt(wallet.Handle(), senderVK, recipientVK, message, cb)
	if err != nil {
		// Send the error immediately
		errChan <- err
	}

	return respChan, errChan
}

type authDecryptResponse struct {
	sender  string
	message []byte
}

func authDecrypt(wallet *wallet.Wallet, recipientVK string, message []byte) (chan *authDecryptResponse, chan error) {
	logger.Debugf("Auth decrypting message - Wallet [%s], RecipientVK [%s] - Message: [%s]", wallet.Name, recipientVK, message)

	respChan := make(chan *authDecryptResponse)
	errChan := make(chan error, 1)

	if recipientVK == "" {
		errChan <- fmt.Errorf("recipient verification key must be specified")
		return respChan, errChan
	}

	cb := func(err error, data callback.Data) {
		if err != nil {
			errChan <- err
		} else {
			senderAndMessage := data.([]interface{})
			respChan <- &authDecryptResponse{
				sender:  senderAndMessage[0].(string),
				message: senderAndMessage[1].([]byte),
			}
		}
	}

	err := indy.AuthDecrypt(wallet.Handle(), recipientVK, message, cb)
	if err != nil {
		// Send the error immediately
		errChan <- err
	}

	return respChan, errChan
}
