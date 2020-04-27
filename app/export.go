package app

import (
	"encoding/json"

	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/cosmos/cosmos-sdk/x/staking"
)

// ExportAppStateAndValidators export the state of the eth peg-zone for a genesis file
func (app *EthereumBridgeApp) ExportAppStateAndValidators(forZeroHeight bool, jailWhiteList []string,
) (appState json.RawMessage, validators []tmtypes.GenesisValidator, cp *abci.ConsensusParams, err error) {

	// as if they could withdraw from the start of the next block
	ctx := app.NewContext(true, abci.Header{Height: app.LastBlockHeight()})

	// if forZeroHeight {
	// 	app.prepForZeroHeightGenesis(ctx, jailWhiteList)
	// }

	genState := app.mm.ExportGenesis(ctx, app.cdc)
	appState, err = codec.MarshalJSONIndent(app.cdc, genState)
	if err != nil {
		return nil, nil, nil, err
	}

	validators = staking.WriteValidators(ctx, app.stakingKeeper)
	return appState, validators, app.BaseApp.GetConsensusParams(ctx), nil
}
