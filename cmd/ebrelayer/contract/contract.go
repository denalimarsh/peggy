package contract

// -------------------------------------------------------
//    Contract : Contains functionality for loading the
//				 smart contract
// -------------------------------------------------------

import (
	"go/build"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// AbiPath : path to the file containing the smart contract's ABI
const AbiPath = "/src/github.com/cosmos/peggy/cmd/ebrelayer/contract/abi/Peggy.abi"

// LoadABI : loads a smart contract as an abi.ABI
func LoadABI() abi.ABI {
	// Open the file containing Peggy contract's ABI
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}

	rawContractAbi, err := ioutil.ReadFile(gopath + AbiPath)
	if err != nil {
		panic(err)
	}

	// Convert the raw abi into a usable format
	contractAbi, err := abi.JSON(strings.NewReader(string(rawContractAbi)))
	if err != nil {
		panic(err)
	}

	return contractAbi
}

// ParseEventSignatures : parse event signatures from compiled contract ABI
func ParseEventSignatures(contractABI abi.ABI) map[string]common.Hash {
	eventSigs := make(map[string]common.Hash)
	contractEvents := contractABI.Events

	for event := range contractEvents {
		eventName := contractEvents[event].Name
		id := contractEvents[event].Id()
		eventSigs[eventName] = id
	}

	return eventSigs
}
