package scanner

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/wwsean08/actions-dependency-graph/pkg/action"
)

var likelySha1, _ = regexp.Compile("^[0-9A-Fa-f]{40}$")

// CheckActionVersionIsLatest checks if the action is referencing the latest  version
// of an action if the git ref is a tag.
//
// The comparison checks the commit sha associated with the two tags to allow for floating
// tags to be used without having to worry about exact matches of tags.
func (s *Scanner) CheckActionVersionIsLatest(action *action.Action) (bool, string, error) {
	if action.Ref == "" {
		return false, "", LocalActionError
	}
	if likelySha1Sum(action.Ref) {
		return false, "", LikelyShaError
	}

	client, err := api.DefaultHTTPClient()
	if err != nil {
		return false, "", err
	}

	shaNow, err := getShaForRef(action.Repo, action.Ref, client)
	if err != nil {
		return false, "", err
	}

	shaLatest, tag, err := getShaForLatestRelease(action, client)
	if err != nil {
		return false, "", err
	}

	return shaNow == shaLatest, tag, nil
}

// getShaForLatestRelease takes an action and client, and gets the commit sha for the latest
// release according to GitHub's release API
func getShaForLatestRelease(action *action.Action, client *http.Client) (string, string, error) {
	type data struct {
		TagName string `json:"tag_name"`
	}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", action.Repo), nil)
	if err != nil {
		return "", "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}

	if resp.StatusCode != 200 {
		return "", "", fmt.Errorf("http error occurred: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	tagData := new(data)
	err = json.Unmarshal(body, tagData)
	if err != nil {
		return "", "", err
	}
	sha, err := getShaForRef(action.Repo, tagData.TagName, client)

	return sha, tagData.TagName, err
}

// getShaForRef takes a repo, ref, and http client, and retrieves the sha commit sha for that reference
func getShaForRef(repo, ref string, client *http.Client) (string, error) {
	type data struct {
		Sha string `json:"sha"`
	}

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://api.github.com/repos/%s/commits/%s", repo, ref), nil)
	if err != nil {
		return "", err
	}
	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("http error occurred: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	respData := new(data)
	err = json.Unmarshal(body, respData)
	if err != nil {
		return "", err
	}

	return respData.Sha, nil
}

// likelySha1Sum checks if it's a 40 character hex string, if so it's likely a sha-1
// string, and we don't need to look up the version
func likelySha1Sum(text string) bool {
	return likelySha1.MatchString(text)
}
