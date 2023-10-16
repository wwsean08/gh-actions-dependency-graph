package action

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/wwsean08/actions-dependency-graph/pkg/common"
	"sigs.k8s.io/yaml"
)

// Action a simplistic representation of a GitHub Action
type Action struct {
	Name             *string    `yaml:"name"`
	Description      *string    `yaml:"description"`
	Runs             *RunsBlock `yaml:"runs"`
	Repo             string
	DependentActions []*Action
}

// RunsBlock represents the runs section of a GitHub Action
type RunsBlock struct {
	Using string         `yaml:"using"`
	Steps []*common.Step `yaml:"steps"`
}

// ParseAction parses an action file and returns an Action.
func ParseAction(file string) (*Action, error) {
	if _, err := os.Stat(file); err != nil {
		return nil, err
	}
	r, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	return parseAction(r)
}

// parseAction takes in a reader and returns an action or error if one occurs while reading/parsing
func parseAction(reader io.Reader) (*Action, error) {
	if reader == nil {
		return nil, fmt.Errorf("reader cannot be nil")
	}
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	action := new(Action)
	err = yaml.Unmarshal(data, action)
	return action, err
}

// GetActionFromRepoPath downloads the action.yaml file from a repository at a specific
// git ref and parses it into an Action struct
func GetActionFromRepoPath(repo, ref, path string) (*Action, error) {
	dlUrl, err := getActionFromRepoPath(repo, ref, fmt.Sprintf("%s/action.yml", path))
	if err != nil {
		dlUrl, err = getActionFromRepoPath(repo, ref, fmt.Sprintf("%s/action.yaml", path))
		if err != nil {
			return nil, errors.Join(fmt.Errorf("unable to find action at either %s/action.yaml or %s/action.yml in repo %s at ref %s", path, path, repo, ref), err)
		}
	}

	resp, err := http.Get(*dlUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	action, err := parseAction(resp.Body)
	if err != nil {
		return nil, err
	}
	action.Repo = repo
	return action, nil
}

// getActionFromRepoPath takes care of the nitty-gritty API calls to get the action file URL
func getActionFromRepoPath(repo, ref, path string) (*string, error) {
	// https://docs.github.com/en/rest/repos/contents?apiVersion=2022-11-28#get-repository-content
	client, err := api.DefaultHTTPClient()
	if err != nil {
		return nil, err
	}

	type contentsData struct {
		DownloadURL string `json:"download_url"`
	}

	// Create request
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://api.github.com/repos/%s/contents%s", repo, path), nil)
	if err != nil {
		return nil, err
	}
	// Set headers
	req.Header.Set("Content-Type", "application/vnd.github.raw")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	// Set query param
	q := req.URL.Query()
	q.Add("ref", ref)
	req.URL.RawQuery = q.Encode()

	// Perform request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("unable to find %s in github repo@ref %s@%s", path, repo, ref)
	}

	respData := new(contentsData)
	respBody, err := io.ReadAll(resp.Body)

	err = json.Unmarshal(respBody, respData)
	if err != nil {
		return nil, err
	}

	return &respData.DownloadURL, nil
}

// IsComposite returns true if the Action is a composite action
func (a *Action) IsComposite() bool {
	return a.Runs.Using == "composite"
}

// IsJavascript returns true if the Action is a javascript action
func (a *Action) IsJavascript() bool {
	return strings.HasPrefix(a.Runs.Using, "node")
}

// GetNodeVersion returns the version of Node.js the Action runs on if it is a javascript
// Action
func (a *Action) GetNodeVersion() (int, error) {
	if !a.IsJavascript() {
		return 0, errors.New("action is not a javascript action")
	}

	textVersion, _ := strings.CutPrefix(a.Runs.Using, "node")
	intVersion, err := strconv.Atoi(textVersion)
	if err != nil {
		return 0, err
	}

	return intVersion, nil
}

// IsDocker returns true is the Action is a docker action
func (a *Action) IsDocker() bool {
	return a.Runs.Using == "docker"
}

// GetDependentActions recursively retrieves the dependent actions to generate a graph of dependencies
func (a *Action) GetDependentActions() error {
	if !a.IsComposite() {
		return NoDependenciesError{msg: "action is not a composite action, and so has no dependent actions"}
	}
	a.DependentActions = make([]*Action, 0)
	for _, step := range a.Runs.Steps {
		if step.Uses == nil {
			// this isn't calling another action, but defining code in line, skip
			continue
		}

		// at this point we know that there is an action call
		repo, path, ref, err := step.ParseUses()
		if err != nil {
			return err
		}

		action, err := GetActionFromRepoPath(repo, ref, path)
		if err != nil {
			return err
		}
		a.DependentActions = append(a.DependentActions, action)
		err = action.GetDependentActions()
		if err != nil {
			// Ignore No Dependencies Error as that is expected eventually
			if _, ok := err.(NoDependenciesError); !ok {
				return err
			}
		}
	}
	return nil
}
