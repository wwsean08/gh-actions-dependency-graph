package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "actions-dependency-graph",
	Annotations: map[string]string{
		cobra.CommandDisplayNameAnnotation: "gh actions-dependency-graph",
	},
	Short: "An application to determine the dependencies of your actions/workflows",
	Long: `An application that recursively walks actions and workflows to determine 
the list of dependent actions so you can feel confident in what you're shipping and 
allow you to identify problematic actions that you may want to avoid.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {}
