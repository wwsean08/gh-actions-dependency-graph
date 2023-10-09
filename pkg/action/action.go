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

type Action struct {
	Name        *string    `yaml:"name"`
	Description *string    `yaml:"description"`
	Runs        *RunsBlock `yaml:"runs"`
}

type RunsBlock struct {
	Using string         `yaml:"using"`
	Steps []*common.Step `yaml:"steps"`
}

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

func (a *Action) IsComposite() bool {
	return a.Runs.Using == "composite"
}

func (a *Action) IsJavascript() bool {
	return strings.HasPrefix(a.Runs.Using, "node")
}

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

func (a *Action) IsDocker() bool {
	return a.Runs.Using == "docker"
}
