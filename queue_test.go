// v0.1.2
// Author: DIEHL E.
// © Dec 2023

package msj

import (
	"github.com/wunderbarb/multistatejob/pkg/go_test"
	"testing"
)

func TestBasicEventQueue_Next(t *testing.T) {
	require, assert := go_test.Describe(t)

	q := NewBasicEventQueue()
	e1 := Event{
		Job:  go_test.Rng.Uint64(),
		Type: go_test.Rng.Int(),
	}
	require.NoError(q.Add(e1))
	e2, err := q.Next()
	require.NoError(err)
	assert.Equal(e1.Job, e2.Job)
	assert.Equal(e1.Type, e2.Type)

}

func TestBasicEventQueue_TryNext(t *testing.T) {
	require, assert := go_test.Describe(t)

	e1 := Event{
		Job:  go_test.Rng.Uint64(),
		Type: go_test.Rng.Int(),
	}
	q := NewBasicEventQueue()
	_, ok := q.TryNext()
	require.False(ok)
	isPanic(q.Add(e1))
	e2, ok := q.TryNext()
	require.True(ok)
	assert.Equal(e1.Job, e2.Job)
	assert.Equal(e1.Type, e2.Type)
}

func isPanic(err error) {
	if err != nil {
		panic(err)
	}
}
