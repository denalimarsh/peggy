package events

// -----------------------------------------------------
// FinalityQueue : stores events in memory for processing once
//				   they've reached the finality threshold
// -----------------------------------------------------

// FinalityThreshold : constant threshold of 6 blocks
const FinalityThreshold = 6

// Queue : is a basic FIFO queue based on a circular list that resizes as needed
type Queue struct {
	events []*LockEvent
	size   int
	head   int
	tail   int
	count  int
}

// NewQueue : returns a new queue with the given initial size
func NewQueue(size int) *Queue {
	return &Queue{
		events: make([]*LockEvent, size),
		size:   size,
	}
}

// Push : adds an event to the queue
func (q *Queue) Push(n *LockEvent) {
	if q.head == q.tail && q.count > 0 {
		events := make([]*LockEvent, len(q.events)+q.size)
		copy(events, q.events[q.head:])
		copy(events[len(q.events)-q.head:], q.events[:q.head])
		q.head = 0
		q.tail = len(q.events)
		q.events = events
	}
	q.events[q.tail] = n
	q.tail = (q.tail + 1) % len(q.events)
	q.count++
}

// Peek : returns the first element from the queue
func (q *Queue) Peek() *LockEvent {
	if q.count == 0 {
		return nil
	}
	event := q.events[q.head]
	q.head = (q.head + 1) % len(q.events)
	q.count--
	return event
}

// Pop : removes and returns an event from the queue in first to last order
func (q *Queue) Pop() *LockEvent {
	if q.count == 0 {
		return nil
	}
	event := q.events[q.head]
	q.head = (q.head + 1) % len(q.events)
	q.count--
	return event
}

// Size : get the current size of the queue
func (q *Queue) Size() int {
	return len(q.events)
}

// IsEmpty : check if the queue has any events
func (q *Queue) IsEmpty() bool {
	return len(q.events) == 0
}

// IsEventProcessing : get the processing status of any event in the queue
func (q *Queue) IsEventProcessing(txHash string) bool {
	for i := 0; i < len(q.events); i++ {
		if q.events[i].TxHash == txHash {
			return true
		}
	}
	return false
}

// GetFinalEvents : Return all the events which have passed the finality threshold
func (q *Queue) GetFinalEvents(currentBlockNumber uint64) ([]LockEvent, int) {
	// Valid events must have occured at least 6 blocks ago
	validEvents := make([]LockEvent, len(q.events))

	for q.Peek().BlockNumber+FinalityThreshold > currentBlockNumber {
		// Add the event to the validated event array
		// validEvent := LockEvent{}
		validEvent := q.Pop()
		validEvents = append(validEvents, validEvent)
	}

	return validEvents, len(validEvents)
}

// func main() {
// 	q := NewQueue(1)
// 	q.Push(&LockEvent{4})
// 	q.Push(&LockEvent{5})
// 	q.Push(&LockEvent{6})
// 	fmt.Println(q.Pop(), q.Pop(), q.Pop())
// }
