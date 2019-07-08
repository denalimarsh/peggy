package events

// -----------------------------------------------------
// 	Events: Events maintains a mapping of events to an array
//		of claims made by validators.
// -----------------------------------------------------

import (
	"fmt"
)

// FinalityQueue : stores events in memory for processing once they've
//				   reached the finality threshold
// var FinalityQueue = make([]*LockEvent, 0)

// // Push : add an element to the back of queue
// func Push(event LockEvent) {
// 	FinalityQueue = append(FinalityQueue, &event)
// }

// // Peek : check the top element in the queue
// func Peek() *LockEvent {
// 	element := FinalityQueue[0]
// 	return element
// }

// // Pop : delete the element from the front of the queue and return it
// func Pop() LockEvent {
// 	element := FinalityQueue[0]
// 	FinalityQueue = FinalityQueue[1:]
// 	return element
// }

// Size : get the current size of the queue
// func Size() int {
// 	return len(FinalityQueue)
// }

// // IsEmpty : check if the queue has any events
// func IsEmpty() bool {
// 	return len(FinalityQueue) == 0
// }

// // IsEventProcessing : get the processing status of any event in the queue
// func IsEventProcessing(txHash string) bool {
// 	for i := 0; i < len(FinalityQueue); i++ {
// 		if FinalityQueue[i].TxHash == txHash {
// 			return true
// 		}
// 	}
// 	return false
// }

// // GetValidEvents : Return all the events which have passed the finality threshold
// func GetValidEvents(currentBlockNumber uint64) ([]LockEvent, int) {
// 	// Valid events must have occured at least 6 blocks ago
// 	validEvents := make([]LockEvent, len(FinalityQueue))

// 	for FinalityQueue[0].BlockNumber+6 > currentBlockNumber {
// 		// Add the event to the validated event array
// 		// validEvent := LockEvent{}
// 		validEvents = append(validEvents, Pop())
// 	}

// 	return validEvents, len(validEvents)
// }

// // PrintEvents : prints all the claims made on this event
// func PrintProcessingEvents() {
// 	for index, event := range FinalityQueue {
// 		PrintEvent(event)
// 	}
// }

// EventRecords : map of transaction hashes to LockEvent structs
var EventRecords = make(map[string]LockEvent)

// NewEventWrite : add a validator's address to the official claims list
func NewEventWrite(txHash string, event LockEvent) {
	EventRecords[txHash] = event
}

// IsEventRecorded : checks the sessions stored events for this transaction hash
func IsEventRecorded(txHash string) bool {
	return EventRecords[txHash].Nonce != nil
}

// PrintEventByTx : prints any witnessed events associated with a given transaction hash
func PrintEventByTx(txHash string) {
	if IsEventRecorded(txHash) {
		PrintEvent(EventRecords[txHash])
	} else {
		fmt.Printf("\nNo records from this sesson for tx: %v\n", txHash)
	}
}

// PrintEvents : prints all the claims made on this event
func PrintEvents() {

	// For each claim, print the validator which submitted the claim
	for txHash, event := range EventRecords {
		fmt.Printf("\nTransaction: %v\n", txHash)
		PrintEvent(event)
	}
}
