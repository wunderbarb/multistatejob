package msj

import (
	"github.com/wunderbarb/multistatejob/testdata"
	"testing"
	"time"

	"github.com/wunderbarb/multistatejob/pkg/go_test"
)

//go:generate protoc  -I=. --proto_path=testdata/ --go_out=testdata/ --go_opt=paths=source_relative msg.proto

func TestNewJob(t *testing.T) {
	require, assert := go_test.Describe(t)

	tt := go_test.Rng.Int63()
	msg := go_test.RandomID()
	msg1 := testdata.Simple{Msg: msg}

	j, err := NewJob(tt, &msg1)
	require.NoError(err)
	assert.NotZero(j.Number)
	assert.Equal(JobPending, j.State)
	assert.WithinDuration(time.Now(), j.Initiated, time.Millisecond)

	var m testdata.Simple
	err = j.GetPayload(&m)
	require.NoError(err)
	assert.Equal(msg, m.Msg)
}

func TestJob_Update(t *testing.T) {
	require, assert := go_test.Describe(t)

	msg := go_test.RandomID()
	msg1 := testdata.Simple{Msg: msg}
	j, _ := NewJob(go_test.Rng.Int63(), &msg1)

	s := JobState(go_test.Rng.Int31())
	require.NoError(j.Update(s, nil))
	assert.Equal(s, j.State)
	var msg2 testdata.Simple
	require.NoError(j.GetPayload(&msg2))
	assert.Equal(msg, msg2.GetMsg())
	assert.Equal(s, j.State)

	msg22 := go_test.RandomID()
	msg21 := testdata.Simple{Msg: msg22}
	s = JobState(go_test.Rng.Int31())

	require.NoError(j.Update(s, &msg21))
	assert.Equal(s, j.State)
	require.NoError(j.GetPayload(&msg2))
	assert.Equal(msg22, msg2.GetMsg())

}
