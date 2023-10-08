package workflow

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
			message:  "expected actions/checkout@v4 to resolve as v4",
		},
		{
			uses:     "foo/bar/baz@main",
			expected: "main",
			message:  "expected foo/bar/baz@main to resolve as main",
		},
		{
			uses:     "",
			expected: "",
			message:  "expected empty string to resolve as empty string",
		},
		{
			uses:     "./foo/action.yaml",
			expected: "",
			message:  "expected local action to not resolve a commitish",
		},
	}

	for _, test := range tests {
		step := Step{
			Uses: &test.uses,
		}
		assert.Equal(t, test.expected, step.GetCommitish(), test.message)
	}

	// Test that nil will return empty string
	step := Step{}
	assert.Equal(t, "", step.GetCommitish())
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
			message:  "expected actions/checkout@v4 to resolve as actions/checkout",
		},
		{
			uses:     "foo/bar/baz@main",
			expected: "foo/bar",
			message:  "expected foo/bar/baz@main to resolve as foo/bar",
		},
		{
			uses:     "",
			expected: "",
			message:  "expected empty string to resolve as empty string",
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
			message:  "expected actions/checkout@v4 to resolve as empty string",
		},
		{
			uses:     "foo/bar/baz@main",
			expected: "baz",
			message:  "expected foo/bar/baz@main to resolve as /baz",
		},
		{
			uses:     "",
			expected: "",
			message:  "expected empty string to resolve as empty string",
		},
		{
			uses:     "./foo/bar/action.yaml",
			expected: "./foo/bar/action.yaml",
			message:  "expected ./foo/bar/action.yaml to resolve as ./foo/bar/action.yaml",
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
