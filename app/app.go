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
	"github.com/cosmos/cosmos-sdk/x/auth/ante" // TODO: is this required
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/capability"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/ibc"
	ibcclient "github.com/cosmos/cosmos-sdk/x/ibc/02-client"
	port "github.com/cosmos/cosmos-sdk/x/ibc/05-port"
	transfer "github.com/cosmos/cosmos-sdk/x/ibc/20-transfer"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"

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
		capability.AppModuleBasic{},
		staking.AppModuleBasic{},
		mint.AppModuleBasic{},
		distr.AppModuleBasic{},
		gov.NewAppModuleBasic(
			paramsclient.ProposalHandler, distr.ProposalHandler, upgradeclient.ProposalHandler,
		),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		ibc.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		oracle.AppModuleBasic{},
		ethbridge.AppModuleBasic{},
		evidence.AppModuleBasic{},
		transfer.AppModuleBasic{},
	)

	// module account permissions
	maccPerms = map[string][]string{
		auth.FeeCollectorName:           nil,
		distr.ModuleName:                nil,
		mint.ModuleName:                 {auth.Minter},
		staking.BondedPoolName:          {auth.Burner, auth.Staking},
		staking.NotBondedPoolName:       {auth.Burner, auth.Staking},
		ethbridge.ModuleName:            {auth.Burner, auth.Minter},
		gov.ModuleName:                  {auth.Burner},
		transfer.GetModuleAccountName(): {auth.Minter, auth.Burner},
	}

	// module accounts that are allowed to receive tokens
	allowedReceivingModAcc = map[string]bool{
		distr.ModuleName: true,
	}
)

// EthereumBridgeApp defines the Ethereum-Cosmos peg-zone application
type EthereumBridgeApp struct {
	*baseapp.BaseApp
	cdc *codec.Codec

	invCheckPeriod uint

	// keys to access the substores
	keys    map[string]*sdk.KVStoreKey
	tkeys   map[string]*sdk.TransientStoreKey
	memKeys map[string]*sdk.MemoryStoreKey

	// subspaces
	subspaces map[string]params.Subspace

	// keepers
	accountKeeper    auth.AccountKeeper
	bankKeeper       bank.Keeper
	capabilityKeeper *capability.Keeper
	stakingKeeper    staking.Keeper
	slashingKeeper   slashing.Keeper
	mintKeeper       mint.Keeper
	distrKeeper      distr.Keeper
	govKeeper        gov.Keeper
	crisisKeeper     crisis.Keeper
	upgradeKeeper    upgrade.Keeper
	paramsKeeper     params.Keeper
	ibcKeeper        *ibc.Keeper
	evidenceKeeper   evidence.Keeper
	transferKeeper   transfer.Keeper
	ethbridgeKeeper  ethbridge.Keeper
	oracleKeeper     oracle.Keeper

	// make scoped keepers public for test purposes
	scopedIBCKeeper      capability.ScopedKeeper
	scopedTransferKeeper capability.ScopedKeeper

	// the module manager
	mm *module.Manager

	// TODO: SimulationManager
}

