package contract

import (
	"io/ioutil"
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestLoadABI : test that contract containing named event is successfully loaded
func TestLoadABI(t *testing.T) {

	//Get the ABI ready
	rawContractAbi, errorMsg := ioutil.ReadFile("./abis/Peggy.abi")
	if errorMsg != nil {
		log.Fatal(errorMsg)
	}

	require.True(t, strings.Contains(string(rawContractAbi), "LogLock"))
}
