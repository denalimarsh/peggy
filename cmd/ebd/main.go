package main

import (
	"encoding/json"
	"io"

	"github.com/spf13/cobra"

	codecstd "github.com/cosmos/cosmos-sdk/codec/std"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/cosmos/cosmos-sdk/x/staking"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/peggy/app"
)

const flagInvCheckPeriod = "inv-check-period"

var invCheckPeriod uint

func main() {
	cdc := codecstd.MakeCodec(app.ModuleBasics)
	appCodec := codecstd.NewAppCodec(cdc)

	// TODO: set custom bech32 prefixes for peggy
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(sdk.Bech32PrefixAccAddr, sdk.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(sdk.Bech32PrefixValAddr, sdk.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(sdk.Bech32PrefixConsAddr, sdk.Bech32PrefixConsPub)
	config.Seal()

	ctx := server.NewDefaultContext()
	cobra.EnableCommandSorting = false

	rootCmd := &cobra.Command{
		Use:               "ebd",
		Short:             "Ethereum Bridge App Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}

	rootCmd.AddCommand(genutilcli.InitCmd(ctx, cdc, app.ModuleBasics, app.DefaultNodeHome))
	rootCmd.AddCommand(genutilcli.CollectGenTxsCmd(ctx, cdc, bank.GenesisBalancesIterator{}, app.DefaultNodeHome))
	rootCmd.AddCommand(
		genutilcli.GenTxCmd(
			ctx, cdc, app.ModuleBasics, staking.AppModuleBasic{},
			bank.GenesisBalancesIterator{}, app.DefaultNodeHome, app.DefaultCLIHome,
		),
	)
	rootCmd.AddCommand(genutilcli.ValidateGenesisCmd(ctx, cdc, app.ModuleBasics))
	rootCmd.AddCommand(AddGenesisAccountCmd(ctx, cdc, appCodec, app.DefaultNodeHome, app.DefaultCLIHome))

	// TODO: testnet cmd
	// rootCmd.AddCommand(flags.NewCompletionCmd(rootCmd, true))
	// rootCmd.AddCommand(testnetCmd(ctx, cdc, app.ModuleBasics, bank.GenesisBalancesIterator{}))
	// rootCmd.AddCommand(replayCmd())
	// rootCmd.AddCommand(debug.Cmd(cdc))

	// TODO:
	server.AddCommands(ctx, cdc, rootCmd, newApp, exportAppStateAndTMValidators)

	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "EB", app.DefaultNodeHome)
	if err := executor.Execute(); err != nil {
		panic(err)
	}
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {
	// var cache sdk.MultiStorePersistentCache

	// if viper.GetBool(server.FlagInterBlockCache) {
	// 	cache = store.NewCommitKVStoreCacheManager()
	// }

	// skipUpgradeHeights := make(map[int64]bool)
	// for _, h := range viper.GetIntSlice(server.FlagUnsafeSkipUpgrades) {
	// 	skipUpgradeHeights[int64(h)] = true
	// }

	// ctx := app.BaseApp.NewContext(true, abci.Header{})

	// return app.NewEthereumBridgeApp(
	// 	logger, db, false, ctx,
	// )
	return app.NewEthereumBridgeApp(logger, db, traceStore, true)
}

func exportAppStateAndTMValidators(
	logger log.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool, jailWhiteList []string,
) (json.RawMessage, []tmtypes.GenesisValidator, *abci.ConsensusParams, error) {

	if height != -1 {
		ebApp := app.NewEthereumBridgeApp(logger, db, traceStore, false)
		if err := ebApp.LoadHeight(height); err != nil {
			return nil, nil, nil, err
		}
		return ebApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
	}

	ebApp := app.NewEthereumBridgeApp(logger, db, traceStore, true)
	return ebApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
}
