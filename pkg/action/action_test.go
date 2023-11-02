package action

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseActionReturnsErrOnNonexistentFile(t *testing.T) {
	action, err := ParseAction("testdata/non-existent-file.404")
	require.Error(t, err)
	require.Nil(t, action)
}

func TestPrivateParseActionReturnsErrOnNil(t *testing.T) {
	action, err := parseAction(nil)
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

func TestPrivateGetActionFromRepoPathReturnsErrOn404(t *testing.T) {
	defer gock.Off()
	gock.New("https://api.github.com").
		Get("/repos/wwsean08/actions-dependency-graph/contents/action.yaml").
		Reply(http.StatusNotFound)

	dlUrl, err := getActionFromRepoPath("wwsean08/actions-dependency-graph", "main", "/action.yaml")
	assert.EqualError(t, err, "unable to find /action.yaml in github repo@ref wwsean08/actions-dependency-graph@main")
	assert.Nil(t, dlUrl)
}

func TestPrivateGetActionFromRepoPathReturnsErrOnHttpError(t *testing.T) {
	defer gock.Off()
	gock.New("https://api.github.com").
		Get("/repos/wwsean08/actions-dependency-graph/contents/action.yaml").
		ReplyError(fmt.Errorf("fake http error"))

	dlUrl, err := getActionFromRepoPath("wwsean08/actions-dependency-graph", "main", "/action.yaml")
	assert.ErrorContains(t, err, "fake http error")
	assert.Nil(t, dlUrl)
}

func TestGetActionFromRepoPathReturnsErrorOnNotFoundTwice(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("repos/actions/checkout/contents/action.yml").
		Reply(http.StatusNotFound)

	gock.New("https://api.github.com").
		Get("repos/actions/checkout/contents/action.yaml").
		Reply(http.StatusNotFound)

	action, err := GetActionFromRepoPath("actions/checkout", "", "")
	assert.ErrorContains(t, err, "unable to find action at")
	assert.Nil(t, action)
}

func TestPrivateGetActionFromRepoPathReturnsValidData(t *testing.T) {
	defer gock.Off()

	data, err := os.Open("testdata/sample-get-contents-response.json")
	require.NoError(t, err)
	require.NotNil(t, data)

	gock.New("https://api.github.com").
		Get("repos/actions/checkout/contents/action.yml").
		Reply(http.StatusOK).
		Body(data)

	dlUrl, err := getActionFromRepoPath("actions/checkout", "", "/action.yml")
	assert.NoError(t, err)
	assert.NotNil(t, dlUrl)
	assert.Equal(t, "https://raw.githubusercontent.com/actions/checkout/main/action.yml", *dlUrl)
}

func TestGetActionFromRepoPathReturnsValidDataForDotYml(t *testing.T) {
	defer gock.Off()

	getContentData, err := os.Open("testdata/sample-get-contents-response.json")
	require.NoError(t, err)
	require.NotNil(t, getContentData)

	actionData, err := os.Open("testdata/js-example.yaml")
	require.NoError(t, err)
	require.NotNil(t, actionData)

	gock.New("https://api.github.com").
		Get("repos/actions/checkout/contents/action.yml").
		Reply(http.StatusOK).
		Body(getContentData)

	gock.New("https://raw.githubusercontent.com").
		Get("/actions/checkout/main/action.yml").
		Reply(http.StatusOK).
		Body(actionData)

	action, err := GetActionFromRepoPath("actions/checkout", "", "")
	assert.NoError(t, err)
	assert.NotNil(t, action)
	assert.True(t, action.IsJavascript())
	assert.Equal(t, "Hello World", *action.Name)
}

func TestGetActionFromRepoPathReturnsValidDataForDotYaml(t *testing.T) {
	defer gock.Off()

	getContentData, err := os.Open("testdata/sample-get-contents-response.json")
	require.NoError(t, err)
	require.NotNil(t, getContentData)

	actionData, err := os.Open("testdata/js-example.yaml")
	require.NoError(t, err)
	require.NotNil(t, actionData)

	gock.New("https://api.github.com").
		Get("repos/actions/checkout/contents/action.yml").
		Reply(http.StatusNotFound)

	gock.New("https://api.github.com").
		Get("repos/actions/checkout/contents/action.yaml").
		Reply(http.StatusOK).
		Body(getContentData)

	gock.New("https://raw.githubusercontent.com").
		Get("/actions/checkout/main/action.yml").
		Reply(http.StatusOK).
		Body(actionData)

	action, err := GetActionFromRepoPath("actions/checkout", "", "")
	assert.NoError(t, err)
	assert.NotNil(t, action)
	assert.True(t, action.IsJavascript())
	assert.Equal(t, "Hello World", *action.Name)
}
