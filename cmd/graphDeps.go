package cmd

import (
	"fmt"
	"os"

	"github.com/goccy/go-graphviz"
	"github.com/spf13/cobra"
	"github.com/wwsean08/actions-dependency-graph/pkg/action"
	"github.com/wwsean08/actions-dependency-graph/pkg/workflow"
)

// graphDepsCmd represents the graphDeps command
var graphDepsCmd = &cobra.Command{
	Use:   "graph-deps",
	Short: "A command to create a graphviz representation of dependencies",
	Long:  ``,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("one input required, the workflow to analyze")
		}
		_, err := os.Stat(args[0])
		return err
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		wf, err := workflow.ParseWorkflow(args[0])
		if err != nil {
			return err
		}
		graphContainer := graphviz.New()
		graph, err := graphContainer.Graph()
		if err != nil {
			return err
		}
		defer graph.Close()
		defer graphContainer.Close()
		graph.SetID(args[0])
		graph.SetLabel(args[0])

		graphIterator := 1
		for key, val := range wf.Jobs {
			jobGraph := graph.SubGraph(key, 1)
			jobGraph.SetLabel(key)
			graphIterator += 1
			for _, step := range val.Steps {
				repo, path, ref, err := step.ParseUses()
				if err != nil {
					// the only possible error is that the uses is blank which is valid,
					// hence the continue statement
					continue
				}
				act, err := action.GetActionFromRepoPath(repo, ref, path)
				if err != nil {
					return err
				}
				err = act.GetDependentActions()
				if err != nil {
					if _, ok := err.(action.NoDependenciesError); !ok {
						// if it is any error but no dependencies error; error out
						return err
					}
				}
				err = act.GenerateGraph(graph, jobGraph)
				if err != nil {
					return err
				}
			}
		}
		outputFile, err := cmd.Flags().GetString("output")
		if err != nil {
			return err
		}
		return graphContainer.RenderFilename(graph, graphviz.XDOT, outputFile)
	},
}

func init() {
	rootCmd.AddCommand(graphDepsCmd)

	graphDepsCmd.Flags().StringP("output", "o", "deps.graphviz", "Used to specify the output file for the GraphViz file")
}
