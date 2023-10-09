package workflow

import (
	"io"
	"os"

	"github.com/wwsean08/actions-dependency-graph/pkg/common"
	"sigs.k8s.io/yaml"
)

type Workflow struct {
	Jobs map[string]Job `yaml:"jobs"`
}

type Job struct {
	Steps []*common.Step `yaml:"steps"`
}

func ParseWorkflow(file string) (*Workflow, error) {
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

	workflow := new(Workflow)
	err = yaml.Unmarshal(data, workflow)
	return workflow, err
}
