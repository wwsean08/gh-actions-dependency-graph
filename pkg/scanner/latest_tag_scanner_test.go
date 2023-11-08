package scanner

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLikelySha1Sum(t *testing.T) {
	type testCase struct {
		text     string
		expected bool
	}
	testCases := []testCase{
		{
			text:     "foo",
			expected: false,
		},
		{
			text:     "abc123",
			expected: false,
		},
		{
			text:     "Lorem ipsum dolor sit amet, consectetur ",
			expected: false,
		},
		{
			text:     "ea2ef62537500fb1ef526076b41a088420e6835a",
			expected: true,
		},
	}

	for _, test := range testCases {
		assert.Equal(t, test.expected, likelySha1Sum(test.text), fmt.Sprintf("%s should have returned %t", test.text, test.expected))
	}
}
