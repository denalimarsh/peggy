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
  "strings"
  "strconv"
  "fmt"

  sdk "github.com/cosmos/cosmos-sdk/types"
  "github.com/swishlabsco/cosmos-ethereum-bridge/cmd/ebrelayer/events"
  "github.com/swishlabsco/cosmos-ethereum-bridge/x/ethbridge/types"
)

func ParseEvent(event *events.LockEvent) (int, string, sdk.AccAddress, sdk.Coins, error) {
  //validator sdk.AccAddress, 

  // witnessClaim := types.EthBridgeClaim{}

  // Nonce type casting (*big.Int -> int)
  nonce, nonceErr := strconv.Atoi(event.Nonce.String())
  if nonceErr != nil {
    fmt.Errorf("%s", nonceErr)
  }
  // witnessClaim.Nonce = nonce

  // EthereumSender type casting (address.common -> string)
  sender := event.From.Hex()

  // CosmosReceiver type casting (bytes[] -> sdk.AccAddress)
  recipient, recipientErr := sdk.AccAddressFromBech32(string(event.To[:]))
  if recipientErr != nil {
    fmt.Errorf("%s", recipientErr)
  }
  // witnessClaim.CosmosReceiver = recipient

  // Validator is already the correct type (sdk.AccAddress)
  // witnessClaim.Validator = validator

  // Amount type casting (*big.Int -> sdk.Coins)
  ethereumCoin := []string {event.Value.String(),"ethereum"}
  weiAmount, coinErr := sdk.ParseCoins(strings.Join(ethereumCoin, ""))
  if coinErr != nil {
    fmt.Errorf("%s", coinErr)
  }
  // witnessClaim.Amount = weiAmount

  // return witnessClaim, nil
  return nonce, sender, recipient, weiAmount, nil
}

func PackageClaim(nonce int, sender string, recipient sdk.AccAddress, amount sdk.Coins, validator sdk.AccAddress,) types.EthBridgeClaim {
  claim := types.EthBridgeClaim {
    Nonce = nonce,
    EthereumSender = sender,
    CosmosReceiver = recipient,
    Validator = validator,
    Amount = amount,
  }
  return claim
}
