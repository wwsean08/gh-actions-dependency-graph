package scanner

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/pkg/errors"
)

// RepoJackingStatus contains the information for if a repo is susceptible to repo-jacking and useful information
type RepoJackingStatus struct {
	original    string
	new         string
	err         error
	susceptible bool
}

type ghRepo struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Owner    struct {
		Login string `json:"login"`
	} `json:"owner"`
}

// RepoJackingCheck checks if a repo is susceptible to repo-jacking, https://blog.aquasec.com/github-dataset-research-reveals-millions-potentially-vulnerable-to-repojacking
func (s *Scanner) RepoJackingCheck(repo string) *RepoJackingStatus {
	client, err := api.DefaultHTTPClient()
	if err != nil {
		return &RepoJackingStatus{
			original: repo,
			err:      err,
		}
	}
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://api.github.com/repos/%s", repo), nil)
	if err != nil {
		return &RepoJackingStatus{
			original: repo,
			err:      errors.Wrap(err, "error creating http request"),
		}
	}

	response, err := client.Do(request)
	if err != nil {
		return &RepoJackingStatus{
			original: repo,
			err:      errors.Wrap(err, "error requesting github repo via api"),
		}
	}

	if response.StatusCode != 200 {
		return &RepoJackingStatus{
			original: repo,
			err:      errors.Wrap(UnexpectedResponseCodeError, fmt.Sprintf("got %d instead of 200 when looking up repo %s", response.StatusCode, repo)),
		}
	}

	ghData := new(ghRepo)
	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		return &RepoJackingStatus{
			original: repo,
			err:      errors.Wrap(err, "error reading response body"),
		}
	}
	err = json.Unmarshal(respBody, ghData)
	if err != nil {
		return &RepoJackingStatus{
			original: repo,
			err:      errors.Wrap(err, "error unmarshalling json body"),
		}
	}

	repoJackingPossible := s.isRepoSusceptibleToRepoJacking(repo, ghData)

	return &RepoJackingStatus{
		original:    repo,
		new:         ghData.FullName,
		susceptible: repoJackingPossible,
		err:         nil,
	}
}

// isRepoSusceptibleToRepoJacking compares the repo to the result and checks if it's susceptible to repo-jacking
func (s *Scanner) isRepoSusceptibleToRepoJacking(repo string, ghData *ghRepo) bool {
	repoOwner := strings.Split(repo, "/")[0]
	return strings.EqualFold(repoOwner, ghData.Owner.Login)
}
