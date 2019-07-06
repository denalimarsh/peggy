package events

// -----------------------------------------------------
//    Event : Creates LockEvents from new events on the ethereum
//			  Ethereum blockchain.
// -----------------------------------------------------

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// LockEvent : struct which represents a single smart contract event
type LockEvent struct {
	Id          [32]byte
	From        common.Address
	To          []byte
	Token       common.Address
	Value       *big.Int
	Nonce       *big.Int
	BlockNumber uint64
	TxHash      string
}

// NewLockEvent : parses LogLock events using go-ethereum's accounts/abi library
func NewLockEvent(contractAbi abi.ABI, eventName string, eventData []byte, blockNumber uint64, txHash string) LockEvent {
	// Check event name
	if eventName != "LogLock" {
		log.Fatal("Only LogLock events are currently supported.")
	}

	// Declare new LockEvent
	event := LockEvent{}

	// Set the event's block number and transaction hash
	event.BlockNumber = blockNumber
	event.TxHash = txHash

	// Parse the event's attributes as Ethereum network variables
	err := contractAbi.Unpack(&event, eventName, eventData)
	if err != nil {
		log.Fatalf("Unpacking: %v", err)
	}

	PrintEvent(event)
	return event
}

// PrintEvent : prints a LockEvent struct's information
func PrintEvent(event LockEvent) {

	// Extract variables into print-friendly format
	id := hex.EncodeToString(event.Id[:])
	sender := event.From.Hex()
	recipient := string(event.To[:])
	token := event.Token.Hex()

	// Print the event's information
	fmt.Printf("\nBlock Number: %v\nTx Hash: %v\nEvent ID: %v\nToken: %v\nSender: %v\nRecipient: %v\nValue: %v\nNonce: %v\n\n",
		event.BlockNumber, event.TxHash, id, token, sender, recipient, event.Value, event.Nonce)
}
