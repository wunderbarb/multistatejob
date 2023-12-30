// v0.1.1
// Author: Wunderbarb
// © Dec 2023

package msj

import (
	"github.com/wunderbarb/multistatejob/pkg/go_test"
	"testing"
)

func TestJobHolderAsMap_Add(t *testing.T) {
	require, assert := go_test.Describe(t)

	jhm := NewJobHolderAsMap()
	j1 := randomJob()
	require.NoError(jhm.Add(j1))
	j2, err := jhm.GetJob(j1.Number)
	require.NoError(err)
	assert.Equal(j1, j2)
	_, err = jhm.GetJob(j1.Number + 1)
	require.Error(err)
	assert.ErrorIs(jhm.Add(j1), ErrJobAlreadyExists)
}

func TestJobHolderAsMap_DeleteJob(t *testing.T) {
	require, assert := go_test.Describe(t)

	jhm := NewJobHolderAsMap()
	j1 := randomJob()
	isPanic(jhm.Add(j1))
	require.NoError(jhm.DeleteJob(j1.Number))
	_, err := jhm.GetJob(j1.Number)
	require.ErrorIs(err, ErrUnknownJob)
	assert.ErrorIs(jhm.DeleteJob(j1.Number), ErrUnknownJob)
}

func TestJobHolderAsMap_UpdateJob(t *testing.T) {
	require, assert := go_test.Describe(t)

	jhm := NewJobHolderAsMap()
	j1 := randomJob()
	isPanic(jhm.Add(j1))
	j1.State = JobState(go_test.Rng.Int())
	require.NoError(jhm.UpdateJob(j1))
	j2, err := jhm.GetJob(j1.Number)
	require.NoError(err)
	assert.Equal(j1, j2)
	assert.ErrorIs(jhm.UpdateJob(Job{Number: j1.Number + 1}), ErrUnknownJob)
}

func randomJob() Job {
	return Job{
		Number:  go_test.Rng.Uint64(),
		State:   JobState(go_test.Rng.Int()),
		Payload: go_test.RandomSlice(256),
	}
}
