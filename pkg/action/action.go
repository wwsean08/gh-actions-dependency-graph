package action

import (
	"errors"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/wwsean08/actions-dependency-graph/pkg/common"
	"sigs.k8s.io/yaml"
)

// Action a simplistic representation of a GitHub Action
type Action struct {
	Name        *string    `yaml:"name"`
	Description *string    `yaml:"description"`
	Runs        *RunsBlock `yaml:"runs"`
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

	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	action := new(Action)
	err = yaml.Unmarshal(data, action)
	return action, err
}

// GetActionFromRepoPath downloads the action.yaml file from a repository at a specific
// git ref and parses it into an Action struct
func GetActionFromRepoPath(owner, repo, ref, path string) (*Action, error) {
	return nil, errors.New("TODO: Implement")
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
