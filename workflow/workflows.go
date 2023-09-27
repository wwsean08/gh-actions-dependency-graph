package workflow

import (
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

type Workflow struct {
	Jobs map[string]Job `yaml:"jobs"`
}

type Job struct {
	Steps []Step `yaml:"steps"`
}

type Step struct {
	Uses *string `yaml:"uses"`
	Name *string `yaml:"name"`
	Id   *string `yaml:"id"`
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
