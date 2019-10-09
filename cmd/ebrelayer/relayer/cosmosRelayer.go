package relayer

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	// TODO: Refactor solc generation of Peggy.go
	peggy "github.com/cosmos/peggy/testnet-contracts/contracts"
)

// TODO: These testing constants will be removed and replaced with Lock event attributes
const (
	WeiAmount = 100
	// CosmosRecipient : hashed address "cosmos1gn8409qq9hnrxde37kuxwx5hrxpfpv8426szuv"
	CosmosRecipient      = "0x636F736D6F7331676E38343039717139686E7278646533376B75787778356872787066707638343236737A7576"
	EthereumTokenAddress = "0x0000000000000000000000000000000000000000"
)

// InitCosmosRelayer : initalizes a relayer which witnesses events on the Cosmos network and relays them to Ethereum
func InitCosmosRelayer(peggyContractAddress common.Address, rawPrivateKey string) error {

	// TODO: Parameterize the provider
	provider := "http://localhost:7545"

	// Start Ethereum client
	client, err := ethclient.Dial(provider)
	if err != nil {
		log.Fatal(err)
	}

	// Parse private key
	privateKey, err := crypto.HexToECDSA(rawPrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	// Parse public key
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Set up tx signature authorization
	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(WeiAmount) // in wei
	auth.GasLimit = uint64(300000)     // 300,000 Gwei in units
	auth.GasPrice = gasPrice

	// Initialize Peggy contract instance
	instance, err := peggy.NewPeggy(peggyContractAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	// Event parameters
	recipient := []byte{}
	copy(recipient[:], []byte(CosmosRecipient))

	tokenAddressString := EthereumTokenAddress
	if !common.IsHexAddress(tokenAddressString) {
		return fmt.Errorf("Invalid contract-address: %v", tokenAddressString)
	}
	tokenAddress := common.HexToAddress(tokenAddressString)

	amount := big.NewInt(WeiAmount)

	// Send transaction to the instance's specified method
	tx, err := instance.Lock(auth, recipient, tokenAddress, amount)
	if err != nil {
		log.Fatal(err)
	}

	// Get the transaction receipt
	receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Tx hash:", tx.Hash().Hex())
	fmt.Println("Status:", receipt.Status, "\n")

	return nil
}
