package txs

// --------------------------------------------------------
//      Parser
//
//      Parses structs containing event information into
//      unsigned transactions for validators to sign, then
//      relays the data packets as transactions on the
//      Cosmos Bridge.
// --------------------------------------------------------

import (
	"errors"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/peggy/cmd/ebrelayer/events"
	ethbridgeTypes "github.com/cosmos/peggy/x/ethbridge/types"
	"github.com/ethereum/go-ethereum/common"
)

// ParsePayload : parses and packages a LockEvent struct with a validator address in an EthBridgeClaim msg
func ParsePayload(valAddr sdk.ValAddress, event *events.LockEvent) (ethbridgeTypes.EthBridgeClaim, error) {

	witnessClaim := ethbridgeTypes.EthBridgeClaim{}

	bridgeContractAddress := ethbridgeTypes.NewEthereumAddress(event.BridgeContractAddress.Hex())

	// Nonce type casting (*big.Int -> int)
	nonce := int(event.Nonce.Int64())

	// Sender type casting (address.common -> string)
	sender := ethbridgeTypes.NewEthereumAddress(event.From.Hex())

	// Recipient type casting ([]bytes -> sdk.AccAddress)
	recipient, err := sdk.AccAddressFromBech32(string(event.To[:]))
	if err != nil {
		return witnessClaim, err
	}
	if recipient.Empty() {
		return witnessClaim, errors.New("empty recipient address")
	}

	// Sender type casting (address.common -> string)
	tokenContractAddress := ethbridgeTypes.NewEthereumAddress(event.TokenContractAddress.Hex())

	// Symbol formatted to lowercase
	symbol := strings.ToLower(event.Symbol)
	if symbol == "eth" && event.TokenContractAddress != common.HexToAddress("0x0000000000000000000000000000000000000000") {
		return witnessClaim, errors.New("symbol \"eth\" must have null address set as token address")
	}

	// Amount type casting (*big.Int -> sdk.Coins)
	coins := sdk.Coins{sdk.NewInt64Coin(symbol, event.Value.Int64())}

	// Package the information in a unique EthBridgeClaim
	witnessClaim.ChainID = event.ChainId
	witnessClaim.BridgeContractAddress = bridgeContractAddress
	witnessClaim.Nonce = nonce
	witnessClaim.TokenContractAddress = tokenContractAddress
	witnessClaim.Symbol = symbol
	witnessClaim.EthereumSender = sender
	witnessClaim.ValidatorAddress = valAddr
	witnessClaim.CosmosReceiver = recipient
	witnessClaim.Amount = coins

	return witnessClaim, nil
}

// type LockEvent struct {
// 	ChainID               string
// 	BridgeContractAddress common.Address
// 	Id                    [32]byte
// 	From                  common.Address
// 	To                    []byte
// 	TokenContractAddress  common.Address
// 	Symbol                string
// 	Value                 *big.Int
// 	Nonce                 *big.Int
// }
