package scanner

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wwsean08/actions-dependency-graph/pkg/action"
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

func TestScanner_CheckActionVersionIsLatestWithBlankAction(t *testing.T) {
	act := new(action.Action)
	scanner := Scanner{}
	isLatest, err := scanner.CheckActionVersionIsLatest(act)
	require.EqualErrorf(t, err, LocalActionError.Error(), "")
	require.False(t, isLatest)
}

func TestScanner_CheckActionVersionIsLatestWithHashReturnsError(t *testing.T) {
	act := new(action.Action)
	act.Ref = "ea2ef62537500fb1ef526076b41a088420e6835a"
	scanner := Scanner{}
	isLatest, err := scanner.CheckActionVersionIsLatest(act)
	require.EqualError(t, err, LikelyShaError.Error())
	require.False(t, isLatest)
}

func TestGetShaForLatestReleaseReturnsErrOnNon200(t *testing.T) {
	defer gock.Off()
	gock.New("https://api.github.com").
		Get("/repos/actions/checkout/releases/latest").
		Reply(http.StatusNotFound)

	act := &action.Action{
		Repo: "actions/checkout",
	}
	client, err := api.DefaultHTTPClient()
	require.NoError(t, err)

	sha, err := getShaForLatestRelease(act, client)
	assert.Error(t, err)
	assert.Zero(t, sha)
}

func TestGetShaForLatestReleaseReturnsErrOnNon200ForGetSha(t *testing.T) {
	defer gock.Off()
	latestFile, err := os.Open("testdata/actions-checkout-latest.json")
	require.NoError(t, err)
	require.NotNil(t, latestFile)
	gock.New("https://api.github.com").
		Get("/repos/actions/checkout/releases/latest").
		Reply(http.StatusOK).
		Body(latestFile)
	gock.New("https://api.github.com").
		Get("/repos/actions.checkout/commits/v4.1.1").
		Reply(http.StatusNotFound)

	act := &action.Action{
		Repo: "actions/checkout",
	}
	client, err := api.DefaultHTTPClient()
	require.NoError(t, err)

	sha, err := getShaForLatestRelease(act, client)
	assert.Error(t, err)
	assert.Zero(t, sha)
}

func TestGetShaForLatestReleaseReturnsSha(t *testing.T) {
	defer gock.Off()
	latestFile, err := os.Open("testdata/actions-checkout-latest.json")
	require.NoError(t, err)
	require.NotNil(t, latestFile)
	commitFile, err := os.Open("testdata/actions-checkout-v4.1.1.json")
	require.NoError(t, err)
	require.NotNil(t, commitFile)

	gock.New("https://api.github.com").
		Get("/repos/actions/checkout/releases/latest").
		Reply(http.StatusOK).
		Body(latestFile)
	gock.New("https://api.github.com").
		Get("/repos/actions.checkout/commits/v4.1.1").
		Reply(http.StatusOK).
		Body(commitFile)

	act := &action.Action{
		Repo: "actions/checkout",
	}
	client, err := api.DefaultHTTPClient()
	require.NoError(t, err)

	sha, err := getShaForLatestRelease(act, client)
	assert.NoError(t, err)
	assert.Equal(t, "b4ffde65f46336ab88eb53be808477a3936bae11", sha)
}

func TestGetShaForRefReturnsRef(t *testing.T) {
	defer gock.Off()
	commitFile, err := os.Open("testdata/actions-checkout-v4.1.1.json")
	require.NoError(t, err)
	require.NotNil(t, commitFile)

	gock.New("https://api.github.com").
		Get("/repos/actions/checkout/commits/v4.1.1").
		Reply(http.StatusOK).
		Body(commitFile)

	client, err := api.DefaultHTTPClient()
	require.NoError(t, err)

	sha, err := getShaForRef("actions/checkout", "v4.1.1", client)
	assert.NoError(t, err)
	assert.Equal(t, "b4ffde65f46336ab88eb53be808477a3936bae11", sha)
}

