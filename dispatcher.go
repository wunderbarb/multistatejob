// v0.1.0

package msj

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
)

var (

	// ErrUnknownEventType represents an error indicating an unknown event type.
	ErrUnknownEventType = errors.New("unknown event type")
	// ErrUnknownState represents an error indicating an unknown job state.
	ErrUnknownState = errors.New("unknown state")
)

// EventHandlerFunction represents a function type that handles an event for a specific job.
// It takes the job (type `Job`), the payload (type `[]byte`), the dispatch engine (type `*DispatchEngine`),
// and any additional arguments and returns an error if it fails
// It is called for a given type of Event.
type EventHandlerFunction func(Job, []byte, *DispatchEngine, ...any) error

type MapEventFunction map[int]EventHandlerFunction

type MapState map[JobState]MapEventFunction

// DispatchEngine is a struct that represents a dispatch engine,
// responsible for managing queues and job operations.
type DispatchEngine struct {
	inQueue  Queue
	outQueue Queue
	jobs     JobHolder
	ms       MapState
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewDispatchEngineInput is a struct that represents a dispatch engine,
// responsible for managing queues, job operations, and map state.
type NewDispatchEngineInput struct {
	In   Queue
	Out  Queue
	Jobs JobHolder
	Ms   MapState
}

func New(nde NewDispatchEngineInput) (DispatchEngine, error) {
	if nde.In == nil {
		return DispatchEngine{}, errors.New("input queue is required")
	}
	if nde.Jobs == nil {
		nde.Jobs = NewJobHolderAsMap()
	}
	return DispatchEngine{
		inQueue:  nde.In,
		outQueue: nde.Out,
		jobs:     nde.Jobs,
		ms:       nde.Ms,
	}, nil
}

func (de *DispatchEngine) Run(ctx context.Context, a ...any) {
	de.ctx, de.cancel = context.WithCancel(ctx)
	go func() {
		for {
			select {
			case <-de.ctx.Done():
				return
			default:
				ev, ok := de.inQueue.TryNext()
				if ok {
					err := de.HandleEvent(ev)
					if err != nil {
						fmt.Printf("err %v\n", err)
					}
				}
			}
		}
	}()
}

func (de *DispatchEngine) Shutdown() {
	de.cancel()
}

// HandleEvent is a method of the DispatchEngine type that handles an event by executing the corresponding handler function.
//
// It takes an Event as the first parameter and accepts variadic arguments of any type.
// The Event contains information about the job and event type, as well as the payload data.
//
// The method starts by retrieving the corresponding job from the JobHolder using the job ID from the Event.
// If the job is not found, it returns an error.
//
// It then retrieves the handler function from the MapState based on the job state and event type.
// If the handler function is not found, it returns an error.
//
// Finally, it calls the handler function with the job, event payload, DispatchEngine instance, and additional arguments passed to HandleEvent.
//
// The method returns any error that occurred during the execution of the handler function.
func (de *DispatchEngine) HandleEvent(ev Event, a ...any) error {
	j, err := de.jobs.GetJob(ev.Job)
	if err != nil {
		return err
	}
	h, err := de.getHandler(j.State, ev.Type)
	if err != nil {
		return err
	}
	return h(j, ev.Payload, de, a...)
}

func (de *DispatchEngine) getHandler(state JobState, eventType int) (EventHandlerFunction, error) {
	m, ok := de.ms[state]
	if !ok {
		return nil, ErrUnknownState
	}
	h, ok := m[eventType]
	if !ok {
		return nil, ErrUnknownEventType
	}
	return h, nil
}
