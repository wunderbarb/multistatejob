// v0.2.0
// Author: DIEHL E.
// © Dec 2023

package msj

import (
	"errors"
)

var (
	// ErrQueueEmpty occurs when the queue is empty.
	ErrQueueEmpty = errors.New("queue is empty")
)

// Queue is the interface of a queue.
type Queue interface {
	// Add pushes an event to the queue.
	Add(Event) error
	// Next pops an event from the queue. It waits until an event appears.
	Next() (Event, error)
	// TryNext pops an event from the queue if present.  It is not blocking.  If there is no event, it returns false.
	TryNext() (Event, bool)
}

// BasicEventQueue is a type that represents a basic event queue. It maintains a list of events and provides methods
// compliant with the Queue interface.
// It uses a buffered channel to convey the events.
type BasicEventQueue struct {
	ch chan Event
}

// NewBasicEventQueue creates a new instance of BasicEventQueue.
func NewBasicEventQueue() *BasicEventQueue {
	const cSizeQueue = 100
	return &BasicEventQueue{
		ch: make(chan Event, cSizeQueue),
	}
}

// Add adds an event to the BasicEventQueue.
// The event is added to the underlying channel of the BasicEventQueue, allowing it to be processed.
// It returns an error if the addition fails.
func (q *BasicEventQueue) Add(e Event) error {
	q.ch <- e
	return nil
}

func (q *BasicEventQueue) Next() (Event, error) {
	e := <-q.ch
	return e, nil
}

func (q *BasicEventQueue) TryNext() (Event, bool) {
	select {
	case e := <-q.ch:
		return e, true
	default:
		return Event{}, false
	}
}
