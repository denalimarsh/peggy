package relayer

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"
	"sync"
	"time"

	sdkContext "github.com/cosmos/cosmos-sdk/client/context"
	ckeys "github.com/cosmos/cosmos-sdk/client/keys"
	aminocodec "github.com/cosmos/cosmos-sdk/codec"
	codecstd "github.com/cosmos/cosmos-sdk/codec/std"
	keys "github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/types"
	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ctypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	abci "github.com/tendermint/tendermint/abci/types"
	tmLog "github.com/tendermint/tendermint/libs/log"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
	libclient "github.com/tendermint/tendermint/rpc/lib/client"

	"github.com/cosmos/peggy/cmd/ebrelayer/contract"
	"github.com/cosmos/peggy/cmd/ebrelayer/txs"
	"github.com/cosmos/peggy/cmd/ebrelayer/types"
	"github.com/cosmos/peggy/x/ethbridge"
)

// TODO: Move relay functionality out of EthereumSub into a new Relayer parent struct
// TODO: Client                  *ethclient.Client

var (
	gas           = uint64(200000)
	gasPrices     = "0.025stake"
	gasAdjustment = 1.0
)

// EthereumSub is an Ethereum listener that can relay txs to Cosmos and Ethereum
type EthereumSub struct {
	HomePath                string
	Cdc                     *contextualStdCodec
	Amino                   *contextualAminoCodec
	EthProvider             string
	Keybase                 keys.Keyring
	RegistryContractAddress common.Address
	ValidatorName           string
	ValidatorAddress        sdk.ValAddress
	Client                  rpcclient.Client
	PrivateKey              *ecdsa.PrivateKey
	Logger                  tmLog.Logger
	CosmosChainID           string
	CliCtx                  sdkContext.CLIContext
}

// NewEthereumSub initializes a new EthereumSub
func NewEthereumSub(homePath string, rpcURL string, cdc *codecstd.Codec, amino *aminocodec.Codec,
	validatorMoniker, chainID, ethProvider string, registryContractAddress common.Address,
	privateKey *ecdsa.PrivateKey, kr keys.Keyring, logger tmLog.Logger,
) (EthereumSub, error) {

	// Initialize a new HTTP tendermint client
	timeout := time.Duration(time.Hour * 24 * 7) // One week timeout
	rpcAddr := fmt.Sprintf("http://localhost:%s", "26657")
	client, err := newRPCClient(rpcAddr, timeout)
	if err != nil {
		return EthereumSub{}, err
	}

	_, err = sdk.ParseDecCoins(gasPrices)
	if err != nil {
		return EthereumSub{}, err
	}

	// TODO: Don't really need these?
	contextualCdc := newContextualStdCodec(cdc, UseSDKContext)
	contextualAminoCdc := newContextualAminoCodec(amino, UseSDKContext)

	// Load validator details
	validatorAddress, validatorName, err := LoadValidatorCredentials(validatorMoniker, kr)
	if err != nil {
		return EthereumSub{}, err
	}

	// Load CLI context
	cliCtx := sdkContext.NewCLIContext().
		WithCodec(amino).
		WithFromAddress(sdk.AccAddress(validatorAddress)).
		WithFromName(validatorName)

	if rpcURL != "" {
		cliCtx = cliCtx.WithNodeURI(rpcURL)
	}
	cliCtx.SkipConfirm = true

	// Validate sender's account (will be required for tx building)
	accGetter := auth.NewAccountRetriever(cdc, cliCtx)
	if err := accGetter.EnsureExists(sdk.AccAddress(validatorAddress)); err != nil {
		return EthereumSub{}, err
	}

	return EthereumSub{
		HomePath:                homePath,
		Cdc:                     contextualCdc,
		Amino:                   contextualAminoCdc,
		EthProvider:             ethProvider,
		RegistryContractAddress: registryContractAddress,
		ValidatorName:           validatorName,
		Keybase:                 kr,
		ValidatorAddress:        validatorAddress,
		Client:                  client,
		PrivateKey:              privateKey,
		Logger:                  logger,
		CosmosChainID:           chainID,
		CliCtx:                  cliCtx,
	}, nil
}

// LoadValidatorCredentials : loads validator's credentials (address, moniker, and passphrase)
func LoadValidatorCredentials(validatorFrom string, kr keys.Keyring) (sdk.ValAddress, string, error) {
	// Get the validator's name and account address using their moniker
	validatorAccAddress, validatorName, err := sdkContext.GetFromFields(kr, validatorFrom, false)
	if err != nil {
		return sdk.ValAddress{}, "", err
	}
	validatorAddress := sdk.ValAddress(validatorAccAddress)

	// Confirm that the key is valid
	_, err = authtxb.MakeSignature(kr, validatorName, ckeys.DefaultKeyPass, authtxb.StdSignMsg{})
	if err != nil {
		return sdk.ValAddress{}, "", err
	}

	return validatorAddress, validatorName, nil
}

