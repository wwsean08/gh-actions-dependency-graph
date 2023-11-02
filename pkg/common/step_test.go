package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStep_GetCommitishReturnsExpectedVersion(t *testing.T) {
	type testData struct {
		uses     string
		expected string
		message  string
	}

	tests := []testData{
		{
			uses:     "actions/checkout@v4",
			expected: "v4",
			message:  "expectedPath actions/checkout@v4 to resolve as v4",
		},
		{
			uses:     "foo/bar/baz@main",
			expected: "main",
			message:  "expectedPath foo/bar/baz@main to resolve as main",
		},
		{
			uses:     "",
			expected: "",
			message:  "expectedPath empty string to resolve as empty string",
		},
		{
			uses:     "./foo/action.yaml",
			expected: "",
			message:  "expectedPath local action to not resolve a commitish",
		},
	}

	for _, test := range tests {
		step := Step{
			Uses: &test.uses,
		}
		assert.Equal(t, test.expected, step.GetGitRef(), test.message)
	}

	// Test that nil will return empty string
	step := Step{}
	assert.Equal(t, "", step.GetGitRef())
}

func TestStep_GetRepoReturnsExpectedRepo(t *testing.T) {
	type testData struct {
		uses     string
		expected string
		message  string
	}

	tests := []testData{
		{
			uses:     "actions/checkout@v4",
			expected: "actions/checkout",
			message:  "expectedPath actions/checkout@v4 to resolve as actions/checkout",
		},
		{
			uses:     "foo/bar/baz@main",
			expected: "foo/bar",
			message:  "expectedPath foo/bar/baz@main to resolve as foo/bar",
		},
		{
			uses:     "",
			expected: "",
			message:  "expectedPath empty string to resolve as empty string",
		},
	}

	for _, test := range tests {
		step := Step{
			Uses: &test.uses,
		}

		assert.Equal(t, test.expected, step.GetRepo(), test.message)
	}

	// test that nil uses returns empty string
	step := Step{}
	assert.Equal(t, "", step.GetRepo())
}

func TestStep_GetRepoPathReturnsExpectedPath(t *testing.T) {
	type testData struct {
		uses     string
		expected string
		message  string
	}

	tests := []testData{
		{
			uses:     "actions/checkout@v4",
			expected: "",
			message:  "expectedPath actions/checkout@v4 to resolve as empty string",
		},
		{
			uses:     "foo/bar/baz@main",
			expected: "baz",
			message:  "expectedPath foo/bar/baz@main to resolve as /baz",
		},
		{
			uses:     "",
			expected: "",
			message:  "expectedPath empty string to resolve as empty string",
		},
		{
			uses:     "./foo/bar/action.yaml",
			expected: "./foo/bar/action.yaml",
			message:  "expectedPath ./foo/bar/action.yaml to resolve as ./foo/bar/action.yaml",
		},
	}

	for _, test := range tests {
		step := Step{
			Uses: &test.uses,
		}

		assert.Equal(t, test.expected, step.GetRepoPath(), test.message)
	}

	// test that nil uses returns empty string
	step := Step{}
	assert.Equal(t, "", step.GetRepoPath())
}

func TestStep_ParseUses(t *testing.T) {
	type testData struct {
		uses         string
		expectedPath string
		expectedRepo string
		expectedRef  string
	}

	tests := []testData{
		{
			uses:         "actions/checkout@v4",
			expectedPath: "",
			expectedRepo: "actions/checkout",
			expectedRef:  "v4",
		},
		{
			uses:         "foo/bar/baz@main",
			expectedPath: "baz",
			expectedRepo: "foo/bar",
			expectedRef:  "main",
		},
		{
			uses:         "./foo/bar/action.yaml",
			expectedPath: "./foo/bar/action.yaml",
			expectedRepo: "",
			expectedRef:  "",
		},
	}

	for _, test := range tests {
		step := Step{
			Uses: &test.uses,
		}
		repo, path, ref, err := step.ParseUses()
		require.NoError(t, err)
		assert.Equal(t, test.expectedRepo, repo)
		assert.Equal(t, test.expectedPath, path)
		assert.Equal(t, test.expectedRef, ref)
	}

	// test that nil uses returns empty string
	step := Step{}
	_, _, _, err := step.ParseUses()
	assert.Error(t, err)
}
