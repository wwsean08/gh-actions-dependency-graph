package cmd

import (
	"fmt"
	"github.com/wwsean08/actions-dependency-graph/pkg/action"
	"github.com/wwsean08/actions-dependency-graph/pkg/scanner"
	"github.com/wwsean08/actions-dependency-graph/pkg/workflow"
	"os"

	"github.com/spf13/cobra"
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Perform a basic security scan, looking for potential issues",
	Long: `This security scan looks for the following issues:
1. EOL versions of NodeJS like node 16.`,
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
		scan := scanner.NewDefaultScanner()
		for key, val := range wf.Jobs {
			fmt.Printf("Scanning for job %s\n", key)
			for _, step := range val.Steps {
				repo, path, ref, err := step.ParseUses()
				if err != nil {
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
				fmt.Printf("Scanning %s/%s@%s\n", act.Repo, act.Path, act.Ref)
				results, errs := scan.Scan(act)
				if len(errs) != 0 {
					fmt.Printf("Errors occured scanning action %s/%s@%s: %v", act.Repo, act.Path, act.Ref, errs)
				}
				fmt.Print(scan.FormatResults(results))
			}
			fmt.Println()
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
