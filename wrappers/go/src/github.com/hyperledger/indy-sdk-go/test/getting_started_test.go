/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"fmt"
	"testing"

	"github.com/hyperledger/indy-sdk-go/anoncreds"
	"github.com/hyperledger/indy-sdk-go/common/indyerror"
	"github.com/hyperledger/indy-sdk-go/common/role"
	"github.com/hyperledger/indy-sdk-go/crypto"
	"github.com/hyperledger/indy-sdk-go/did"
	"github.com/hyperledger/indy-sdk-go/ledger"
	"github.com/hyperledger/indy-sdk-go/pool"
	"github.com/hyperledger/indy-sdk-go/test/assert"
	"github.com/hyperledger/indy-sdk-go/test/json"
	"github.com/hyperledger/indy-sdk-go/wallet"
)

var poolConfig = `{"genesis_txn": "./testdata/docker_pool_transactions_genesis"}`

func TestGettingStarted(t *testing.T) {
	fmt.Println("Getting started -> started")
	poolName := "pool1"

	// Delete any pool that was created by previous test
	pool.Delete(poolName)

	fmt.Println("Open Pool Ledger")
	err := pool.Create(poolName, poolConfig)
	assert.NoError(t, err)

	p, err := pool.Open(poolName, "")
	assert.NoError(t, err)
	defer p.Close()

	fmt.Println("==============================")
	fmt.Println("=== Getting Trust Anchor credentials for Faber, Acme, Thrift and Government  ==")
	fmt.Println("------------------------------")
	fmt.Println(`"Sovrin Steward" -> Create wallet`)

	stewardWalletName := "sovrin_steward_wallet"
	stewardWallet, err := getWallet(poolName, stewardWalletName)
	assert.NoError(t, err)
	defer stewardWallet.Close()

	fmt.Println(`"Sovrin Steward" -> Create and store in Wallet DID from seed`)

	stewardDIDInfo, err := did.CreateAndStoreMyDID(stewardWallet, `{"seed": "000000000000000000000000Steward1"}`)
	assert.NoError(t, err)

	fmt.Println("==============================")
	fmt.Println("== Getting Trust Anchor credentials - Government Onboarding  ==")
	fmt.Println("------------------------------")

	governmentWallet, stewardGovernmentDIDInfo, governmentStewardDIDInfo, _, err := onboarding(p, "Sovrin Steward", stewardWallet, stewardDIDInfo.DID, "Government", nil, "government_wallet")
	assert.NoError(t, err)

	fmt.Println("==============================")
	fmt.Println("== Getting Trust Anchor credentials - Government getting Verinym  ==")
	fmt.Println("------------------------------")

	governmentDIDInfo, err := getVerinym(t, p, "Sovrin Steward", stewardWallet, stewardDIDInfo.DID, stewardGovernmentDIDInfo.VerKey, "Government", governmentWallet, governmentStewardDIDInfo, role.TrustAnchor)
	assert.NoError(t, err)

	fmt.Println("==============================")
	fmt.Println("== Getting Trust Anchor credentials - Faber Onboarding  ==")
	fmt.Println("------------------------------")
	faberWallet, stewardFaberDIDInfo, faberStewardDIDInfo, _, err := onboarding(p, "Sovrin Steward", stewardWallet, stewardDIDInfo.DID, "Faber", nil, "faber_wallet")
	assert.NoError(t, err)

	fmt.Println("==============================")
	fmt.Println("== Getting Trust Anchor credentials - Faber getting Verinym  ==")
	fmt.Println("------------------------------")

	faberDIDInfo, err := getVerinym(t, p, "Sovrin Steward", stewardWallet, stewardDIDInfo.DID, stewardFaberDIDInfo.VerKey, "Faber", faberWallet, faberStewardDIDInfo, role.TrustAnchor)
	assert.NoError(t, err)

	fmt.Println("==============================")
	fmt.Println("== Getting Trust Anchor credentials - Acme Onboarding  ==")
	fmt.Println("------------------------------")

	acmeWallet, stewardAcmeDIDInfo, acmeStewardDIDInfo, _, err := onboarding(p, "Sovrin Steward", stewardWallet, stewardDIDInfo.DID, "Acme", nil, "acme_wallet")
	assert.NoError(t, err)

	fmt.Println("==============================")
	fmt.Println("== Getting Trust Anchor credentials - Acme getting Verinym  ==")
	fmt.Println("------------------------------")

	acmeDIDInfo, err := getVerinym(t, p, "Sovrin Steward", stewardWallet, stewardDIDInfo.DID, stewardAcmeDIDInfo.VerKey, "Acme", acmeWallet, acmeStewardDIDInfo, role.TrustAnchor)
	assert.NoError(t, err)

	fmt.Println("==============================")
	fmt.Println("== Getting Trust Anchor credentials - Thrift Onboarding  ==")
	fmt.Println("------------------------------")

	thriftWallet, stewardThriftDIDInfo, thriftStewardDIDInfo, _, err := onboarding(p, "Sovrin Steward", stewardWallet, stewardDIDInfo.DID, "Thrift", nil, "thrift_wallet")
	assert.NoError(t, err)

	fmt.Println("==============================")
	fmt.Println("== Getting Trust Anchor credentials - Thrift getting Verinym  ==")
	fmt.Println("------------------------------")

	thriftDIDInfo, err := getVerinym(t, p, "Sovrin Steward", stewardWallet, stewardDIDInfo.DID, stewardThriftDIDInfo.VerKey, "Thrift", thriftWallet, thriftStewardDIDInfo, role.TrustAnchor)
	assert.NoError(t, err)

	fmt.Println("==============================")
	fmt.Println("=== Credential Schemas Setup ==")
	fmt.Println("------------------------------")

	fmt.Println(`"Government" -> Create "Job-Certificate" Schema`)

	jobCertificateSchemaID, jobCertificateSchema, err := anoncreds.IssuerCreateSchema(governmentDIDInfo.DID, "Job-Certificate", "0.2", `["first_name", "last_name", "salary", "employee_status","experience"]`)
	assert.NoError(t, err)

	fmt.Println(`"Government" -> Send "Job-Certificate" Schema to Ledger`)
	_, err = sendSchema(p, governmentWallet, governmentDIDInfo.DID, jobCertificateSchema)
	assert.NoError(t, err)

	fmt.Println(`"Government" -> Create "Transcript" Schema`)
	transcriptSchemaID, transcriptSchema, err := anoncreds.IssuerCreateSchema(governmentDIDInfo.DID, "Transcript", "1.2", `["first_name","last_name","degree","status","year","average","ssn"]`)
	assert.NoError(t, err)

	fmt.Println(`"Government" -> Send "Transcript" Schema to Ledger`)
	_, err = sendSchema(p, governmentWallet, governmentDIDInfo.DID, transcriptSchema)
	assert.NoError(t, err)

	fmt.Println("==============================")
	fmt.Println("=== Faber Credential Definition Setup ==")
	fmt.Println("------------------------------")

	fmt.Println(`"Faber" -> Get "Transcript" Schema from Ledger`)
	_, transcriptSchema, err = getSchema(p, faberDIDInfo.DID, transcriptSchemaID)
	assert.NoError(t, err)

	fmt.Println(`"Faber" -> Create and store in Wallet "Faber Transcript" Credential Definition`)
	faberTranscriptCredDefID, faberTranscriptCredDefJSON, err := anoncreds.IssuerCreateAndStoreCredentialDef(faberWallet, faberDIDInfo.DID, transcriptSchema, "TAG1", "CL", `{"support_revocation": false}`)
	assert.NoError(t, err)

	fmt.Println(`"Faber" -> Send  "Faber Transcript" Credential Definition to Ledger`)
	_, err = sendCredDef(p, faberWallet, faberDIDInfo.DID, faberTranscriptCredDefJSON)
	assert.NoError(t, err)

	fmt.Println("==============================")
	fmt.Println("=== Acme Credential Definition Setup ==")
	fmt.Println("------------------------------")

	fmt.Println(`"Acme" ->  Get from Ledger "Job-Certificate" Schema`)
	_, jobCertificateSchema, err = getSchema(p, acmeDIDInfo.DID, jobCertificateSchemaID)
	assert.NoError(t, err)

	fmt.Println(`"Acme" -> Create and store in Wallet "Acme Job-Certificate" Credential Definition`)
	acmeJobCertificateCredDefID, acmeJobCertificateCredDefJSON, err := anoncreds.IssuerCreateAndStoreCredentialDef(acmeWallet, acmeDIDInfo.DID, jobCertificateSchema, "TAG1", "CL", `{"support_revocation": false}`)
	assert.NoError(t, err)

	fmt.Println(`"Acme" -> Send "Acme Job-Certificate" Credential Definition to Ledger`)
	_, err = sendCredDef(p, acmeWallet, acmeDIDInfo.DID, acmeJobCertificateCredDefJSON)
	assert.NoError(t, err)

	fmt.Println("==============================")
	fmt.Println("=== Getting Transcript with Faber ==")
	fmt.Println("==============================")
	fmt.Println("== Getting Transcript with Faber - Onboarding ==")
	fmt.Println("------------------------------")

	aliceWallet, faberAliceDIDInfo, aliceFaberDIDInfo, faberAliceConnectionResponseJSON, err := onboarding(p, "Faber", faberWallet, faberDIDInfo.DID, "Alice", nil, "alice_wallet")
	assert.NoError(t, err)
	faberAliceConnectionResponse := json.AsMap(faberAliceConnectionResponseJSON)

	fmt.Println("==============================")
	fmt.Println("== Getting Transcript with Faber - Getting Transcript Credential ==")
	fmt.Println("------------------------------")

	fmt.Println(`"Faber" -> Create "Transcript" Credential Offer for Alice`)
	transcriptCredOfferJSON, err := anoncreds.IssuerCreateCredentialOffer(faberWallet, faberTranscriptCredDefID)
	assert.NoError(t, err)

	fmt.Println(`"Faber" -> Get key for Alice did`)
	aliceFaberVerKey, err := did.KeyForDID(p, acmeWallet, faberAliceConnectionResponse["did"].String())
	assert.NoErrorf(t, err, "Error received from KeyForDID - DID [%s], Wallet [%s]: %s", faberAliceConnectionResponse["did"].String(), acmeWallet.Name)

	fmt.Println(`"Faber" -> Authcrypt "Transcript" Credential Offer for Alice`)
	authCryptedTranscriptCredOffer, err := crypto.AuthCrypt(faberWallet, faberAliceDIDInfo.VerKey, aliceFaberVerKey, []byte(transcriptCredOfferJSON))
	assert.NoError(t, err)
	fmt.Println(`"Faber" -> Send authcrypted "Transcript" Credential Offer to Alice`)

	fmt.Println(`"Alice" -> Authdecrypted "Transcript" Credential Offer from Faber`)
	faberAliceVerKey, authDecryptedTranscriptCredOfferJSON, err := crypto.AuthDecrypt(aliceWallet, aliceFaberDIDInfo.VerKey, authCryptedTranscriptCredOffer)
	assert.NoError(t, err)

	authDecryptedTranscriptCredOffer := json.AsMap(string(authDecryptedTranscriptCredOfferJSON))

	fmt.Println(`"Alice" -> Create and store "Alice" Master Secret in Wallet`)
	aliceMasterSecretID, err := anoncreds.ProverCreateMasterSecret(aliceWallet, "")
	assert.NoError(t, err)

	fmt.Println(`"Alice" -> Get "Faber Transcript" Credential Definition from Ledger`)

	faberTranscriptCredDefID, faberTranscriptCredDef, err := getCredDef(p, aliceFaberDIDInfo.DID, authDecryptedTranscriptCredOffer["cred_def_id"].String())
	assert.NoError(t, err)

	fmt.Println(`"Alice" -> Create "Transcript" Credential Request for Faber`)
	transcriptCredRequestJSON, transcriptCredRequestMetadataJSON, err := anoncreds.ProverCreateCredentialReq(aliceWallet, aliceFaberDIDInfo.DID, string(authDecryptedTranscriptCredOfferJSON), faberTranscriptCredDef, aliceMasterSecretID)
	assert.NoError(t, err)

	fmt.Println(`"Alice" -> Authcrypt "Transcript" Credential Request for Faber`)
	authCryptedTranscriptCredRequest, err := crypto.AuthCrypt(aliceWallet, aliceFaberDIDInfo.VerKey, faberAliceVerKey, []byte(transcriptCredRequestJSON))
	assert.NoError(t, err)

	fmt.Println(`"Alice" -> Send authcrypted "Transcript" Credential Request to Faber`)

	fmt.Println(`"Faber" -> Authdecrypt "Transcript" Credential Request from Alice`)
	aliceFaberVerKey, authDecryptedTranscriptCredRequestJSON, err := crypto.AuthDecrypt(faberWallet, faberAliceDIDInfo.VerKey, authCryptedTranscriptCredRequest)
	assert.NoError(t, err)

	fmt.Println(`"Faber" -> Create "Transcript" Credential for Alice`)

	degree := "Bachelor of Science, Marketing"
	status := "graduated"
	ssn := "123-45-6789"
	firstName := "Alice"
	lastName := "Garcia"

	transcriptCredValuesJSON := json.Map{
		"first_name": json.Map{
			"raw":     json.Str(firstName),
			"encoded": json.Str("1139481716457488690172217916278103335"),
		}.M(),
		"last_name": json.Map{
			"raw":     json.Str(lastName),
			"encoded": json.Str("5321642780241790123587902456789123452"),
		}.M(),
		"degree": json.Map{
			"raw":     json.Str(degree),
			"encoded": json.Str("12434523576212321"),
		}.M(),
		"status": json.Map{
			"raw":     json.Str(status),
			"encoded": json.Str("2213454313412354"),
		}.M(),
		"ssn": json.Map{
			"raw":     json.Str(ssn),
			"encoded": json.Str("3124141231422543541"),
		}.M(),
		"year": json.Map{
			"raw":     json.Str("2015"),
			"encoded": json.Str("2015"),
		}.M(),
		"average": json.Map{
			"raw":     json.Str("5"),
			"encoded": json.Str("5"),
		}.M(),
	}.JSON()

	transcriptCredJSON, _, _, err := anoncreds.IssuerCreateCredential(faberWallet, transcriptCredOfferJSON, string(authDecryptedTranscriptCredRequestJSON), transcriptCredValuesJSON, "", 0)
	assert.NoError(t, err)

	fmt.Println(`"Faber" -> Authcrypt "Transcript" Credential for Alice`)
	authCryptedTranscriptCredJSON, err := crypto.AuthCrypt(faberWallet, faberAliceDIDInfo.VerKey, aliceFaberVerKey, []byte(transcriptCredJSON))
	assert.NoError(t, err)

	fmt.Println(`"Faber" -> Send authcrypted "Transcript" Credential to Alice`)

	fmt.Println(`"Alice" -> Authdecrypted "Transcript" Credential from Faber`)
	_, authDecryptedTranscriptCredJSON, err := crypto.AuthDecrypt(aliceWallet, aliceFaberDIDInfo.VerKey, authCryptedTranscriptCredJSON)
	assert.NoError(t, err)

	fmt.Println(`"Alice" -> Store "Transcript" Credential from Faber`)
	_, err = anoncreds.ProverStoreCredential(aliceWallet, "", transcriptCredRequestMetadataJSON, string(authDecryptedTranscriptCredJSON), faberTranscriptCredDef, "")
	assert.NoError(t, err)

	fmt.Println("==============================")
	fmt.Println("=== Apply for the job with Acme ==")
	fmt.Println("==============================")
	fmt.Println("== Apply for the job with Acme - Onboarding ==")
	fmt.Println("------------------------------")

	aliceWallet, acmeAliceDIDInfo, aliceAcmeDIDInfo, acmeAliceConnectionResponseJSON, err := onboarding(p, "Acme", acmeWallet, acmeDIDInfo.DID, "Alice", aliceWallet, "alice_wallet")
	assert.NoError(t, err)

	acmeAliceConnectionResponse := json.AsMap(acmeAliceConnectionResponseJSON)

	fmt.Println("==============================")
	fmt.Println("== Apply for the job with Acme - Transcript proving ==")
	fmt.Println("------------------------------")

	fmt.Println(`"Acme" -> Create "Job-Application" Proof Request`)

	jobApplicationProofRequestJSON := json.Map{
		"nonce":   json.Str("1432422343242122312411212"),
		"name":    json.Str("Job-Application"),
		"version": json.Str("0.1"),
		"requested_attributes": json.Map{
			"attr1_referent": json.Map{
				"name": json.Str("first_name"),
			}.M(),
			"attr2_referent": json.Map{
				"name": json.Str("last_name"),
			}.M(),
			"attr3_referent": json.Map{
				"name": json.Str("degree"),
				"restrictions": json.Slice{
					json.Map{
						"cred_def_id": json.Str(faberTranscriptCredDefID),
					}.M(),
				}.M(),
			}.M(),
			"attr4_referent": json.Map{
				"name": json.Str("status"),
				"restrictions": json.Slice{
					json.Map{
						"cred_def_id": json.Str(faberTranscriptCredDefID),
					}.M(),
				}.M(),
			}.M(),
			"attr5_referent": json.Map{
				"name": json.Str("ssn"),
				"restrictions": json.Slice{
					json.Map{
						"cred_def_id": json.Str(faberTranscriptCredDefID),
					}.M(),
				}.M(),
			}.M(),
			"attr6_referent": json.Map{
				"name": json.Str("phone_number"),
			}.M(),
		}.M(),
		"requested_predicates": json.Map{
			"predicate1_referent": json.Map{
				"name":    json.Str("average"),
				"p_type":  json.Str(">="),
				"p_value": json.Int(4),
				"restrictions": json.Slice{
					json.Map{
						"cred_def_id": json.Str(faberTranscriptCredDefID),
					}.M(),
				}.M(),
			}.M(),
		}.M(),
	}.JSON()

	fmt.Println(`"Acme" -> Get key for Alice did`)
	aliceAcmeVerKey, err := did.KeyForDID(p, acmeWallet, acmeAliceConnectionResponse["did"].String())
	assert.NoError(t, err)

	fmt.Println(`"Acme" -> Authcrypt "Job-Application" Proof Request for Alice`)
	authCryptedJobApplicationProofRequestJSON, err := crypto.AuthCrypt(acmeWallet, acmeAliceDIDInfo.VerKey, aliceAcmeVerKey, []byte(jobApplicationProofRequestJSON))
	assert.NoError(t, err)

	fmt.Println(`"Acme" -> Send authcrypted "Job-Application" Proof Request to Alice`)

	fmt.Println(`"Alice" -> Authdecrypt "Job-Application" Proof Request from Acme`)
	acmeAliceVerKey, authDecryptedJobApplicationProofRequestJSON, err := crypto.AuthDecrypt(aliceWallet, aliceAcmeDIDInfo.VerKey, authCryptedJobApplicationProofRequestJSON)
	assert.NoError(t, err)

	fmt.Println(`"Alice" -> Get credentials for "Job-Application" Proof Request`)
	credsForJobApplicationProofRequestJSON, err := anoncreds.ProverGetCredentialsForProofReq(aliceWallet, string(authDecryptedJobApplicationProofRequestJSON))
	assert.NoError(t, err)

	credsForJobApplicationProofRequestAttrs := json.AsMap(credsForJobApplicationProofRequestJSON)["attrs"].AsMap()
	credsForJobApplicationProofRequestPredicates := json.AsMap(credsForJobApplicationProofRequestJSON)["predicates"].AsMap()
	credForAttr1 := credsForJobApplicationProofRequestAttrs["attr1_referent"].Idx(0).Val("cred_info")
	credForAttr2 := credsForJobApplicationProofRequestAttrs["attr2_referent"].Idx(0).Val("cred_info")
	credForAttr3 := credsForJobApplicationProofRequestAttrs["attr3_referent"].Idx(0).Val("cred_info")
	credForAttr4 := credsForJobApplicationProofRequestAttrs["attr4_referent"].Idx(0).Val("cred_info")
	credForAttr5 := credsForJobApplicationProofRequestAttrs["attr5_referent"].Idx(0).Val("cred_info")
	credForPredicate1 := credsForJobApplicationProofRequestPredicates["predicate1_referent"].Idx(0).Val("cred_info")

	credsForJobApplicationProof := make(json.Map)
	credsForJobApplicationProof[credForAttr1.Val("referent").String()] = credForAttr1
	credsForJobApplicationProof[credForAttr2.Val("referent").String()] = credForAttr2
	credsForJobApplicationProof[credForAttr3.Val("referent").String()] = credForAttr3
	credsForJobApplicationProof[credForAttr4.Val("referent").String()] = credForAttr4
	credsForJobApplicationProof[credForAttr5.Val("referent").String()] = credForAttr5
	credsForJobApplicationProof[credForPredicate1.Val("referent").String()] = credForPredicate1

	schemasJSON, credDefsJSON, revocStatesJSON, err := proverGetEntitiesFromLedger(p, aliceFaberDIDInfo.DID, credsForJobApplicationProof, "Alice")
	assert.NoError(t, err)

	fmt.Println(`"Alice" -> Create "Job-Application" Proof`)
	jobApplicationRequestedCredsJSON := json.Map{
		"self_attested_attributes": json.Map{
			"attr1_referent": json.Str("Alice"),
			"attr2_referent": json.Str("Garcia"),
			"attr6_referent": json.Str("123-45-6789"),
		}.M(),
		"requested_attributes": json.Map{
			"attr3_referent": json.Map{
				"cred_id":  credForAttr3.Val("referent"),
				"revealed": json.Bool(true)}.M(),
			"attr4_referent": json.Map{
				"cred_id":  credForAttr4.Val("referent"),
				"revealed": json.Bool(true)}.M(),
			"attr5_referent": json.Map{
				"cred_id":  credForAttr5.Val("referent"),
				"revealed": json.Bool(true)}.M(),
		}.M(),
		"requested_predicates": json.Map{
			"predicate1_referent": json.Map{
				"cred_id": credForPredicate1.Val("referent")}.M(),
		}.M(),
	}.JSON()

	jobApplicationProofJSON, err := anoncreds.ProverCreateProof(aliceWallet, string(authDecryptedJobApplicationProofRequestJSON), jobApplicationRequestedCredsJSON, aliceMasterSecretID, schemasJSON, credDefsJSON, revocStatesJSON)
	assert.NoError(t, err)

	fmt.Printf(`"Alice" -> Authcrypt "Job-Application" Proof for Acme` + "\n")
	authCryptedJobApplicationProofJSON, err := crypto.AuthCrypt(aliceWallet, aliceAcmeDIDInfo.VerKey, acmeAliceDIDInfo.VerKey, []byte(jobApplicationProofJSON))
	assert.NoError(t, err)

	fmt.Printf(`"Alice" -> Send authcrypted "Job-Application" Proof to Acme` + "\n")

	fmt.Printf(`"Acme" -> Authdecrypted "Job-Application" Proof from Alice` + "\n")
	_, decryptedJobApplicationProofJSON, err := crypto.AuthDecrypt(acmeWallet, acmeAliceDIDInfo.VerKey, authCryptedJobApplicationProofJSON)
	assert.NoError(t, err)

	decryptedJobApplicationProof := json.AsMap(string(decryptedJobApplicationProofJSON))
	schemasJSON, credDefsJSON, revocRefDefsJSON, revocRegsJSON, err := verifierGetEntitiesFromLedger(p, acmeDIDInfo.DID, decryptedJobApplicationProof["identifiers"].Slice(), "Acme")
	assert.NoError(t, err)

	fmt.Printf(`"Acme" -> Verify "Job-Application" Proof from Alice` + "\n")

	requestedProof := decryptedJobApplicationProof["requested_proof"]
	revealedAttrs := requestedProof.Val("revealed_attrs").AsMap()
	assert.Equalf(t, degree, revealedAttrs["attr3_referent"].Val("raw").String(), "Invalid value for degree")
	assert.Equalf(t, status, revealedAttrs["attr4_referent"].Val("raw").String(), "Invalid value for status")
	assert.Equalf(t, ssn, revealedAttrs["attr5_referent"].Val("raw").String(), "Invalid value for ssn")

	selfAttestedAttrs := requestedProof.Val("self_attested_attrs").AsMap()
	assert.Equalf(t, firstName, selfAttestedAttrs["attr1_referent"].String(), "Invalid value for firstName")
	assert.Equalf(t, lastName, selfAttestedAttrs["attr2_referent"].String(), "Invalid value for lastName")
	assert.Equalf(t, ssn, selfAttestedAttrs["attr6_referent"].String(), "Invalid value for ssn")

	valid, err := anoncreds.VerifierVerifyProof(string(jobApplicationProofRequestJSON), string(decryptedJobApplicationProofJSON), schemasJSON, credDefsJSON, revocRefDefsJSON, revocRegsJSON)
	assert.NoError(t, err)
	assert.Equalf(t, true, valid, `"Acme" -> "Job-Application" Proof from Alice is NOT valid!`)
	fmt.Printf(`"Acme" -> "Job-Application" Proof from Alice is valid!` + "\n")

	fmt.Printf("==============================\n")
	fmt.Printf("== Apply for the job with Acme - Getting Job-Certificate Credential ==\n")
	fmt.Printf("------------------------------\n")

	fmt.Printf(`"Acme" -> Create "Job-Certificate" Credential Offer for Alice`)
	jobCertificateCredOfferJSON, err := anoncreds.IssuerCreateCredentialOffer(acmeWallet, acmeJobCertificateCredDefID)
	assert.NoError(t, err)

	fmt.Printf(`"Acme" -> Get key for Alice did` + "\n")
	aliceAcmeVerKey, err = did.KeyForDID(p, acmeWallet, acmeAliceConnectionResponse["did"].String())
	assert.NoError(t, err)

	fmt.Printf(`"Acme" -> Authcrypt "Job-Certificate" Credential Offer for Alice` + "\n")
	authCryptedJobCertificateCredOffer, err := crypto.AuthCrypt(acmeWallet, acmeAliceDIDInfo.VerKey, aliceAcmeVerKey, []byte(jobCertificateCredOfferJSON))
	assert.NoError(t, err)

	fmt.Printf(`"Acme" -> Send authcrypted "Job-Certificate" Credential Offer to Alice` + "\n")

	fmt.Printf(`"Alice" -> Authdecrypted "Job-Certificate" Credential Offer from Acme` + "\n")
	acmeAliceVerKey, authDecryptedJobCertificateCredOfferJSON, err := crypto.AuthDecrypt(aliceWallet, aliceAcmeDIDInfo.VerKey, authCryptedJobCertificateCredOffer)
	assert.NoError(t, err)

	fmt.Printf(`"Alice" -> Get "Acme Job-Certificate" Credential Definition from Ledger` + "\n")
	authDecryptedJobCertificateCredOffer := json.AsMap(string(authDecryptedJobCertificateCredOfferJSON))
	_, acmeJobCertificateCredDef, err := getCredDef(p, aliceAcmeDIDInfo.DID, authDecryptedJobCertificateCredOffer["cred_def_id"].String())
	assert.NoError(t, err)

	fmt.Printf(`"Alice" -> Create "Job-Certificate" Credential Request for Acme` + "\n")
	jobCertificateCredRequestJSON, jobCertificateCredRequestMetadataJSON, err := anoncreds.ProverCreateCredentialReq(aliceWallet, aliceAcmeDIDInfo.DID, string(authDecryptedJobCertificateCredOfferJSON), acmeJobCertificateCredDef, aliceMasterSecretID)
	assert.NoError(t, err)

	fmt.Printf(`"Alice" -> Authcrypt "Job-Certificate" Credential Request for Acme` + "\n")
	authCryptedJobCertificateCredRequestJSON, err := crypto.AuthCrypt(aliceWallet, aliceAcmeDIDInfo.VerKey, acmeAliceVerKey, []byte(jobCertificateCredRequestJSON))
	assert.NoError(t, err)

	fmt.Printf(`"Alice" -> Send authcrypted "Job-Certificate" Credential Request to Acme` + "\n")

	fmt.Printf(`"Acme" -> Authdecrypt "Job-Certificate" Credential Request from Alice`)
	aliceAcmeVerKey, authDecryptedJobCertificateCredRequestJSON, err := crypto.AuthDecrypt(acmeWallet, acmeAliceDIDInfo.VerKey, authCryptedJobCertificateCredRequestJSON)
	assert.NoError(t, err)

	fmt.Printf(`"Acme" -> Create "Job-Certificate" Credential for Alice` + "\n")
	aliceJobCertificateCredValuesJSON := json.Map{
		"first_name": json.Map{
			"raw":     json.Str("Alice"),
			"encoded": json.Str("245712572474217942457235975012103335"),
		}.M(),
		"last_name": json.Map{
			"raw":     json.Str("Garcia"),
			"encoded": json.Str("312643218496194691632153761283356127"),
		}.M(),
		"employee_status": json.Map{
			"raw":     json.Str("Permanent"),
			"encoded": json.Str("2143135425425143112321314321"),
		}.M(),
		"salary": json.Map{
			"raw":     json.Str("2400"),
			"encoded": json.Str("2400"),
		}.M(),
		"experience": json.Map{
			"raw":     json.Str("10"),
			"encoded": json.Str("10"),
		}.M(),
	}.JSON()

	jobCertificateCredJSON, _, _, err := anoncreds.IssuerCreateCredential(acmeWallet, jobCertificateCredOfferJSON, string(authDecryptedJobCertificateCredRequestJSON), aliceJobCertificateCredValuesJSON, "", 0)
	assert.NoError(t, err)

	fmt.Printf(`"Acme" ->  Authcrypt "Job-Certificate" Credential for Alice` + "\n")
	authCryptedJobCertificateCredRequestJSON, err = crypto.AuthCrypt(acmeWallet, acmeAliceDIDInfo.VerKey, aliceAcmeVerKey, []byte(jobCertificateCredJSON))
	assert.NoError(t, err)

	fmt.Printf(`"Acme" ->  Send authcrypted "Job-Certificate" Credential to Alice` + "\n")

	fmt.Printf(`"Alice" -> Authdecrypted "Job-Certificate" Credential from Acme` + "\n")
	_, authDecryptedJobCertificateCredJSON, err := crypto.AuthDecrypt(aliceWallet, aliceAcmeDIDInfo.VerKey, authCryptedJobCertificateCredRequestJSON)
	assert.NoError(t, err)

	fmt.Printf(`"Alice" -> Store "Job-Certificate" Credential` + "\n")
	_, err = anoncreds.ProverStoreCredential(aliceWallet, "", jobCertificateCredRequestMetadataJSON, string(authDecryptedJobCertificateCredJSON), acmeJobCertificateCredDefJSON, "")
	assert.NoError(t, err)

	fmt.Printf("==============================\n")
	fmt.Printf("=== Apply for the loan with Thrift ==\n")
	fmt.Printf("==============================\n")
	fmt.Printf("== Apply for the loan with Thrift - Onboarding ==\n")
	fmt.Printf("------------------------------\n")

	aliceWallet, thriftAliceDID, aliceThriftDID, thriftAliceConnectionResponseJSON, err := onboarding(p, "Thrift", thriftWallet, thriftDIDInfo.DID, "Alice", aliceWallet, "alice_wallet")
	assert.NoError(t, err)

	fmt.Printf("==============================\n")
	fmt.Printf("== Apply for the loan with Thrift - Transcript and Job-Certificate proving  ==\n")
	fmt.Printf("------------------------------\n")

	fmt.Printf(`"Thrift" -> Create "Loan-Application-KYC" Proof Request` + "\n")
	applyLoanKYCProofRequestJSON := json.Map{
		"nonce":   json.Str("123432421212"),
		"name":    json.Str("Loan-Application-KYC"),
		"version": json.Str("0.1"),
		"requested_attributes": json.Map{
			"attr1_referent": json.Map{
				"name": json.Str("first_name"),
			}.M(),
			"attr2_referent": json.Map{
				"name": json.Str("last_name"),
			}.M(),
			"attr3_referent": json.Map{
				"name": json.Str("ssn"),
			}.M(),
		}.M(),
		"requested_predicates": json.Map{}.M(),
	}.JSON()

	fmt.Printf(`"Thrift" -> Get key for Alice did` + "\n")
	thriftAliceConnectionResponse := json.AsMap(thriftAliceConnectionResponseJSON)
	aliceThriftVerKey, err := did.KeyForDID(p, thriftWallet, thriftAliceConnectionResponse["did"].String())
	assert.NoError(t, err)

	fmt.Printf(`"Thrift" -> Authcrypt "Loan-Application-KYC" Proof Request for Alice` + "\n")
	authCryptedApplyLoanKYCProofRequestJSON, err := crypto.AuthCrypt(thriftWallet, thriftAliceDID.VerKey, aliceThriftVerKey, []byte(applyLoanKYCProofRequestJSON))
	assert.NoError(t, err)

	fmt.Printf(`"Thrift" -> Send authcrypted "Loan-Application-KYC" Proof Request to Alice` + "\n")

	fmt.Printf(`"Alice" -> Authdecrypt "Loan-Application-KYC" Proof Request from Thrift` + "\n")
	thriftAliceVerKey, authDecryptedApplyLoanKYCProofRequestJSON, err := crypto.AuthDecrypt(aliceWallet, aliceThriftDID.VerKey, authCryptedApplyLoanKYCProofRequestJSON)
	assert.NoError(t, err)

	fmt.Printf(`"Alice" -> Get credentials for "Loan-Application-KYC" Proof Request` + "\n")
	credsJSONForApplyLoanKYCProofRequestJSON, err := anoncreds.ProverGetCredentialsForProofReq(aliceWallet, string(authDecryptedApplyLoanKYCProofRequestJSON))
	assert.NoError(t, err)

	credsForApplyLoanKYCProofRequestAttrs := json.AsMap(credsJSONForApplyLoanKYCProofRequestJSON)["attrs"].AsMap()

	credForAttr1 = credsForApplyLoanKYCProofRequestAttrs["attr1_referent"].Idx(0).Val("cred_info")
	credForAttr2 = credsForApplyLoanKYCProofRequestAttrs["attr2_referent"].Idx(0).Val("cred_info")
	credForAttr3 = credsForApplyLoanKYCProofRequestAttrs["attr3_referent"].Idx(0).Val("cred_info")

	referentForAttr1 := credForAttr1.Val("referent").String()
	referentForAttr2 := credForAttr2.Val("referent").String()
	referentForAttr3 := credForAttr3.Val("referent").String()

	credsForApplyLoanKYCProof := make(json.Map)
	credsForApplyLoanKYCProof[referentForAttr1] = credForAttr1
	credsForApplyLoanKYCProof[referentForAttr2] = credForAttr2
	credsForApplyLoanKYCProof[referentForAttr3] = credForAttr3

	schemasJSON, credDefsJSON, revocStatesJSON, err = proverGetEntitiesFromLedger(p, aliceThriftDID.DID, credsForApplyLoanKYCProof, "Alice")
	assert.NoError(t, err)

	fmt.Printf(`"Alice" -> Create "Loan-Application-KYC" Proof` + "\n")

	applyLoanKYCRequestedCredsJSON := json.Map{
		"self_attested_attributes": json.Map{}.M(),
		"requested_attributes": json.Map{
			"attr1_referent": json.Map{
				"cred_id":  json.Str(referentForAttr1),
				"revealed": json.Bool(true),
			}.M(),
			"attr2_referent": json.Map{
				"cred_id":  json.Str(referentForAttr2),
				"revealed": json.Bool(true),
			}.M(),
			"attr3_referent": json.Map{
				"cred_id":  json.Str(referentForAttr3),
				"revealed": json.Bool(true),
			}.M(),
		}.M(),
		"requested_predicates": json.Map{}.M(),
	}.JSON()

	aliceApplyLoanKYCProofJSON, err := anoncreds.ProverCreateProof(aliceWallet, string(authDecryptedApplyLoanKYCProofRequestJSON), applyLoanKYCRequestedCredsJSON, aliceMasterSecretID, schemasJSON, credDefsJSON, revocStatesJSON)
	assert.NoError(t, err)

	fmt.Printf(`"Alice" -> Authcrypt "Loan-Application-KYC" Proof for Thrift` + "\n")
	authCryptedAliceApplyLoanKYCProofJSON, err := crypto.AuthCrypt(aliceWallet, aliceThriftDID.VerKey, thriftAliceVerKey, []byte(aliceApplyLoanKYCProofJSON))
	assert.NoError(t, err)

	fmt.Printf(`"Alice" -> Send authcrypted "Loan-Application-KYC" Proof to Thrift` + "\n")

	fmt.Printf(`"Thrift" -> Authdecrypted "Loan-Application-KYC" Proof from Alice` + "\n")
	_, authDecryptedAliceApplyLoanKYCProofJSON, err := crypto.AuthDecrypt(thriftWallet, thriftAliceDID.VerKey, authCryptedAliceApplyLoanKYCProofJSON)
	assert.NoError(t, err)

	authDecryptedAliceApplyLoanKYCProof := json.AsMap(string(authDecryptedAliceApplyLoanKYCProofJSON))

	fmt.Printf(`"Thrift" -> Get Schemas, Credential Definitions and Revocation Registries from Ledger required for Proof verifying` + "\n")

	fmt.Printf(`"Thrift" -> Verify "Loan-Application-KYC" Proof from Alice` + "\n")

	authDecryptedAliceApplyLoanKYCProofRevealedAttrs := authDecryptedAliceApplyLoanKYCProof["requested_proof"].Val("revealed_attrs").AsMap()
	assert.Equalf(t, firstName, authDecryptedAliceApplyLoanKYCProofRevealedAttrs["attr1_referent"].Val("raw").String(), "Invalid value for firstName")
	assert.Equalf(t, lastName, authDecryptedAliceApplyLoanKYCProofRevealedAttrs["attr2_referent"].Val("raw").String(), "Invalid value for lastName")
	assert.Equalf(t, ssn, authDecryptedAliceApplyLoanKYCProofRevealedAttrs["attr3_referent"].Val("raw").String(), "Invalid value for lastName")

	schemasJSON, credDefsJSON, revocDefsJSON, revocRegsJSON, err := verifierGetEntitiesFromLedger(p, thriftDIDInfo.DID, authDecryptedAliceApplyLoanKYCProof["identifiers"].Slice(), "Thrift")
	assert.NoError(t, err)

	valid, err = anoncreds.VerifierVerifyProof(string(applyLoanKYCProofRequestJSON), string(authDecryptedAliceApplyLoanKYCProofJSON), schemasJSON, credDefsJSON, revocDefsJSON, revocRegsJSON)
	assert.NoError(t, err)
	assert.Equalf(t, true, valid, `"Thrift" -> Verify "Loan-Application-KYC" Proof from Alice is NOT valid`)
	fmt.Printf(`"Thrift" -> Verify "Loan-Application-KYC" Proof from Alice is valid!` + "\n")

	fmt.Printf("==============================" + "\n")

	fmt.Printf(`"Sovrin Steward" -> Close and Delete wallet` + "\n")
	stewardWallet.Close()
	wallet.Delete(stewardWallet.Name, "")

	fmt.Printf(`"Government" -> Close and Delete wallet` + "\n")
	governmentWallet.Close()
	wallet.Delete(governmentWallet.Name, "")

	fmt.Printf(`"Faber" -> Close and Delete wallet` + "\n")
	faberWallet.Close()
	wallet.Delete(faberWallet.Name, "")

	fmt.Printf(`"Acme" -> Close and Delete wallet` + "\n")
	acmeWallet.Close()
	wallet.Delete(acmeWallet.Name, "")

	fmt.Printf(`"Thrift" -> Close and Delete wallet` + "\n")
	thriftWallet.Close()
	wallet.Delete(thriftWallet.Name, "")

	fmt.Printf(`"Alice" -> Close and Delete wallet` + "\n")
	aliceWallet.Close()
	wallet.Delete(aliceWallet.Name, "")

	fmt.Printf(`Close and Delete pool` + "\n")
	p.Close()
	pool.Delete(p.Name)

	fmt.Printf(`Getting started -> done` + "\n")
}

