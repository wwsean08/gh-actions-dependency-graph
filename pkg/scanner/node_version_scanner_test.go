package scanner

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wwsean08/actions-dependency-graph/pkg/action"
)

func TestIsNodeVersionEOL(t *testing.T) {
	type testMap struct {
		using       string
		expected    bool
		errExpected bool
	}
	testData := []testMap{
		{
			using:       "node16",
			expected:    true,
			errExpected: false,
		},
		{
			using:       "node20",
			expected:    false,
			errExpected: false,
		},
		{
			using:       "composite",
			expected:    false,
			errExpected: true,
		},
		{
			using:       "docker",
			expected:    false,
			errExpected: true,
		},
	}

	for _, test := range testData {
		act := new(action.Action)
		act.Runs = new(action.RunsBlock)
		act.Runs.Using = test.using
		scan := &Scanner{}
		vuln, err := scan.IsNodeVersionEOL(act)
		if test.errExpected {
			assert.Error(t, err, fmt.Sprintf("Expected error for %v", test))
		} else {
			assert.NoError(t, err, fmt.Sprintf("Expected no error for %v", test))
		}

		assert.Equal(t, test.expected, vuln, fmt.Sprintf("Expected %t for %v", test.expected, test))
	}
}
