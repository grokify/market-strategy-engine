package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mse",
	Short: "Market Strategy Engine CLI",
	Long: `Market Strategy Engine (mse) is a strategic capability gap analysis system.

It helps you:
  - Map competitors across product lines and market segments
  - Score capabilities with segment-specific weightings
  - Compute weighted gaps to identify strategic priorities
  - Generate prioritized roadmaps for market expansion

Use "mse [command] --help" for more information about a command.`,
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags can be added here
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
}
