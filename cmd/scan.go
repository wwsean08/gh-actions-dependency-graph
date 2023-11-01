package cmd

import (
	"encoding/json"
	"errors"
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
		if err != nil {
			return err
		}

		outputMode, err := cmd.Flags().GetString("output")
		if err != nil {
			return err
		}

		switch outputMode {
		case "text":
			break
		case "json":
			break
		default:
			return errors.New("output flag must be either text or json")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		wf, err := workflow.ParseWorkflow(args[0])
		if err != nil {
			return err
		}
		scan := scanner.NewDefaultScanner()
		allResults := map[string][]scanner.Results{}
		outputMode, err := cmd.Flags().GetString("output")
		if err != nil {
			return err
		}

		for key, val := range wf.Jobs {
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

				if outputMode == "text" {
					fmt.Printf("Scanning %s/%s@%s\n", act.Repo, act.Path, act.Ref)
				}
				results, errs := scan.Scan(act)
				if len(errs) != 0 {
					fmt.Printf("Errors occured scanning action %s/%s@%s: %v", act.Repo, act.Path, act.Ref, errs)
				}
				allResults[key] = append(allResults[key], *results)
			}
		}

		switch outputMode {
		case "json":
			resultBytes, err := json.Marshal(allResults)
			if err != nil {
				return err
			}
			fmt.Printf("%s", resultBytes)
		case "text":
			for key, val := range allResults {
				fmt.Printf("Results for %s job\n", key)
				for _, result := range val {
					fmt.Printf("Results for %s\n", result.Action)
					fmt.Println(scan.FormatResults(&result))
				}
				fmt.Println()
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)

	scanCmd.Flags().StringP("output", "o", "text", "Used to set the output to text or json.  Acceptable values: [\"text\", \"json\"].")
}
