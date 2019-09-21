package client

import (
	"github.com/cosmos/cosmos-sdk/client"
	ethbridgecmd "github.com/cosmos/peggy/x/ethbridge/client/cli"
	"github.com/spf13/cobra"
	amino "github.com/tendermint/go-amino"
)

// ModuleClient exports all client functionality from this module
type ModuleClient struct {
	queryRoute string
	cdc        *amino.Codec
}

func NewModuleClient(queryRoute string, cdc *amino.Codec) ModuleClient {
	return ModuleClient{queryRoute, cdc}
}

// GetQueryCmd returns the cli query commands for this module
func (mc ModuleClient) GetQueryCmd() *cobra.Command {
	// Group ethbridge queries under a subcommand
	ethBBridgeQueryCmd := &cobra.Command{
		Use:   "ethbridge",
		Short: "Querying commands for the ethbridge module",
	}

	ethBBridgeQueryCmd.AddCommand(client.GetCommands(
		ethbridgecmd.GetCmdGetEthBridgeProphecy(mc.queryRoute, mc.cdc),
	)...)

	return ethBBridgeQueryCmd
}

// GetTxCmd returns the transaction commands for this module
func (mc ModuleClient) GetTxCmd() *cobra.Command {
	ethBridgeTxCmd := &cobra.Command{
		Use:   "ethbridge",
		Short: "EthBridge transactions subcommands",
	}

	ethBridgeTxCmd.AddCommand(client.PostCommands(
		ethbridgecmd.GetCmdCreateEthBridgeClaim(mc.cdc),
	)...)

	return ethBridgeTxCmd
}
