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

	// peggy "/src/github.com/cosmos/peggy/testnet-contracts/contracts"
	// peggy "../../../testnet-contracts/contracts" // for demo
	peggy "github.com/cosmos/peggy/testnet-contracts/contracts"
)

// const peggy = "/src/github.com/cosmos/peggy/testnet-contracts/contracts"

func InitCosmosRelayer(peggyContractAddress common.Address, rawPrivateKey string) error {

	provider := "http://localhost:7545"

	// Star our ethereum client
	client, err := ethclient.Dial(provider)
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA(rawPrivateKey)
	if err != nil {
		log.Fatal(err)
	}

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

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(5)     // in wei
	auth.GasLimit = uint64(300000) // 300,000 Gwei in units
	auth.GasPrice = gasPrice

	// address := common.HexToAddress(peggyContractAddress)
	instance, err := peggy.NewPeggy(peggyContractAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	recipient := []byte{}
	copy(recipient[:], []byte("cosmos1gn8409qq9hnrxde37kuxwx5hrxpfpv8426szuv"))

	tokenAddressString := "0x7B95B6EC7EbD73572298cEf32Bb54FA408207359"
	if !common.IsHexAddress(tokenAddressString) {
		return fmt.Errorf("Invalid contract-address: %v", tokenAddressString)
	}
	tokenAddress := common.HexToAddress(tokenAddressString)

	amount := big.NewInt(0)
	amount.SetBytes([]byte("50"))

	tx, err := instance.Lock(auth, recipient, tokenAddress, amount)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx sent: %s", tx.Hash().Hex()) // tx sent: 0x8d490e535678e9a24360e955d75b27ad307bdfb97a1dca51d0f3035dcee3e870

	// result, err := instance.Items(nil, key)
	// if err != nil {
	//     log.Fatal(err)
	// }

	// fmt.Println(string(result[:])) // "bar"

	return nil
}
