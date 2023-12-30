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
