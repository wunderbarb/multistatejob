package msj

import (
	"github.com/wunderbarb/multistatejob/pkg/go_test"
	"github.com/wunderbarb/multistatejob/testdata"
	"testing"
)

func TestNewEvent(t *testing.T) {
	require, assert := go_test.Describe(t)

	tt := go_test.Rng.Int()
	msg := go_test.RandomID()
	msg1 := testdata.Simple{Msg: msg}

	e, err := NewEvent(0, tt, &msg1)
	require.NoError(err)
	assert.Equal(e.Type, tt)

	var m testdata.Simple
	err = e.GetPayload(&m)
	require.NoError(err)
	assert.Equal(msg, m.Msg)
}
