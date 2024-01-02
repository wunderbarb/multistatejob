// v0.1.1
// Author: Wunderbarb

package msj

import (
	"context"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"time"
)

var (

	// ErrInputQueueNeeded represents an error indicating the need for an input queue.
	ErrInputQueueNeeded = errors.New("input queue requested")
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
	log      *zap.Logger
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewDispatchEngineInput is a struct that represents the input parameters for creating a new DispatchEngine.
// In is a mandatory interface that represents a queue for incoming events.
// Out is an optional interface that represents a queue for outgoing events.
// Jobs is an interface that represents a holder for jobs.
// Ms is a map that stores event handler functions based on the job state.
// Log is an optional logger for logging.
type NewDispatchEngineInput struct {
	In   Queue
	Out  Queue
	Jobs JobHolder
	Ms   MapState
	Log  *zap.Logger
}

// New is a function that creates a new instance of DispatchEngine.
// It takes a NewDispatchEngineInput struct as the input parameter, which contains the necessary configuration options for initializing the DispatchEngine instance.
// If the Jobs object is not provided, it creates a new JobHolderAsMap.
// If the Log object is not provided, it creates a new Logger with the filename "dispatcher.log" and WithVerbose option.
// Finally, it creates a new instance of DispatchEngine by assigning the input parameters to the corresponding fields of DispatchEngine struct and returns it along with nil error.
func New(nde NewDispatchEngineInput) (DispatchEngine, error) {
	if nde.In == nil {
		return DispatchEngine{}, ErrInputQueueNeeded
	}
	if nde.Jobs == nil {
		nde.Jobs = NewJobHolderAsMap()
	}
	if nde.Log == nil {
		l, err := NewLogger("dispatcher.log", WithVerbose())
		if err != nil { // SHOULD NEVER HAPPEN
			return DispatchEngine{}, err
		}
		nde.Log = l.Zap()
	}
	return DispatchEngine{
		inQueue:  nde.In,
		outQueue: nde.Out,
		jobs:     nde.Jobs,
		ms:       nde.Ms,
		log:      nde.Log,
	}, nil
}

func (de *DispatchEngine) JobCompleted(j Job) error {
	j.State = JobCompleted
	j.Ended = time.Now()
	de.log.Info("completed", zap.Uint64("job", j.Number))
	err := de.jobs.UpdateJob(j)
	if err != nil {
		return err
	}
	if de.outQueue == nil {
		return nil
	}
	ev, err := NewEvent(j.Number, EventCompleted, nil)
	if err != nil {
		// SHOULD NEVER HAPPEN
		return err
	}
	err = de.outQueue.Add(ev)
	if err != nil {
		// SHOULD NEVER HAPPEN
		return err
	}

	return nil
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
					de.log.Info("HandleEvent", zap.Int("event type", ev.Type),
						zap.Uint64("job", ev.Job))
					err := de.HandleEvent(ev, a...)
					if err != nil {
						de.log.Error("HandleEvent", zap.Error(err), zap.Int("event type", ev.Type),
							zap.Uint64("job", ev.Job))
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