func onboarding(
	pool *pool.Pool,
	from string, fromWallet *wallet.Wallet, fromDID string,
	to string, toWallet *wallet.Wallet, toWalletName string) (toWalletRet *wallet.Wallet, fromToDIDInfo *did.Info, toFromDID *did.Info, decryptedConnectionResponseJSON string, err error) {

	fmt.Printf(`"%s" -> Create and store in Wallet "{%s} {%s}" DID`+"\n", from, from, to)

	fromToDIDInfo, err = did.CreateAndStoreMyDID(fromWallet, "{}")
	if err != nil {
		err = fmt.Errorf("error creating and storing DID - Wallet [%s]: %s", fromWallet.Name, err)
		return
	}

	fmt.Printf(`"%s" -> Send Nym to Ledger for "{%s} {%s}" DID`+"\n", from, from, to)

	_, err = sendNYM(pool, fromWallet, fromDID, fromToDIDInfo.DID, fromToDIDInfo.VerKey, nil)
	if err != nil {
		err = fmt.Errorf("error sending NYM - Wallet [%s], FromDID [%s], FromToDID [%s], FromToVerKey [%s]: %s", fromWallet.Name, fromDID, fromToDIDInfo.DID, fromToDIDInfo.VerKey, err)
		return
	}

	toWalletRet = toWallet
	if toWalletRet == nil {
		fmt.Printf(`"%s" -> Create wallet`+"\n", to)
		toWalletRet, err = getWallet(pool.Name, toWalletName)
		if err != nil {
			err = fmt.Errorf("error getting wallet - Wallet [%s]: %s", toWalletName, err)
			return
		}
	}

	fmt.Printf(`"%s" -> Create and store in Wallet "%s %s" DID`+"\n", to, to, from)
	toFromDID, err = did.CreateAndStoreMyDID(toWalletRet, "{}")
	if err != nil {
		err = fmt.Errorf("error creating and storing DID - Wallet [%s]: %s", toWalletRet.Name, err)
		return
	}

	fmt.Printf(`"%s" -> Send connection request to %s with "%s %s" DID and nonce`+"\n", from, to, from, to)
	connectionRequest := json.Map{
		"did":   json.Str(fromToDIDInfo.DID),
		"nonce": json.Int(123456789),
	}

	// (If this were real then the connection request would encrypted and sent to the "to" endpoint)

	// The "to" endpoint receives the connection request and sends back an encrypted response...
	fmt.Printf(`"%s" -> Get key for did from "%s" connection request`+"\n", to, from)
	var fromToVerKey string
	fromToVerKey, err = did.KeyForDID(pool, toWalletRet, fromToDIDInfo.DID)
	if err != nil {
		err = fmt.Errorf("error getting key for DID - DID [%s], Wallet [%s]: %s", fromToDIDInfo.DID, toWalletRet.Name, err)
		return
	}
	fmt.Printf(`"%s" -> Anoncrypt connection response for "%s" with "%s %s" DID, verkey and nonce`+"\n", to, from, to, from)
	connectionResponseJSON := json.Map{
		"did":    json.Str(toFromDID.DID),
		"verKey": json.Str(toFromDID.VerKey),
		"nonce":  json.Int(connectionRequest["nonce"].Int64()),
	}.JSON()

	var anoncryptedConnectionResponse []byte
	anoncryptedConnectionResponse, err = crypto.AnonCrypt(fromToVerKey, []byte(connectionResponseJSON))
	if err != nil {
		err = fmt.Errorf("error encrypting connectionResponse [%s]: %s", connectionResponseJSON, err)
		return
	}
	fmt.Printf(`"%s" -> Send anoncrypted connection response to "%s"`+"\n", to, from)

	// The "from" endpoint receives the connection response, decrypts it, and authenticates it (using the nonce)
	fmt.Printf(`"%s" -> Anondecrypt connection response from "%s"`+"\n", from, to)
	decryptedConnectionResponseJSONBytes, err := crypto.AnonDecrypt(fromWallet, fromToDIDInfo.VerKey, anoncryptedConnectionResponse)
	if err != nil {
		err = fmt.Errorf("error decrypting connectionResponse: %s", err)
		return
	}
	decryptedConnectionResponseJSON = string(decryptedConnectionResponseJSONBytes)
	decryptedConnectionResponse := json.AsMap(decryptedConnectionResponseJSON)

	fmt.Printf(`"%s" -> Authenticates "%s" by comparison of Nonce`+"\n", from, to)
	if decryptedConnectionResponse["nonce"].Int64() != connectionRequest["nonce"].Int64() {
		err = fmt.Errorf("Expecting decrypted connection response nonce to be %d but got %d", connectionRequest["nonce"].Int64(), decryptedConnectionResponse["nonce"].Int64())
		return
	}

	fmt.Printf(`"%s" -> Send Nym to Ledger for "%s %s" DID`+"\n", from, to, from)
	_, err = sendNYM(pool, fromWallet, fromDID, toFromDID.DID, toFromDID.VerKey, nil)
	if err != nil {
		err = fmt.Errorf("error sending NYM - FromWallet [%s], FromDID [%s], ToFromDID [%s], ToFromVerKey [%s]: %s", fromWallet.Name, fromDID, toFromDID.DID, toFromDID.VerKey, err)
		return
	}
	return
}

