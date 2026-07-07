package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "repolift",
	Short: "A declarative workspace manager for your git repositories",
	Long: `Repolift allows you to declaratively define your local folder structure
and the git repositories that should be cloned into them.`,
	// If no subcommand is given, default to the 'sync' command.
	RunE: func(cmd *cobra.Command, args []string) error {
		// We want to run the sync command by default
		syncCmd.SetArgs(args)
		return syncCmd.Execute()
	},
	// Add aliases for the main command
	Aliases: []string{"rlift", "rl"},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// We no longer need the init() function that added the applyCmd here,
// because syncCmd is now added in its own file (cmd/sync.go).