// Start an Ethereum chain subscription
func (sub EthereumSub) Start() {
	client, err := SetupWebsocketEthClient(sub.EthProvider)
	if err != nil {
		sub.Logger.Error(err.Error())
		os.Exit(1)
	}
	sub.Logger.Info("Started Ethereum websocket with provider:", sub.EthProvider)

	clientChainID, err := client.NetworkID(context.Background())
	if err != nil {
		sub.Logger.Error(err.Error())
		os.Exit(1)
	}

	// We will check logs for new events
	logs := make(chan ctypes.Log)

	// Start BridgeBank subscription, prepare contract ABI and LockLog event signature
	bridgeBankAddress, subBridgeBank := sub.startContractEventSub(logs, client, txs.BridgeBank)
	bridgeBankContractABI := contract.LoadABI(txs.BridgeBank)
	eventLogLockSignature := bridgeBankContractABI.Events[types.LogLock.String()].Id().Hex()

	// Start CosmosBridge subscription, prepare contract ABI and LogNewProphecyClaim event signature
	cosmosBridgeAddress, subCosmosBridge := sub.startContractEventSub(logs, client, txs.CosmosBridge)
	cosmosBridgeContractABI := contract.LoadABI(txs.CosmosBridge)
	eventLogNewProphecyClaimSignature := cosmosBridgeContractABI.Events[types.LogNewProphecyClaim.String()].Id().Hex()

	for {
		select {
		// Handle any errors
		case err := <-subBridgeBank.Err():
			sub.Logger.Error(err.Error())
		case err := <-subCosmosBridge.Err():
			sub.Logger.Error(err.Error())
		// vLog is raw event data
		case vLog := <-logs:
			sub.Logger.Info(fmt.Sprintf("Witnessed tx %s on block %d\n", vLog.TxHash.Hex(), vLog.BlockNumber))
			var err error
			switch vLog.Topics[0].Hex() {
			case eventLogLockSignature:
				err = sub.handleLogLock(clientChainID, bridgeBankAddress, bridgeBankContractABI,
					types.LogLock.String(), vLog)
			case eventLogNewProphecyClaimSignature:
				err = sub.handleLogNewProphecyClaim(cosmosBridgeAddress, cosmosBridgeContractABI,
					types.LogNewProphecyClaim.String(), vLog)
			}
			// TODO: Check local events store for status, if retryable, attempt relay again
			if err != nil {
				sub.Logger.Error(err.Error())
			}
		}
	}
}

// startContractEventSub : starts an event subscription on the specified Peggy contract
func (sub EthereumSub) startContractEventSub(logs chan ctypes.Log, client *ethclient.Client,
	contractName txs.ContractRegistry) (common.Address, ethereum.Subscription) {
	// Get the contract address for this subscription
	subContractAddress, err := txs.GetAddressFromBridgeRegistry(client, sub.RegistryContractAddress, contractName)
	if err != nil {
		sub.Logger.Error(err.Error())
	}

	// We need the address in []bytes for the query
	subQuery := ethereum.FilterQuery{
		Addresses: []common.Address{subContractAddress},
	}

	// Start the contract subscription
	contractSub, err := client.SubscribeFilterLogs(context.Background(), subQuery, logs)
	if err != nil {
		sub.Logger.Error(err.Error())
	}
	sub.Logger.Info(fmt.Sprintf("Subscribed to %v contract at address: %s", contractName, subContractAddress.Hex()))
	return subContractAddress, contractSub
}

// handleLogLock unpacks a LogLock event, converts it to a ProphecyClaim, and relays a tx to Cosmos
func (sub EthereumSub) handleLogLock(clientChainID *big.Int, contractAddress common.Address,
	contractABI abi.ABI, eventName string, cLog ctypes.Log) error {
	// Parse the event's attributes via contract ABI
	event := types.LockEvent{}
	err := contractABI.Unpack(&event, eventName, cLog.Data)
	if err != nil {
		sub.Logger.Error("error unpacking: %v", err)
	}
	event.BridgeContractAddress = contractAddress
	event.EthereumChainID = clientChainID
	sub.Logger.Info(event.String())

	// Add the event to the record
	types.NewEventWrite(cLog.TxHash.Hex(), event)

	prophecyClaim, err := txs.LogLockToEthBridgeClaim(sub.ValidatorAddress, &event)
	if err != nil {
		return err
	}

	// Packages the claim as a Tendermint message
	msg := ethbridge.NewMsgCreateEthBridgeClaim(prophecyClaim)
	err = msg.ValidateBasic()
	if err != nil {
		return err
	}

	// Send msg
	txRes, err := sub.SendMsg(msg)
	if err != nil {
		return err
	}

	res, err := sub.Amino.MarshalJSON(txRes)
	if err != nil {
		return err
	}
	fmt.Println(string(res))
	return nil
}