func sendNYM(pool *pool.Pool, wallet *wallet.Wallet, did string, newDID string, newKey string, role *role.Role) (string, error) {
	nymReq, err := ledger.BuildNYMRequest(did, newDID, newKey, nil, role)
	if err != nil {
		return "", fmt.Errorf("error building NYM request - NewDID [%s], NewKey [%s], role [%s]: %s", newDID, newKey, role, err)
	}
	return ledger.SignAndSubmitRequest(pool, wallet, did, nymReq)
}

func getWallet(poolName, walletName string) (*wallet.Wallet, error) {
	err := wallet.Create(poolName, walletName, "", "", "")
	if err != nil && indyerror.Code(err) != indyerror.WalletAlreadyExistsError {
		return nil, fmt.Errorf("error creating wallet - Wallet [%s]: %s", walletName, err)
	}
	return wallet.Open(walletName, "", "")
}

func getVerinym(t *testing.T, pool *pool.Pool, from string, fromWallet *wallet.Wallet, fromDID string, fromToKey string, to string, toWallet *wallet.Wallet, toFromDIDInfo *did.Info, role *role.Role) (*did.Info, error) {
	fmt.Printf(`"%s" -> Create and store in Wallet "%s" new DID`+"\n", to, to)
	toDID, err := did.CreateAndStoreMyDID(toWallet, "{}")
	if err != nil {
		return nil, fmt.Errorf("error creating and storing DID - Wallet [%s]: %s", fromWallet.Name, err)
	}

	fmt.Printf(`"%s" -> Authcrypt "%s DID info" for "%s"`+"\n", to, to, from)
	didInfoJSON := json.Map{
		"did":    json.Str(toDID.DID),
		"verkey": json.Str(toDID.VerKey),
	}.JSON()

	authCryptedDIDInfoJSON, err := crypto.AuthCrypt(toWallet, toFromDIDInfo.VerKey, fromToKey, []byte(didInfoJSON))
	if err != nil {
		return nil, fmt.Errorf("error encrypting DID Info - DIDInfo [%s]: %s", didInfoJSON, err)
	}

	fmt.Printf(`"%s" -> Send authcrypted "%s DID info" to %s`+"\n", to, to, from)
	fmt.Printf(`"%s" -> Authdecrypted "%s DID info" from %s`+"\n", from, to, to)

	senderVerkey, authDecryptedDIDInfoJSON, err := crypto.AuthDecrypt(fromWallet, fromToKey, authCryptedDIDInfoJSON)
	if err != nil {
		return nil, fmt.Errorf("error decrypting DID Info - CryptedDIDInfo [%s]: %s", authCryptedDIDInfoJSON, err)
	}

	authDecryptedDIDInfoMsg := json.AsMap(string(authDecryptedDIDInfoJSON))
	authDecryptedDIDInfo := &did.Info{
		DID:    authDecryptedDIDInfoMsg["did"].String(),
		VerKey: authDecryptedDIDInfoMsg["verkey"].String(),
	}

	fmt.Printf(`"%s" -> Authenticate %s by comparison of Verkeys`+"\n", from, to)

	key, err := did.KeyForDID(pool, fromWallet, toFromDIDInfo.DID)
	if err != nil {
		return nil, fmt.Errorf("error getting key for DID [%s] in wallet [%s]: %s", toFromDIDInfo.DID, fromWallet.Name, err)
	}
	if senderVerkey != key {
		return nil, fmt.Errorf("expecting sender key to be [%#x] but got [%#x]", senderVerkey, key)
	}

	fmt.Printf(`"%s" -> Send Nym to Ledger for "%s DID" with %s Role`+"\n", from, to, role)
	_, err = sendNYM(pool, fromWallet, fromDID, authDecryptedDIDInfo.DID, authDecryptedDIDInfo.VerKey, role)
	if err != nil {
		return nil, fmt.Errorf("error sending NYM - FromWallet [%s], fromDID [%s], DecryptedDIDInfo.DID [%s], DecryptedDIDInfo.VerKey [%s]: %s", fromWallet.Name, fromDID, authDecryptedDIDInfo.DID, authDecryptedDIDInfo.VerKey, err)
	}

	return toDID, nil
}

