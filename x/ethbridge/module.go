package ethbridge

import (
	"encoding/json"

	"github.com/cosmos/peggy/x/ethbridge/client"
	"github.com/cosmos/peggy/x/ethbridge/types"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic defines the basic application module used by the ethbridge module.
type AppModuleBasic struct {
	cdc codec.Marshaler
}

var _ module.AppModuleBasic = AppModuleBasic{}

// Name returns the ethbridge module's name.
func (AppModuleBasic) Name() string {
	return ModuleName
}

// RegisterCodec registers the ethbridge module's types for the given codec.
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	types.RegisterCodec(cdc)
}

// DefaultGenesis returns default genesis state as raw bytes for the ethbridge
// module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONMarshaler) json.RawMessage {
	return nil
}

// ValidateGenesis performs genesis state validation for the ethbridge module.
func (AppModuleBasic) ValidateGenesis(_ codec.JSONMarshaler, _ json.RawMessage) error {
	return nil
}

// RegisterRESTRoutes registers the REST routes for the ethbridge module.
func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	client.RegisterRESTRoutes(ctx, rtr, StoreKey)
}

// GetTxCmd returns the root tx command for the ethbridge module.
func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return client.GetTxCmd(StoreKey, cdc)
}

// GetQueryCmd returns no root query command for the ethbridge module.
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return client.GetQueryCmd(StoreKey, cdc)
}

//____________________________________________________________________________

// AppModuleSimulation defines the module simulation functions used by the ethbridge module.
type AppModuleSimulation struct{}

// AppModule implements an application module for the ethbridge module.
type AppModule struct {
	AppModuleBasic
	AppModuleSimulation

	keeper        Keeper
	BankKeeper    types.BankKeeper
	AccountKeeper types.AccountKeeper
	OracleKeeper  types.OracleKeeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(
	cdc codec.Marshaler, oracleKeeper types.OracleKeeper,
	bankKeeper types.BankKeeper, accountKeeper types.AccountKeeper,
	keeper Keeper) AppModule {

	return AppModule{
		AppModuleBasic:      AppModuleBasic{cdc: cdc},
		AppModuleSimulation: AppModuleSimulation{},

		keeper:        keeper,
		BankKeeper:    bankKeeper,
		AccountKeeper: accountKeeper,
		OracleKeeper:  oracleKeeper,
	}
}

// Name returns the ethbridge module's name.
func (AppModule) Name() string {
	return ModuleName
}

// RegisterInvariants registers the ethbridge module invariants.
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
}

// Route returns the message routing key for the ethbridge module.
func (AppModule) Route() string {
	return RouterKey
}

// NewHandler returns an sdk.Handler for the ethbridge module.
func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.AccountKeeper, am.keeper)
}

// QuerierRoute returns the ethbridge module's querier route name.
func (AppModule) QuerierRoute() string {
	return QuerierRoute
}

// NewQuerierHandler returns the ethbridge module sdk.Querier.
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier(am.OracleKeeper)
}

// InitGenesis performs genesis initialization for the ethbridge module. It returns
// no validator updates.
func (am AppModule) InitGenesis(_ sdk.Context, _ codec.JSONMarshaler, _ json.RawMessage) []abci.ValidatorUpdate {
	// bridgeAccount := supply.NewEmptyModuleAccount(ModuleName, supply.Burner, supply.Minter)
	// am.SupplyKeeper.SetModuleAccount(ctx, bridgeAccount)
	// ---------------------------------------------
	// TODO: replace module account with bank account
	// app.AccountKeeper.SetParams(ctx, auth.DefaultParams())
	// am.BankKeeper.SetSendEnabled(ctx, true)
	return nil
}

// ExportGenesis returns the exported genesis state as raw bytes for the ethbridge
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context, _ codec.JSONMarshaler) json.RawMessage {
	return nil
}

// BeginBlock returns the begin blocker for the ethbridge module.
func (AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}

// EndBlock returns the end blocker for the ethbridge module. It returns no validator
// updates.
func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return nil
}