// Unpacks a handleLogNewProphecyClaim event, builds a new OracleClaim, and relays it to Ethereum
func (sub EthereumSub) handleLogNewProphecyClaim(contractAddress common.Address, contractABI abi.ABI,
	eventName string, cLog ctypes.Log) error {
	// Parse the event's attributes via contract ABI
	event := types.ProphecyClaimEvent{}
	err := contractABI.Unpack(&event, eventName, cLog.Data)
	if err != nil {
		sub.Logger.Error("error unpacking: %v", err)
	}
	sub.Logger.Info(event.String())

	oracleClaim, err := txs.ProphecyClaimToSignedOracleClaim(event, sub.PrivateKey)
	if err != nil {
		return err
	}
	return txs.RelayOracleClaimToEthereum(sub.EthProvider, contractAddress, types.LogNewProphecyClaim,
		oracleClaim, sub.PrivateKey)
}

// SendMsg wraps the msg in a stdtx, signs and sends it
func (sub EthereumSub) SendMsg(datagram sdk.Msg) (sdk.TxResponse, error) {
	return sub.SendMsgs([]sdk.Msg{datagram})
}

// SendMsgs wraps the msgs in a stdtx, signs and sends it
func (sub EthereumSub) SendMsgs(datagrams []sdk.Msg) (res sdk.TxResponse, err error) {
	var out []byte
	if out, err = sub.BuildAndSignTx(datagrams); err != nil {
		return res, err
	}
	return sub.BroadcastTxCommit(out)
}

// BuildAndSignTx takes messages and builds, signs and marshals a sdk.Tx to prepare it for broadcast
func (sub EthereumSub) BuildAndSignTx(datagram []sdk.Msg) ([]byte, error) {
	// Fetch account and sequence numbers for the account
	acc, err := auth.NewAccountRetriever(sub.Cdc, sub.CliCtx).GetAccount(sdk.AccAddress(sub.ValidatorAddress))
	if err != nil {
		return nil, err
	}

	gasCoins, err := sdk.ParseDecCoins(gasPrices)
	if err != nil {
		return nil, err
	}

	defer UseSDKContext()()
	txBldr := auth.NewTxBuilder(
		auth.DefaultTxEncoder(sub.CliCtx.Codec), acc.GetAccountNumber(),
		acc.GetSequence(), gas, gasAdjustment, false, sub.CosmosChainID,
		"", sdk.NewCoins(), gasCoins,
	)

	return txBldr.WithKeybase(sub.Keybase).BuildAndSign(sub.ValidatorName, ckeys.DefaultKeyPass, datagram)
}

// BroadcastTxCommit takes the marshaled transaction bytes and broadcasts them
func (sub EthereumSub) BroadcastTxCommit(txBytes []byte) (sdk.TxResponse, error) {
	res, err := sub.CliCtx.BroadcastTxCommit(txBytes)
	return res, err
}

func (sub EthereumSub) QueryABCI(req abci.RequestQuery) (res abci.ResponseQuery, err error) {
	opts := rpcclient.ABCIQueryOptions{
		Height: req.GetHeight(),
		Prove:  req.Prove,
	}

	result, err := sub.Client.ABCIQueryWithOptions(req.Path, req.Data, opts)
	if err != nil {
		// retry queries on EOF
		if strings.Contains(err.Error(), "EOF") {
			return sub.QueryABCI(req)
		}
		return res, err
	}

	if !result.Response.IsOK() {
		return res, errors.New(result.Response.Log)
	}

	return result.Response, nil
}

// QueryWithData satisfies auth.NodeQuerier interface and used for fetching account details
func (sub EthereumSub) QueryWithData(p string, d []byte) (byt []byte, i int64, err error) {
	var res abci.ResponseQuery
	if res, err = sub.QueryABCI(abci.RequestQuery{Path: p, Height: 0, Data: d}); err != nil {
		return byt, i, err
	}

	return res.Value, res.Height, nil
}

func newRPCClient(addr string, timeout time.Duration) (*rpchttp.HTTP, error) {
	httpClient, err := libclient.DefaultHTTPClient(addr)
	if err != nil {
		return nil, err
	}

	httpClient.Timeout = timeout
	rpcClient, err := rpchttp.NewWithClient(addr, "/websocket", httpClient)
	if err != nil {
		return nil, err
	}

	return rpcClient, nil
}

// TODO: We'll want this for chains with prefixes other than 'cosmos'
var sdkContextMutex sync.Mutex

// UseSDKContext uses a custom Bech32 account prefix and returns a restore func
func UseSDKContext() func() {
	// Return a function that resets and unlocks.
	return func() {
		// defer sdkContextMutex.Unlock()
		// config.SetBech32PrefixForAccount(account, pubaccount)
	}
}
