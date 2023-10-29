package cmd

import (
	"fmt"
	"github.com/wwsean08/actions-dependency-graph/pkg/action"
	"github.com/wwsean08/actions-dependency-graph/pkg/workflow"
	"os"

	"github.com/spf13/cobra"
)

// listDepsCmd represents the listDeps command
var listDepsCmd = &cobra.Command{
	Use:   "list-deps",
	Short: "List the dependent actions for the workflow jobs",
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
		for key, val := range wf.Jobs {
			fmt.Printf("Dependencies for %s:\n", key)
			for _, step := range val.Steps {
				if step.Uses != nil {
					fmt.Printf("%s\n", *step.Uses)
				} else if step.Name != nil {
					fmt.Printf("%s\n", *step.Name)
				} else if step.Id != nil {
					fmt.Printf("%s\n", *step.Id)
				}
				repo, path, ref, err := step.ParseUses()
				if err != nil {
					fmt.Printf("Inline code step\n")
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
					fmt.Printf("no dependencies for this action\n")
					continue
				}
				act.DFSPrint(1)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listDepsCmd)
}
