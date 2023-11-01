package action

import (
	"fmt"
	"strings"

	"github.com/goccy/go-graphviz/cgraph"
)

func (a *Action) DFSPrint(tab int) {
	fmt.Println(fmt.Sprintf("%s%s", strings.Repeat("\t", tab), a.Repo))
	if len(a.DependentActions) != 0 {
		for _, subAction := range a.DependentActions {
			subAction.DFSPrint(tab + 1)
		}
	}
}

// GenerateGraph is used to recursively generate a graph
func (a *Action) GenerateGraph(parentGraph, jobGraph *cgraph.Graph) error {
	return a.generateGraph(nil, parentGraph, jobGraph)
}

func (a *Action) generateGraph(parent *cgraph.Node, parentGraph, jobGraph *cgraph.Graph) error {
	if a.Runs == nil || a.Runs.Using == "" {
		return nil
	}
	act, err := jobGraph.CreateNode(fmt.Sprintf("%s/%s@%s", a.Repo, a.Path, a.Ref))
	if err != nil {
		return err
	}

	label := fmt.Sprintf("%s@%s", a.Repo, a.Ref)

	if a.Path != "" {
		label = fmt.Sprintf("%s/%s@%s", a.Repo, a.Path, a.Ref)
	}

	if a.IsJavascript() {
		act.SetColor("green")
		act.SetLabel(fmt.Sprintf("%s\n%s", label, a.Runs.Using))
	} else if a.IsDocker() {
		act.SetColor("blue")
		if a.Runs.Image == nil {
			act.SetLabel(label)
		} else {
			dockerfile := "local Dockerfile"
			if *a.Runs.Image != "Dockerfile" {
				dockerfile = *a.Runs.Image
			}
			act.SetLabel(fmt.Sprintf("%s\n%s", label, dockerfile))
		}
	}

	if parent != nil {
		_, err = jobGraph.CreateEdge("", parent, act)
		if err != nil {
			return err
		}
	}
	if len(a.DependentActions) != 0 {
		for _, subAction := range a.DependentActions {
			return subAction.generateGraph(act, parentGraph, jobGraph)
		}
	}

	return nil
}
