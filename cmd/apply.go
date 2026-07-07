package cmd

import (
	"fmt"
	"os"

	"github.com/repolift/repolift/config"
	"github.com/repolift/repolift/internal/workspace"
	"github.com/spf13/cobra"
)

var configFile string

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply the workspace configuration",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load(configFile)
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Loaded configuration with %d workspaces.\n", len(cfg.Workspaces))
		
		manager := workspace.NewManager(cfg)
		if err := manager.Apply(); err != nil {
			fmt.Printf("Error applying workspace: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("🎉 Workspace applied successfully!")
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)
	applyCmd.Flags().StringVarP(&configFile, "file", "f", "repolift.yaml", "config file")
}
