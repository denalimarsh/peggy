package txs

import (
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	// TODO: use these in TestGenerateClaimMessage
	// TestProphecyID       = 20
	// TestCosmosSender     = "cosmos1qwnw2r9ak79536c4dqtrtk2pl2nlzpqh763rls"
	// TestEthereumReceiver = "0x7B95B6EC7EbD73572298cEf32Bb54FA408207359"
	// TestTokenAddress     = "0x0000000000000000000000000000000000000000"
	// TestAmount           = 5
	// TestValidator        = "0xc230f38FF05860753840e0d7cbC66128ad308B67"
	TestAddrHex = "970e8128ab834e8eac17ab8e3812f010678cf791"
	TestPrivHex = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232032"
)

// func TestGenerateClaimMessage(t *testing.T) string {
// }

func TestSignClaim(t *testing.T) {
	// Set and get env variables to replicate relayer
	os.Setenv("ETHEREUM_PRIVATE_KEY", TestPrivHex)
	rawKey := os.Getenv("ETHEREUM_PRIVATE_KEY")

	// Load signer's private key and address
	key, _ := crypto.HexToECDSA(rawKey)
	signerAddr := common.HexToAddress(TestAddrHex)

	// TODO: use GenerateClaimMessage(...)
	msg := crypto.Keccak256Hash([]byte("foo"))

	// Sign the message
	sig, err := SignClaim(msg.Bytes(), key)
	if err != nil {
		t.Errorf("sign error: %s", err)
	}

	// Recover signer's public key and address
	recoveredPub, err := crypto.Ecrecover(msg.Bytes(), sig)
	if err != nil {
		t.Errorf("ecrecover error: %s", err)
	}
	pubKey, _ := crypto.UnmarshalPubkey(recoveredPub)
	recoveredAddr := crypto.PubkeyToAddress(*pubKey)

	// Confirm that the recovered address is correct
	if signerAddr != recoveredAddr {
		t.Errorf("address mismatch: want: %x have: %x", signerAddr, recoveredAddr)
	}
}

// func TestInvalidSign(t *testing.T) {
// 	if _, err := crypto.Sign(make([]byte, 1), nil); err == nil {
// 		t.Errorf("expected sign with hash 1 byte to error")
// 	}
// 	if _, err := crypto.Sign(make([]byte, 33), nil); err == nil {
// 		t.Errorf("expected sign with hash 33 byte to error")
// 	}
// }

// func TestLoadPrivateKey(t *testing.T) {
// 	key, err := LoadPrivateKey()
// 	require.NoError(t, err)

// 	_ = key

// 	// require.True(t, strings.Contains(key, "unrecognized ethbridge message type: "))
// }

// func TestLoadSender(t *testing.T) {

// }