func sendSchema(pool *pool.Pool, wallet *wallet.Wallet, did string, schema string) (string, error) {
	schemaRequest, err := ledger.BuildSchemaRequest(did, schema)
	if err != nil {
		return "", fmt.Errorf("error building schema request - Schema [%s], DID [%s]: %s", schema, did, err)
	}
	return ledger.SignAndSubmitRequest(pool, wallet, did, schemaRequest)
}

func getSchema(pool *pool.Pool, did string, schemaID string) (id, schemaJSON string, err error) {
	var getSchemaRequest string
	getSchemaRequest, err = ledger.BuildGetSchemaRequest(did, schemaID)
	if err != nil {
		err = fmt.Errorf("error building get-schema request - SchemaID [%s], DID [%s]: %s", schemaID, did, err)
		return
	}

	var getSchemaResponse string
	getSchemaResponse, err = ledger.SubmitRequest(pool, getSchemaRequest)
	if err != nil {
		err = fmt.Errorf("error submitting get-schema request - GetSchemaRequest [%s]: %s", getSchemaRequest, err)
		return
	}

	id, schemaJSON, err = ledger.ParseGetSchemaResponse(getSchemaResponse)
	if err != nil {
		err = fmt.Errorf("error parsing get-schema response - GetSchemaResponse [%s]: %s", getSchemaResponse, err)
	}
	return
}

