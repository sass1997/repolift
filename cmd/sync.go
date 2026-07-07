package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/sass1997/repolift/config"
	"github.com/sass1997/repolift/internal/executor"
	"github.com/sass1997/repolift/internal/planner"
	"github.com/spf13/cobra"
)

// Re-declaring the package-level variables for the flags.
var (
	configFile  string
	autoApprove bool
)

// syncCmd logic remains the same...
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Syncs the workspaces by planning and applying the desired state",
	Long: `Compares the desired state from the configuration file with the actual state 
of the filesystem. It then generates a plan of actions (clone, create dir, etc.) 
and prompts for confirmation before executing.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Determine config path
		finalConfigPath := determineConfigPath()

		// 2. Load config
		cfg, err := config.Load(finalConfigPath)
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("🔄 Loaded configuration from %s\n", finalConfigPath)

		// 3. Generate Plan
		p := planner.New(cfg)
		plan, err := p.GeneratePlan()
		if err != nil {
			fmt.Printf("Error generating plan: %v\n", err)
			os.Exit(1)
		}

		if len(plan.Actions) == 0 {
			fmt.Println("✅ Everything is already in sync. Nothing to do.")
			return
		}

		// 4. Display Plan and ask for confirmation
		fmt.Println("\nRepolift has generated the following execution plan:")
		fmt.Println("----------------------------------------------------")
		for _, action := range plan.Actions {
			fmt.Printf("[%s] %s\n", action.Type, action.Description)
		}
		fmt.Println("----------------------------------------------------")

		if !autoApprove {
			if !askForConfirmation("Do you want to perform these actions?") {
				fmt.Println("🚫 Plan execution aborted by user.")
				return
			}
		}

		// 5. Execute Plan
		fmt.Println("\n🚀 Applying plan...")
		if err := executor.ExecutePlan(plan); err != nil {
			fmt.Printf("\n❌ Error applying plan: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("\n🎉 Workspace synced successfully!")
	},
}

// determineConfigPath now correctly uses the XDG standard.
func determineConfigPath() string {
	// 1. User-provided flag has highest priority.
	if configFile != "" {
		return configFile
	}

	// 2. Get the XDG-compliant default path.
	defaultPath, err := config.GetDefaultConfigPath()
	if err != nil {
		fmt.Printf("Error: Could not determine default config path: %v\n", err)
		os.Exit(1)
	}

	// 3. Check if the default config exists.
	if _, err := os.Stat(defaultPath); err == nil {
		return defaultPath
	}

	// 4. As a fallback, check for a local repolift.yaml for convenience.
	localPath := "repolift.yaml"
	if _, err := os.Stat(localPath); err == nil {
		return localPath
	}

	// 5. If no configuration is found, guide the user.
	fmt.Printf("Error: No configuration file found.\n")
	fmt.Printf("Checked for default file at: %s\n", defaultPath)
	fmt.Printf("Checked for local fallback at: %s\n", localPath)
	fmt.Println("Please create a config file or specify one with the -f flag.")
	os.Exit(1)
	return "" // Unreachable
}

// askForConfirmation logic remains the same...
func askForConfirmation(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s [y/n]: ", prompt)
		response, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			return false
		}
		response = strings.ToLower(strings.TrimSpace(response))
		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}

func init() {
	rootCmd.AddCommand(syncCmd)
	syncCmd.Flags().StringVarP(&configFile, "file", "f", "", "config file (e.g., repolift.yaml)")
	syncCmd.Flags().BoolVar(&autoApprove, "auto-approve", false, "skip interactive approval before applying")
}