// NewEthereumBridgeApp is a constructor function for EthereumBridgeApp
func NewEthereumBridgeApp(
	logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool,
	invCheckPeriod uint, skipUpgradeHeights map[int64]bool, home string,
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
		auth.StoreKey, bank.StoreKey, staking.StoreKey,
		mint.StoreKey, distr.StoreKey, slashing.StoreKey, gov.StoreKey,
		oracle.StoreKey, params.StoreKey, ibc.StoreKey, upgrade.StoreKey,
		evidence.StoreKey, transfer.StoreKey, capability.StoreKey,
	)

	tkeys := sdk.NewTransientStoreKeys(params.TStoreKey)
	memKeys := sdk.NewMemoryStoreKeys(capability.MemStoreKey)

	// Here you initialize your application with the store keys it requires
	app := &EthereumBridgeApp{
		BaseApp:        bApp,
		cdc:            cdc,
		invCheckPeriod: invCheckPeriod,
		keys:           keys,
		tkeys:          tkeys,
		memKeys:        memKeys,
		subspaces:      make(map[string]params.Subspace),
	}

	// init params keeper and subspaces
	app.paramsKeeper = params.NewKeeper(appCodec, keys[params.StoreKey], tkeys[params.TStoreKey])
	app.subspaces[auth.ModuleName] = app.paramsKeeper.Subspace(auth.DefaultParamspace)
	app.subspaces[bank.ModuleName] = app.paramsKeeper.Subspace(bank.DefaultParamspace)
	app.subspaces[staking.ModuleName] = app.paramsKeeper.Subspace(staking.DefaultParamspace)
	app.subspaces[mint.ModuleName] = app.paramsKeeper.Subspace(mint.DefaultParamspace)
	app.subspaces[distr.ModuleName] = app.paramsKeeper.Subspace(distr.DefaultParamspace)
	app.subspaces[slashing.ModuleName] = app.paramsKeeper.Subspace(slashing.DefaultParamspace)
	app.subspaces[gov.ModuleName] = app.paramsKeeper.Subspace(gov.DefaultParamspace).WithKeyTable(gov.ParamKeyTable())
	app.subspaces[crisis.ModuleName] = app.paramsKeeper.Subspace(crisis.DefaultParamspace)

	// set the BaseApp's parameter store
	bApp.SetParamStore(app.paramsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(std.ConsensusParamsKeyTable()))

	// add capability keeper and ScopeToModule for ibc module
	app.capabilityKeeper = capability.NewKeeper(appCodec, keys[capability.StoreKey], memKeys[capability.MemStoreKey])
	scopedIBCKeeper := app.capabilityKeeper.ScopeToModule(ibc.ModuleName)
	scopedTransferKeeper := app.capabilityKeeper.ScopeToModule(transfer.ModuleName)

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
	app.mintKeeper = mint.NewKeeper(
		appCodec, keys[mint.StoreKey], app.subspaces[mint.ModuleName], &stakingKeeper,
		app.accountKeeper, app.bankKeeper, auth.FeeCollectorName,
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
	app.crisisKeeper = crisis.NewKeeper(
		app.subspaces[crisis.ModuleName], invCheckPeriod, app.bankKeeper, auth.FeeCollectorName,
	)
	app.upgradeKeeper = upgrade.NewKeeper(skipUpgradeHeights, keys[upgrade.StoreKey], appCodec, home)

	// register the proposal types
	govRouter := gov.NewRouter()
	govRouter.AddRoute(gov.RouterKey, gov.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(app.paramsKeeper)).
		AddRoute(distr.RouterKey, distr.NewCommunityPoolSpendProposalHandler(app.distrKeeper)).
		AddRoute(upgrade.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(app.upgradeKeeper))
	app.govKeeper = gov.NewKeeper(
		appCodec, keys[gov.StoreKey], app.subspaces[gov.ModuleName], app.accountKeeper, app.bankKeeper,
		&stakingKeeper, govRouter,
	)

	// Set hooks
	app.stakingKeeper = *stakingKeeper.SetHooks(
		staking.NewMultiStakingHooks(app.distrKeeper.Hooks(), app.slashingKeeper.Hooks()),
	)

	// Create IBC Keeper
	app.ibcKeeper = ibc.NewKeeper(
		app.cdc, keys[ibc.StoreKey], app.stakingKeeper, scopedIBCKeeper,
	)

	// Create Transfer Keepers
	app.transferKeeper = transfer.NewKeeper(
		app.cdc, keys[transfer.StoreKey],
		app.ibcKeeper.ChannelKeeper, &app.ibcKeeper.PortKeeper,
		app.accountKeeper, app.bankKeeper,
		scopedTransferKeeper,
	)
	transferModule := transfer.NewAppModule(app.transferKeeper)

	// Create static IBC router, add transfer route, then set and seal it
	ibcRouter := port.NewRouter()
	ibcRouter.AddRoute(transfer.ModuleName, transferModule)
	app.ibcKeeper.SetRouter(ibcRouter)

	// create evidence keeper with router
	evidenceKeeper := evidence.NewKeeper(
		appCodec, keys[evidence.StoreKey], &app.stakingKeeper, app.slashingKeeper,
	)
	evidenceRouter := evidence.NewRouter().
		AddRoute(ibcclient.RouterKey, ibcclient.HandlerClientMisbehaviour(app.ibcKeeper.ClientKeeper))

	evidenceKeeper.SetRouter(evidenceRouter)
	app.evidenceKeeper = *evidenceKeeper

	app.mm = module.NewManager(
		genutil.NewAppModule(app.accountKeeper, app.stakingKeeper, app.BaseApp.DeliverTx),
		auth.NewAppModule(appCodec, app.accountKeeper),
		bank.NewAppModule(appCodec, app.bankKeeper, app.accountKeeper),
		capability.NewAppModule(*app.capabilityKeeper),
		crisis.NewAppModule(&app.crisisKeeper),
		gov.NewAppModule(appCodec, app.govKeeper, app.accountKeeper, app.bankKeeper),
		mint.NewAppModule(appCodec, app.mintKeeper, app.accountKeeper),
		slashing.NewAppModule(appCodec, app.slashingKeeper, app.accountKeeper, app.bankKeeper, app.stakingKeeper),
		distr.NewAppModule(appCodec, app.distrKeeper, app.accountKeeper, app.bankKeeper, app.stakingKeeper),
		staking.NewAppModule(appCodec, app.stakingKeeper, app.accountKeeper, app.bankKeeper),
		upgrade.NewAppModule(app.upgradeKeeper),
		evidence.NewAppModule(appCodec, app.evidenceKeeper),
		oracle.NewAppModule(appCodec, app.oracleKeeper),
		ethbridge.NewAppModule(appCodec, app.oracleKeeper, app.bankKeeper, app.accountKeeper, app.ethbridgeKeeper),
		ibc.NewAppModule(app.ibcKeeper),
		params.NewAppModule(app.paramsKeeper),
		transferModule,
	)

	app.mm.SetOrderBeginBlockers(
		upgrade.ModuleName, mint.ModuleName, distr.ModuleName, slashing.ModuleName,
		evidence.ModuleName, staking.ModuleName, ibc.ModuleName,
	)
	app.mm.SetOrderEndBlockers(crisis.ModuleName, gov.ModuleName, staking.ModuleName)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	app.mm.SetOrderInitGenesis(
		auth.ModuleName, distr.ModuleName, staking.ModuleName, bank.ModuleName,
		slashing.ModuleName, genutil.ModuleName, ethbridge.ModuleName,
	)

	app.mm.SetOrderInitGenesis(
		capability.ModuleName, auth.ModuleName, distr.ModuleName, staking.ModuleName, bank.ModuleName,
		slashing.ModuleName, gov.ModuleName, mint.ModuleName, crisis.ModuleName, ibc.ModuleName,
		genutil.ModuleName, ethbridge.ModuleName, evidence.ModuleName, transfer.ModuleName,
	)

	app.mm.RegisterInvariants(&app.crisisKeeper)
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter())

	// TODO: app.sm = module.NewSimulationManager(...)
	// TODO: app.sm.RegisterStoreDecoders()

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)
	app.MountMemoryStores(memKeys)

	// initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetAnteHandler(
		ante.NewAnteHandler(
			app.accountKeeper, app.bankKeeper, *app.ibcKeeper,
			ante.DefaultSigVerificationGasConsumer,
		),
	)
	app.SetEndBlocker(app.EndBlocker)

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			tmos.Exit(err.Error())
		}
	}

	ctx := app.BaseApp.NewContext(true, abci.Header{})
	app.capabilityKeeper.InitializeAndSeal(ctx)

	app.scopedIBCKeeper = scopedIBCKeeper
	app.scopedTransferKeeper = scopedTransferKeeper

	return app
}

// Name returns the name of the App
func (app *EthereumBridgeApp) Name() string { return app.BaseApp.Name() }

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

// Codec returns EthereumBridgeApps's codec
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

// TODO: func (app *EthereumBridgeApp) SimulationManager()
