package action

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Some initial TDD
func TestParseActionReturnsErr(t *testing.T) {
	action, err := ParseAction("")
	require.Error(t, err)
	require.Nil(t, action)
}