func getCredDef(pool *pool.Pool, did, schemaID string) (string, string, error) {
	getCredDefRequest, err := ledger.BuildGetCredDefRequest(did, schemaID)
	if err != nil {
		return "", "", err
	}
	getCredDefResponse, err := ledger.SubmitRequest(pool, getCredDefRequest)
	if err != nil {
		return "", "", err
	}
	return ledger.ParseGetCredDefResponse(getCredDefResponse)
}

func sendCredDef(pool *pool.Pool, wallet *wallet.Wallet, did, credDefJSON string) (string, error) {
	credDefRequest, err := ledger.BuildCredDefRequest(did, credDefJSON)
	if err != nil {
		return "", err
	}
	return ledger.SignAndSubmitRequest(pool, wallet, did, credDefRequest)
}

func proverGetEntitiesFromLedger(pool *pool.Pool, did string, identifiers json.Map, actor string) (schemasJSON string, credDefsJSON string, revStatesJSON string, err error) {
	schemas := make(json.Map)
	credDefs := make(json.Map)
	revStates := make(json.Map)

	for _, identifier := range identifiers {
		item := identifier.AsMap()
		fmt.Printf(`"%s" -> Get Schema from Ledger`+"\n", actor)
		receivedSchemaID, receivedSchemaJSON, err := getSchema(pool, did, item["schema_id"].String())
		if err != nil {
			return "", "", "", err
		}

		schemas[receivedSchemaID] = json.New(receivedSchemaJSON)

		fmt.Printf(`"%s" -> Get Credential Definition from Ledger`+"\n", actor)

		receivedCredDefID, receivedCredDefJSON, err := getCredDef(pool, did, item["cred_def_id"].String())
		if err != nil {
			return "", "", "", err
		}

		credDefs[receivedCredDefID] = json.New(receivedCredDefJSON)

		// TODO: Create Revocation States
	}

	schemasJSON = schemas.JSON()
	credDefsJSON = credDefs.JSON()
	revStatesJSON = revStates.JSON()

	return
}

