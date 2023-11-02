package action

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNoDependenciesError_IsError(t *testing.T) {
	err := NoDependenciesError{
		msg: "error",
	}
	require.Error(t, err)
}

func TestNoDependenciesError_Error(t *testing.T) {
	err := NoDependenciesError{
		msg: "error",
	}
	assert.Equal(t, "error", err.Error())
}
