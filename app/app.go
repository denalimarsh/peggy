package app

import (
	"io"
	"os"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmos "github.com/tendermint/tendermint/libs/os"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	codecstd "github.com/cosmos/cosmos-sdk/codec/std"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/std"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/cosmos/peggy/x/ethbridge"
	"github.com/cosmos/peggy/x/oracle"
)

const (
	appName = "EthereumBridge"
)

var (
	// DefaultCLIHome default home directories for ebcli
	DefaultCLIHome = os.ExpandEnv("$HOME/.ebcli")

	// DefaultNodeHome sets the folder where the application data and configuration will be stored
	DefaultNodeHome = os.ExpandEnv("$HOME/.ebd")

	// ModuleBasics the module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		genutil.AppModuleBasic{},
		bank.AppModuleBasic{},
		staking.AppModuleBasic{},
		distr.AppModuleBasic{},
		params.AppModuleBasic{},
		bank.AppModuleBasic{},
		oracle.AppModuleBasic{},
		ethbridge.AppModuleBasic{},
	)

	// module account permissions
	maccPerms = map[string][]string{
		auth.FeeCollectorName:     nil,
		distr.ModuleName:          nil,
		staking.BondedPoolName:    {auth.Burner, auth.Staking},
		staking.NotBondedPoolName: {auth.Burner, auth.Staking},
		ethbridge.ModuleName:      {auth.Burner, auth.Minter},
	}

	// module accounts that are allowed to receive tokens
	allowedReceivingModAcc = map[string]bool{
		distr.ModuleName: true,
	}
)

// TODO:
// // MakeCodec generates the necessary codecs for Amino
// func MakeCodec() *codec.Codec {
// 	var cdc = codec.New()
// 	ModuleBasics.RegisterCodec(cdc)
// 	vesting.RegisterCodec(cdc)
// 	sdk.RegisterCodec(cdc)
// 	codec.RegisterCrypto(cdc)
// 	return cdc
// }

// EthereumBridgeApp defines the Ethereum-Cosmos peg-zone application
type EthereumBridgeApp struct {
	*baseapp.BaseApp
	cdc *codec.Codec

	// keys to access the substores
	keys  map[string]*sdk.KVStoreKey
	tkeys map[string]*sdk.TransientStoreKey

	// subspaces
	subspaces map[string]params.Subspace

	// keepers
	// TODO: add governance keeper
	accountKeeper   auth.AccountKeeper
	bankKeeper      bank.Keeper
	stakingKeeper   staking.Keeper
	slashingKeeper  slashing.Keeper
	distrKeeper     distr.Keeper
	paramsKeeper    params.Keeper
	ethbridgeKeeper ethbridge.Keeper
	oracleKeeper    oracle.Keeper

	// the module manager
	mm *module.Manager
}

