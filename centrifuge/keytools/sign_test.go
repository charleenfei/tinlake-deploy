// +build unit

package keytools

import (
	"fmt"
	"os"
	"testing"

	"github.com/CentrifugeInc/go-centrifuge/centrifuge/keytools/secp256k1"
	"github.com/CentrifugeInc/go-centrifuge/centrifuge/utils"
	"github.com/stretchr/testify/assert"
)

func TestSignMessage(t *testing.T) {

	publicKeyFile := "publicKey"
	privateKeyFile := "privateKey"
	testMsg := []byte("test")

	GenerateSigningKeyPair(publicKeyFile, privateKeyFile, CurveSecp256K1)
	privateKey, err := utils.ReadKeyFromPemFile(privateKeyFile, PrivateKey)
	assert.Nil(t, err)
	publicKey, err := utils.ReadKeyFromPemFile(publicKeyFile, PublicKey)
	assert.Nil(t, err)
	signature, err := SignMessage(privateKey, testMsg, CurveSecp256K1, false)
	assert.Nil(t, err)
	correct := VerifyMessage(publicKey, testMsg, signature, CurveSecp256K1, false)

	os.Remove(publicKeyFile)
	os.Remove(privateKeyFile)

	assert.True(t, correct, "signature or verification didn't work correctly")

}

func TestSignAndVerifyMessageEthereum(t *testing.T) {

	publicKeyFile := "publicKey"
	privateKeyFile := "privateKey"
	testMsg := []byte("Centrifuge likes Ethereum")

	GenerateSigningKeyPair(publicKeyFile, privateKeyFile, CurveSecp256K1)
	privateKey, err := utils.ReadKeyFromPemFile(privateKeyFile, PrivateKey)
	assert.Nil(t, err)
	signature, err := SignMessage(privateKey, testMsg, CurveSecp256K1, true)
	assert.Nil(t, err)

	publicKey, _ := utils.ReadKeyFromPemFile(publicKeyFile, PublicKey)
	address := secp256k1.GetAddress(publicKey)

	fmt.Println("privateKey: ", utils.ByteArrayToHex(privateKey))
	fmt.Println("publicKey: ", utils.ByteArrayToHex(publicKey))
	fmt.Println("address:", address)
	fmt.Println("msg:", string(testMsg[:]))
	fmt.Println("msg in hex:", utils.ByteArrayToHex(testMsg))
	fmt.Println("hash of msg: ", utils.ByteArrayToHex(secp256k1.SignHash(testMsg)))
	fmt.Println("signature:", utils.ByteArrayToHex(signature))
	fmt.Println("Generated Signature can also be verified at https://etherscan.io/verifySig")

	correct := VerifyMessage(publicKey, testMsg, signature, CurveSecp256K1, true)

	os.Remove(publicKeyFile)
	os.Remove(privateKeyFile)

	assert.True(t, correct, "signature or verification didn't work correctly")

}