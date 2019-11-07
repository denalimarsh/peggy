package txs

import (
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	solsha3 "github.com/miguelmota/go-solidity-sha3"
	"github.com/stretchr/testify/require"
)

const (
	TestCoinString = "9921ruby"
	TestCoinSymbol = "ruby"
	TestCoinAmount = 9921
)

func TestGenerateClaimMessage(t *testing.T) {
	// Create new test ProphecyClaimEvent
	prophecyClaimEvent := CreateTestProphecyClaimEvent(t)
	// Generate claim message from ProphecyClaim
	message := GenerateClaimMessage(prophecyClaimEvent)

	// Confirm that the generated message matches the expected generated message
	require.Equal(t, message.Hex(), TestExpectedMessage)
}

func TestPrepareMessageForSigning(t *testing.T) {
	// Create new test ProphecyClaimEvent
	prophecyClaimEvent := CreateTestProphecyClaimEvent(t)
	// Generate claim message from ProphecyClaim
	message := GenerateClaimMessage(prophecyClaimEvent)

	// Simulate message hashing, prefixing
	hashedMsg := solsha3.SoliditySHA3(solsha3.String(message.Hex()))
	prefixedMessage := solsha3.SoliditySHA3(solsha3.String("\x19Ethereum Signed Message:\n32"), solsha3.Bytes32(hashedMsg))

	// Prepare the message for signing
	preparedMessage := PrepareMsgForSigning(message.Hex())

	// Confirm that the prefixed message matches the prepared message
	require.Equal(t, preparedMessage, prefixedMessage)
}

func TestSignClaim(t *testing.T) {
	// Set and get env variables to replicate relayer
	os.Setenv(EthereumPrivateKey, TestPrivHex)
	rawKey := os.Getenv(EthereumPrivateKey)

	// Load signer's private key and address
	key, _ := crypto.HexToECDSA(rawKey)
	signerAddr := common.HexToAddress(TestAddrHex)

	// Create new test ProphecyClaimEvent
	prophecyClaimEvent := CreateTestProphecyClaimEvent(t)

	// Generate claim message from ProphecyClaim
	message := GenerateClaimMessage(prophecyClaimEvent)

	// Prepare the message (required for signature verification on contract)
	prefixedHashedMsg := PrepareMsgForSigning(message.Hex())

	// Sign the message using the validator's private key
	signature, err := SignClaim(prefixedHashedMsg, key)
	require.NoError(t, err)

	// Recover signer's public key and address
	recoveredPub, err := crypto.Ecrecover(prefixedHashedMsg, signature)
	require.NoError(t, err)
	pubKey, _ := crypto.UnmarshalPubkey(recoveredPub)
	recoveredAddr := crypto.PubkeyToAddress(*pubKey)

	// Confirm that the recovered address is correct
	require.Equal(t, recoveredAddr, signerAddr)
}

func TestGetSymbolFromCoin(t *testing.T) {
	testCoinAmount := big.NewInt(int64(TestCoinAmount))

	// Parse the coin into (symbol, amount)
	symbol, amount := GetSymbolAmountFromCoin(TestCoinString)

	// Confirm that the symbol  is correct
	require.Equal(t, symbol, TestCoinSymbol)

	// Confirm that the amount  is correct
	require.Equal(t, amount, testCoinAmount)
}

// TODO: "Error loading .env file"
// func TestLoadPrivateKey(t *testing.T) {
// 	// Set env variable "ETHEREUM_PRIVATE_KEY"
// 	os.Setenv(EthereumPrivateKey, TestPrivHex)

// 	// Load the validators private key from config
// 	key, err := LoadPrivateKey()
// 	require.NoError(t, err)

// 	// Get env variable "ETHEREUM_PRIVATE_KEY"
// 	testKey := os.Getenv(EthereumPrivateKey)

// 	// Confirm that the key matches the testKey
// 	require.Equal(t, key, testKey)
// }
