package txs

import (
	"math/big"
	"testing"

	"github.com/cosmos/peggy/cmd/ebrelayer/events"
	"github.com/ethereum/go-ethereum/common"
)

const (
	EthereumPrivateKey    = "ETHEREUM_PRIVATE_KEY"
	TestProphecyID        = 20
	TestClaimType         = 1
	TestValidatorAddress  = "0xc230f38FF05860753840e0d7cbC66128ad308B67"
	TestEthereumReceiver  = "0x7B95B6EC7EbD73572298cEf32Bb54FA408207359"
	TestTokenAddress      = "0x0000000000000000000000000000000000000000"
	TestSymbol            = "eth"
	TestAmount            = 5
	TestExpectedMessage   = "0xfc3c746e966d5f48af553b166b0870b0fa6b6921b353fba67de4e2230392f48b"
	TestExpectedSignature = "0xac349f2452d50d14e11f72de8fc7acde0b47f280a47792470198dcff59358e42425315c0db810dc5d2a7ba5eda7d9cf35cea4f13d550bfa03484df739249c4d401"
	TestAddrHex           = "970e8128ab834e8eac17ab8e3812f010678cf791"
	TestPrivHex           = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232032"
)

// CreateTestProphecyClaimEvent : creates a sample ProphecyClaimEvent for testing purposes
func CreateTestProphecyClaimEvent(t *testing.T) events.NewProphecyClaimEvent {
	testProphecyID := big.NewInt(int64(TestProphecyID))
	testEthereumReceiver := common.HexToAddress(TestEthereumReceiver)
	testValidatorAddress := common.HexToAddress(TestValidatorAddress)
	testTokenAddress := common.HexToAddress(TestTokenAddress)
	testAmount := big.NewInt(int64(TestAmount))

	// Create new ProphecyClaimEvent
	prophecyClaimEvent := events.NewProphecyClaimEvent{
		ProphecyID:       testProphecyID,
		ClaimType:        TestClaimType,
		EthereumReceiver: testEthereumReceiver,
		ValidatorAddress: testValidatorAddress,
		TokenAddress:     testTokenAddress,
		Symbol:           TestSymbol,
		Amount:           testAmount,
	}
	return prophecyClaimEvent
}
