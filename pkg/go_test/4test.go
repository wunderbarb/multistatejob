// v0.1.0

package go_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ysmood/gotrace"
	"strings"
	"testing"
)

var testCounter int

// Describe displays the rank of the test, the name of the function
// and its optional description provided by 'msg'.  It initializes an assert
// and a require function and returns them.
func Describe(t *testing.T, msg ...string) (*require.Assertions,
	*assert.Assertions) {

	dispMsg := ""
	if len(msg) != 0 {
		dispMsg = msg[0]
	}
	name := strings.TrimPrefix(strings.TrimPrefix(t.Name(), "Test"), "_")
	fmt.Printf("Test %d: %s %s\n", testCounter, name, dispMsg)
	testCounter++
	return require.New(t), assert.New(t)
}

// NoLeak verifies whether there was no new running go routines after the current test.
func NoLeak(t *testing.T) {
	ign := gotrace.IgnoreCurrent()
	gotrace.CheckLeak(t, 0, ign)
}
