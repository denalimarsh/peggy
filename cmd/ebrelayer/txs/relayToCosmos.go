package txs

// ------------------------------------------------------------
//	Relay : Builds and encodes EthBridgeClaim Msgs with the
//  	specified variables, before presenting the unsigned
//      transaction to validators for optional signing.
//      Once signed, the data packets are sent as transactions
//      on the Cosmos Bridge.
// ------------------------------------------------------------

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/cosmos/peggy/x/ethbridge"
	"github.com/cosmos/peggy/x/ethbridge/types"
)

// RelayLockToCosmos : RelayLockToCosmos applies validator's signature to an EthBridgeClaim message
//		containing information about an event on the Ethereum blockchain before relaying to the Bridge
func RelayLockToCosmos(
	chainID string,
	cdc *codec.Codec,
	validatorAddress sdk.ValAddress,
	moniker string,
	passphrase string,
	claim *types.EthBridgeClaim,
	rpcURL string,
) error {
	cliCtx := context.NewCLIContext().
		WithCodec(cdc).
		WithFromAddress(sdk.AccAddress(validatorAddress)).
		WithFromName(moniker)

	if rpcURL != "" {
		cliCtx = cliCtx.WithNodeURI(rpcURL)
	}

	cliCtx.SkipConfirm = true

	txBldr := authtxb.NewTxBuilderFromCLI().
		WithTxEncoder(utils.GetTxEncoder(cdc)).
		WithChainID(chainID)

	accountRetriever := authtypes.NewAccountRetriever(cliCtx)

	err := accountRetriever.EnsureExists((sdk.AccAddress(claim.ValidatorAddress)))
	if err != nil {
		return err
	}

	msg := ethbridge.NewMsgCreateEthBridgeClaim(*claim)

	err = msg.ValidateBasic()
	if err != nil {
		return err
	}

	// Prepare tx
	txBldr, err = utils.PrepareTxBuilder(txBldr, cliCtx)
	if err != nil {
		return err
	}

	// Build and sign the transaction
	txBytes, err := txBldr.BuildAndSign(moniker, passphrase, []sdk.Msg{msg})
	if err != nil {
		return err
	}

	// Broadcast to a Tendermint node
	res, err := cliCtx.BroadcastTxSync(txBytes)
	if err != nil {
		return err
	}

	cliCtx.PrintOutput(res)
	return nil
}
