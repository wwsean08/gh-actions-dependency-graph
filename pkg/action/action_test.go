package action

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseActionReturnsErrOnNonexistentFile(t *testing.T) {
	action, err := ParseAction("testdata/non-existent-file.404")
	require.Error(t, err)
	require.Nil(t, action)
}

func TestParseActionParsesCompositeActionCorrectly(t *testing.T) {
	action, err := ParseAction("testdata/composite-example.yaml")
	require.NoError(t, err)
	require.NotNil(t, action)

	assert.True(t, action.IsComposite())
}

func TestParseActionParsesDockerActionCorrectly(t *testing.T) {
	action, err := ParseAction("testdata/docker-example.yaml")
	require.NoError(t, err)
	require.NotNil(t, action)

	assert.True(t, action.IsDocker())
}

func TestParseActionParsesJavascriptActionCorrectly(t *testing.T) {
	action, err := ParseAction("testdata/js-example.yaml")
	require.NoError(t, err)
	require.NotNil(t, action)

	assert.True(t, action.IsJavascript())
}

func TestAction_IsComposite(t *testing.T) {
	type testData struct {
		using    string
		expected bool
		message  string
	}
	tests := []testData{
		{
			using:    "composite",
			expected: true,
			message:  "expected composite to resolve as true",
		},
		{
			using:    "node16",
			expected: false,
			message:  "expected node16 to resolve as false",
		},
		{
			using:    "node20",
			expected: false,
			message:  "expected node20 to resolve as false",
		},
		{
			using:    "docker",
			expected: false,
			message:  "expected docker to resolve as false",
		},
	}

	for _, test := range tests {
		action := Action{
			Runs: &RunsBlock{
				Using: test.using,
			},
		}
		assert.Equal(t, test.expected, action.IsComposite(), test.message)
	}
}

func TestAction_IsDocker(t *testing.T) {
	type testData struct {
		using    string
		expected bool
		message  string
	}
	tests := []testData{
		{
			using:    "composite",
			expected: false,
			message:  "expected composite to resolve as false",
		},
		{
			using:    "node16",
			expected: false,
			message:  "expected node16 to resolve as false",
		},
		{
			using:    "node20",
			expected: false,
			message:  "expected node20 to resolve as false",
		},
		{
			using:    "docker",
			expected: true,
			message:  "expected docker to resolve as true",
		},
	}

	for _, test := range tests {
		action := Action{
			Runs: &RunsBlock{
				Using: test.using,
			},
		}
		assert.Equal(t, test.expected, action.IsDocker(), test.message)
	}
}

func TestAction_IsJavascript(t *testing.T) {
	type testData struct {
		using    string
		expected bool
		message  string
	}
	tests := []testData{
		{
			using:    "composite",
			expected: false,
			message:  "expected composite to resolve as false",
		},
		{
			using:    "node16",
			expected: true,
			message:  "expected node16 to resolve as true",
		},
		{
			using:    "node20",
			expected: true,
			message:  "expected node20 to resolve as true",
		},
		{
			using:    "docker",
			expected: false,
			message:  "expected docker to resolve as false",
		},
	}

	for _, test := range tests {
		action := Action{
			Runs: &RunsBlock{
				Using: test.using,
			},
		}
		assert.Equal(t, test.expected, action.IsJavascript(), test.message)
	}
}

func TestAction_GetNodeVersion(t *testing.T) {
	type testData struct {
		using       string
		expected    int
		shouldError bool
		message     string
	}
	tests := []testData{
		{
			using:       "composite",
			expected:    0,
			shouldError: true,
			message:     "expected composite to resolve as 0, error",
		},
		{
			using:       "node16",
			expected:    16,
			shouldError: false,
			message:     "expected node16 to resolve as 16, nil",
		},
		{
			using:       "node20",
			expected:    20,
			shouldError: false,
			message:     "expected node20 to resolve as 20, nil",
		},
		{
			using:       "nodefoo",
			expected:    0,
			shouldError: true,
			message:     "expected node20 to resolve as 0, error",
		},
		{
			using:       "docker",
			expected:    0,
			shouldError: true,
			message:     "expected docker to resolve as 0, error",
		},
	}

	for _, test := range tests {
		action := Action{
			Runs: &RunsBlock{
				Using: test.using,
			},
		}
		version, err := action.GetNodeVersion()
		if test.shouldError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
		assert.Equal(t, test.expected, version)
	}
}
