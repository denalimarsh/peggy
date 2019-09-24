package querier

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/peggy/x/ethbridge/types"
	keep "github.com/cosmos/peggy/x/oracle/keeper"
	oracletypes "github.com/cosmos/peggy/x/oracle/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the oracle Querier
const (
	QueryEthProphecy = "prophecies"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper keep.Keeper, cdc *codec.Codec, codespace sdk.CodespaceType) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryEthProphecy:
			return queryEthProphecy(ctx, cdc, req, keeper, codespace)
		default:
			return nil, sdk.ErrUnknownRequest("unknown ethbridge query endpoint")
		}
	}
}

func queryEthProphecy(ctx sdk.Context, cdc *codec.Codec, req abci.RequestQuery, keeper keep.Keeper, codespace sdk.CodespaceType) (res []byte, errSdk sdk.Error) {
	var params types.QueryEthProphecyParams

	err := cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return []byte{}, sdk.ErrInternal(sdk.AppendMsgToErr("failed to parse params: %s", err.Error()))
	}
	id := strconv.Itoa(params.ChainID) + strconv.Itoa(params.Nonce) + params.EthereumSender.String()
	prophecy, errSdk := keeper.GetProphecy(ctx, id)
	if errSdk != nil {
		return []byte{}, oracletypes.ErrProphecyNotFound(codespace)
	}
	// TODO: Pass these values as real values
	bridgeAdress := types.NewEthereumAddress("0x0000000000000000000000000000000000000000")
	tokenContract := types.NewEthereumAddress("0x0000000000000000000000000000000000000000")
	bridgeClaims, errSdk := types.MapOracleClaimsToEthBridgeClaims(3, bridgeAdress, params.Nonce, "eth", tokenContract, params.EthereumSender, prophecy.ValidatorClaims, types.CreateEthClaimFromOracleString)
	if errSdk != nil {
		return []byte{}, errSdk
	}

	response := types.NewQueryEthProphecyResponse(prophecy.ID, prophecy.Status, bridgeClaims)

	bz, err := cdc.MarshalJSONIndent(response, "", "  ")
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return bz, nil
}
