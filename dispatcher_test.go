// v0.1.1
// Author: Wunderbarb
// ©, Jan 2024

package msj

import (
	"context"
	"errors"
	"github.com/wunderbarb/multistatejob/pkg/go_test"
	"testing"
	"time"
)

func TestNewDispatchEngine(t *testing.T) {
	require, assert := go_test.Describe(t)

	nde := NewDispatchEngineInput{}
	_, err := New(nde)
	require.Error(err)
	nde.In = NewBasicEventQueue()
	de, err := New(nde)
	require.NoError(err)
	require.NotNil(de)
	assert.NotNil(de.jobs)
}

func TestDispatchEngine_HandleEvent(t *testing.T) {
	require, assert := go_test.Describe(t)

	a := 1
	fn := func(_ Job, _ []byte, _ *DispatchEngine, _ ...any) error {
		a++
		return nil
	}
	de, j, tt := randomDispatchEngineAndJob(fn)
	ev, _ := NewEvent(j.Number, tt, nil)
	require.NoError(de.HandleEvent(ev))
	assert.Equal(2, a)
	ev1, _ := NewEvent(j.Number+1, tt, nil)
	assert.Error(de.HandleEvent(ev1))
	ev2, _ := NewEvent(j.Number, tt+1, nil)
	assert.Error(de.HandleEvent(ev2))
}

func TestDispatchEngine_getHandler(t *testing.T) {
	require, assert := go_test.Describe(t)

	fn := func(_ Job, _ []byte, _ *DispatchEngine, _ ...any) error {
		return nil
	}
	tt := go_test.Rng.Int()
	s := JobState(go_test.Rng.Int())
	m := MapState{
		s: MapEventFunction{
			tt: fn,
		},
	}
	de := DispatchEngine{
		ms: m,
	}
	_, err := de.getHandler(s, tt)
	require.NoError(err)
	_, err = de.getHandler(s+1, tt)
	assert.Error(err)
	_, err = de.getHandler(s, tt+1)
	assert.Error(err)
}

func TestDispatchEngine_Run(t *testing.T) {
	require, assert := go_test.Describe(t)
	go_test.NoLeak(t)

	a := 1
	fn := func(_ Job, _ []byte, _ *DispatchEngine, _ ...any) error {
		a++
		return nil
	}
	de, tt := randomDispatchEngine(fn)

	de.Run(context.Background())
	j, _ := NewJob(go_test.Rng.Int63(), nil)
	isPanic(de.jobs.Add(j))
	ev, _ := NewEvent(j.Number, tt, nil)
	require.NoError(de.inQueue.Add(ev))
	assert.Eventually(func() bool {
		return a == 2
	}, time.Second, 20*time.Millisecond)
	de.Shutdown()
}

func TestDispatchEngine_JobCompleted(t *testing.T) {
	require, assert := go_test.Describe(t)

	a := 1
	fn := func(_ Job, _ []byte, _ *DispatchEngine, _ ...any) error {
		a++
		return nil
	}
	de, j, _ := randomDispatchEngineAndJob(fn)

	require.NoError(de.JobCompleted(j))
	j1, err := de.jobs.GetJob(j.Number)
	require.NoError(err)
	assert.Equal(JobCompleted, j1.State)
	assert.WithinDuration(time.Now(), j1.Ended, time.Millisecond)
	oq := NewBasicEventQueue()
	de.outQueue = oq
	j2, _ := NewJob(go_test.Rng.Int63(), nil)
	isPanic(de.jobs.Add(j2))
	require.NoError(de.JobCompleted(j2))
	j1, err = de.jobs.GetJob(j2.Number)
	require.NoError(err)
	assert.Equal(JobCompleted, j1.State)
	assert.WithinDuration(time.Now(), j1.Ended, time.Millisecond)
	ev, err := oq.Next()
	require.NoError(err)
	assert.Equal(EventCompleted, ev.Type)
	assert.Equal(j2.Number, ev.Job)
}

func TestDispatchEngine_JobFailed(t *testing.T) {
	require, assert := go_test.Describe(t)

	a := 1
	fn := func(_ Job, _ []byte, _ *DispatchEngine, _ ...any) error {
		a++
		return nil
	}
	de, j, _ := randomDispatchEngineAndJob(fn)

	err1 := errors.New("for the test")
	require.NoError(de.JobFailed(j, err1))
	j1, err := de.jobs.GetJob(j.Number)
	require.NoError(err)
	assert.Equal(JobFailed, j1.State)
	assert.WithinDuration(time.Now(), j1.Ended, time.Millisecond)
	oq := NewBasicEventQueue()
	de.outQueue = oq
	j2, _ := NewJob(go_test.Rng.Int63(), nil)
	isPanic(de.jobs.Add(j2))
	require.NoError(de.JobFailed(j2, err1))
	j1, err = de.jobs.GetJob(j2.Number)
	require.NoError(err)
	assert.Equal(JobFailed, j1.State)
	assert.WithinDuration(time.Now(), j1.Ended, time.Millisecond)
	ev, err := oq.Next()
	require.NoError(err)
	assert.Equal(EventFailed, ev.Type)
	assert.Equal(j2.Number, ev.Job)
	var msg1 JobFailedMsg
	require.NoError(ev.GetPayload(&msg1))
	assert.Equal(err1.Error(), msg1.GetMsg())
}
func randomDispatchEngine(fn EventHandlerFunction) (DispatchEngine, int) {
	tt := go_test.Rng.Int()
	m := MapState{
		JobPending: MapEventFunction{
			tt: fn,
		},
	}

	nde := NewDispatchEngineInput{
		In:   NewBasicEventQueue(),
		Out:  nil,
		Jobs: nil,
		Ms:   m,
	}
	de, err := New(nde)
	isPanic(err)
	return de, tt
}

func randomDispatchEngineAndJob(fn EventHandlerFunction) (DispatchEngine, Job, int) {
	de, tt := randomDispatchEngine(fn)
	j, _ := NewJob(go_test.Rng.Int63(), nil)
	isPanic(de.jobs.Add(j))
	return de, j, tt
}