func TestScanner_CheckActionVersionIsLatestReturnsTrue(t *testing.T) {
	defer gock.Off()
	latestFile, err := os.Open("testdata/actions-checkout-latest.json")
	require.NoError(t, err)
	require.NotNil(t, latestFile)
	commitFile, err := os.Open("testdata/actions-checkout-v4.1.1.json")
	require.NoError(t, err)
	require.NotNil(t, commitFile)

	gock.New("https://api.github.com").
		Get("/repos/actions/checkout/releases/latest").
		Reply(http.StatusOK).
		Body(latestFile)
	gock.New("https://api.github.com").
		Get("/repos/actions/checkout/commits/v4.1.1").
		Times(2).
		Reply(http.StatusOK).
		Body(commitFile)

	act := &action.Action{
		Repo: "actions/checkout",
		Ref:  "v4.1.1",
	}
	scanner := &Scanner{}

	isLatest, err := scanner.CheckActionVersionIsLatest(act)
	assert.NoError(t, err)
	assert.True(t, isLatest)
}

func TestScanner_CheckActionVersionIsLatestReturnsFalse(t *testing.T) {
	defer gock.Off()
	latestFile, err := os.Open("testdata/actions-checkout-latest.json")
	require.NoError(t, err)
	require.NotNil(t, latestFile)
	commitFile411, err := os.Open("testdata/actions-checkout-v4.1.1.json")
	require.NoError(t, err)
	require.NotNil(t, commitFile411)
	commitFile410, err := os.Open("testdata/actions-checkout-v4.1.0.json")
	require.NoError(t, err)
	require.NotNil(t, commitFile410)

	gock.New("https://api.github.com").
		Get("/repos/actions/checkout/releases/latest").
		Reply(http.StatusOK).
		Body(latestFile)
	gock.New("https://api.github.com").
		Get("/repos/actions/checkout/commits/v4.1.1").
		Reply(http.StatusOK).
		Body(commitFile411)
	gock.New("https://api.github.com").
		Get("/repos/actions/checkout/commits/v4.1.0").
		Reply(http.StatusOK).
		Body(commitFile410)

	act := &action.Action{
		Repo: "actions/checkout",
		Ref:  "v4.1.0",
	}
	scanner := &Scanner{}

	isLatest, err := scanner.CheckActionVersionIsLatest(act)
	assert.NoError(t, err)
	assert.False(t, isLatest)
}

func TestScanner_CheckActionVersionIsLatestReturnsErrOnLatest404(t *testing.T) {
	defer gock.Off()
	latestFile, err := os.Open("testdata/actions-checkout-latest.json")
	require.NoError(t, err)
	require.NotNil(t, latestFile)
	commitFile411, err := os.Open("testdata/actions-checkout-v4.1.1.json")
	require.NoError(t, err)
	require.NotNil(t, commitFile411)
	commitFile410, err := os.Open("testdata/actions-checkout-v4.1.0.json")
	require.NoError(t, err)
	require.NotNil(t, commitFile410)

	gock.New("https://api.github.com").
		Get("/repos/actions/checkout/releases/latest").
		Reply(http.StatusNotFound)

	act := &action.Action{
		Repo: "actions/checkout",
		Ref:  "v4.1.0",
	}
	scanner := &Scanner{}

	isLatest, err := scanner.CheckActionVersionIsLatest(act)
	assert.Error(t, err)
	assert.False(t, isLatest)
}

func TestScanner_CheckActionVersionIsLatestReturnsErrorGettingVersionForAction(t *testing.T) {
	defer gock.Off()
	latestFile, err := os.Open("testdata/actions-checkout-latest.json")
	require.NoError(t, err)
	require.NotNil(t, latestFile)
	commitFile410, err := os.Open("testdata/actions-checkout-v4.1.0.json")
	require.NoError(t, err)
	require.NotNil(t, commitFile410)

	gock.New("https://api.github.com").
		Get("/repos/actions/checkout/releases/latest").
		Reply(http.StatusOK).
		Body(latestFile)
	gock.New("https://api.github.com").
		Get("/repos/actions/checkout/commits/v4.1.1").
		Reply(http.StatusNotFound)
	gock.New("https://api.github.com").
		Get("/repos/actions/checkout/commits/v4.1.0").
		Reply(http.StatusOK).
		Body(commitFile410)

	act := &action.Action{
		Repo: "actions/checkout",
		Ref:  "v4.1.0",
	}
	scanner := &Scanner{}

	isLatest, err := scanner.CheckActionVersionIsLatest(act)
	assert.Error(t, err)
	assert.False(t, isLatest)
}