// NewEthereumBridgeApp is a constructor function for EthereumBridgeApp
func NewEthereumBridgeApp(
	logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool,
	baseAppOptions ...func(*baseapp.BaseApp),
) *EthereumBridgeApp {

	// First define the top level codec that will be shared by the different modules
	cdc := codecstd.MakeCodec(ModuleBasics)
	appCodec := codecstd.NewAppCodec(cdc)

	// BaseApp handles interactions with Tendermint through the ABCI protocol
	bApp := baseapp.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetAppVersion(version.Version)

	keys := sdk.NewKVStoreKeys(
		auth.StoreKey, staking.StoreKey, distr.StoreKey,
		oracle.StoreKey, params.StoreKey,
	)
	tkeys := sdk.NewTransientStoreKeys(params.TStoreKey)

	// Here you initialize your application with the store keys it requires
	app := &EthereumBridgeApp{
		BaseApp:   bApp,
		cdc:       cdc,
		keys:      keys,
		tkeys:     tkeys,
		subspaces: make(map[string]params.Subspace),
	}

	// init params keeper and subspaces
	app.paramsKeeper = params.NewKeeper(appCodec, keys[params.StoreKey], tkeys[params.TStoreKey])
	app.subspaces[auth.ModuleName] = app.paramsKeeper.Subspace(auth.DefaultParamspace)
	app.subspaces[bank.ModuleName] = app.paramsKeeper.Subspace(bank.DefaultParamspace)
	app.subspaces[staking.ModuleName] = app.paramsKeeper.Subspace(staking.DefaultParamspace)
	app.subspaces[distr.ModuleName] = app.paramsKeeper.Subspace(distr.DefaultParamspace)

	// set the BaseApp's parameter store
	bApp.SetParamStore(app.paramsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(std.ConsensusParamsKeyTable()))

	// add keepers
	app.accountKeeper = auth.NewAccountKeeper(
		appCodec, keys[auth.StoreKey], app.subspaces[auth.ModuleName], auth.ProtoBaseAccount, maccPerms,
	)
	app.bankKeeper = bank.NewBaseKeeper(
		appCodec, keys[bank.StoreKey], app.accountKeeper, app.subspaces[bank.ModuleName], app.BlacklistedAccAddrs(),
	)
	stakingKeeper := staking.NewKeeper(
		appCodec, keys[staking.StoreKey], app.accountKeeper, app.bankKeeper, app.subspaces[staking.ModuleName],
	)
	app.distrKeeper = distr.NewKeeper(
		appCodec, keys[distr.StoreKey], app.subspaces[distr.ModuleName], app.accountKeeper, app.bankKeeper,
		&stakingKeeper, auth.FeeCollectorName, app.ModuleAccountAddrs(),
	)
	app.slashingKeeper = slashing.NewKeeper(
		appCodec, keys[slashing.StoreKey], &stakingKeeper, app.subspaces[slashing.ModuleName],
	)
	app.oracleKeeper = oracle.NewKeeper(
		cdc, keys[oracle.StoreKey], &stakingKeeper, oracle.DefaultConsensusNeeded,
	)
	app.ethbridgeKeeper = ethbridge.NewKeeper(
		cdc, keys[ethbridge.StoreKey], app.bankKeeper, app.oracleKeeper,
	)

	// Set hooks
	app.stakingKeeper = *stakingKeeper.SetHooks(
		staking.NewMultiStakingHooks(app.distrKeeper.Hooks(), app.slashingKeeper.Hooks()),
	)

	app.mm = module.NewManager(
		genutil.NewAppModule(app.accountKeeper, app.stakingKeeper, app.BaseApp.DeliverTx),
		auth.NewAppModule(appCodec, app.accountKeeper),
		bank.NewAppModule(appCodec, app.bankKeeper, app.accountKeeper),
		slashing.NewAppModule(appCodec, app.slashingKeeper, app.accountKeeper, app.bankKeeper, app.stakingKeeper),
		distr.NewAppModule(appCodec, app.distrKeeper, app.accountKeeper, app.bankKeeper, app.stakingKeeper),
		staking.NewAppModule(appCodec, app.stakingKeeper, app.accountKeeper, app.bankKeeper),
		oracle.NewAppModule(appCodec, app.oracleKeeper),
		ethbridge.NewAppModule(appCodec, app.oracleKeeper, app.bankKeeper, app.accountKeeper, app.ethbridgeKeeper),
		params.NewAppModule(app.paramsKeeper),
	)

	// TODO: mint.ModuleName,
	app.mm.SetOrderBeginBlockers(distr.ModuleName, slashing.ModuleName, staking.ModuleName)
	app.mm.SetOrderEndBlockers(staking.ModuleName)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// TODO:  mint.ModuleName,
	app.mm.SetOrderInitGenesis(
		auth.ModuleName, distr.ModuleName, slashing.ModuleName,
		staking.ModuleName, bank.ModuleName, genutil.ModuleName,
		ethbridge.ModuleName,
	)

	// TODO: add simulator support

	app.mm.RegisterRoutes(app.Router(), app.QueryRouter())

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)

	// initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)

	// TODO: SET ANTE HANDLER INCLUDING IBCKEEPER
	// app.SetAnteHandler(
	// 	ante.NewAnteHandler(
	// 		app.accountKeeper, app.bankKeeper,
	// 		auth.DefaultSigVerificationGasConsumer,
	// 	),
	// )
	app.SetEndBlocker(app.EndBlocker)

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			tmos.Exit(err.Error())
		}
	}

	// TODO: capabilityKeeper
	// ctx := app.BaseApp.NewContext(true, abci.Header{})
	// app.capabilityKeeper.InitializeAndSeal(ctx)

	return app
}

// InitChainer application update at chain initialization
func (app *EthereumBridgeApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState simapp.GenesisState
	app.cdc.MustUnmarshalJSON(req.AppStateBytes, &genesisState)
	return app.mm.InitGenesis(ctx, app.cdc, genesisState)
}

// BeginBlocker application updates every begin block
func (app *EthereumBridgeApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// EndBlocker application updates every end block
func (app *EthereumBridgeApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

// LoadHeight loads a particular height
func (app *EthereumBridgeApp) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *EthereumBridgeApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[auth.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// BlacklistedAccAddrs returns all the app's module account addresses black listed for receiving tokens.
func (app *EthereumBridgeApp) BlacklistedAccAddrs() map[string]bool {
	blacklistedAddrs := make(map[string]bool)
	for acc := range maccPerms {
		blacklistedAddrs[auth.NewModuleAddress(acc).String()] = !allowedReceivingModAcc[acc]
	}

	return blacklistedAddrs
}

// Codec returns simapp's codec
func (app *EthereumBridgeApp) Codec() *codec.Codec {
	return app.cdc
}

// GetKey returns the KVStoreKey for the provided store key
func (app *EthereumBridgeApp) GetKey(storeKey string) *sdk.KVStoreKey {
	return app.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key
func (app *EthereumBridgeApp) GetTKey(storeKey string) *sdk.TransientStoreKey {
	return app.tkeys[storeKey]
}

// GetMaccPerms returns a copy of the module account permissions
func GetMaccPerms() map[string][]string {
	dupMaccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		dupMaccPerms[k] = v
	}
	return dupMaccPerms
}