func verifierGetEntitiesFromLedger(pool *pool.Pool, did string, identifiers json.Slice, actor string) (schemasJSON, credDefsJSON, revocRefDefsJSON, revocRegsJSON string, err error) {
	schemas := make(json.Map)
	credDefs := make(json.Map)
	revocRefDefs := make(json.Map)
	revocRegs := make(json.Map)

	for _, identifier := range identifiers {
		item := identifier.AsMap()

		fmt.Printf(`"%s" -> Get Schema from Ledger`+"\n", actor)
		receivedSchemaID, receivedSchemaJSON, err := getSchema(pool, did, item["schema_id"].String())
		if err != nil {
			return "", "", "", "", err
		}

		schemas[receivedSchemaID] = json.New(receivedSchemaJSON)

		fmt.Printf(`"%s" -> Get Credential Definition from Ledger`+"\n", actor)
		receivedCredDefID, receivedCredDefJSON, err := getCredDef(pool, did, item["cred_def_id"].String())
		if err != nil {
			return "", "", "", "", err
		}

		credDefs[receivedCredDefID] = json.New(receivedCredDefJSON)

		if _, ok := item["rev_reg_seq_no"]; ok {
			//TODO: Get Revocation Definitions and Revocation Registries
		}
	}

	schemasJSON = schemas.JSON()
	credDefsJSON = credDefs.JSON()
	revocRefDefsJSON = revocRefDefs.JSON()
	revocRegsJSON = revocRegs.JSON()

	return
}
