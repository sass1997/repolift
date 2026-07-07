package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sass1997/repolift/config"
	"github.com/sass1997/repolift/internal/importer"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var configFromWorkspaceCmd = &cobra.Command{
	Use:   "config-from-workspace [path]",
	Short: "Generates a repolift configuration from an existing workspace directory",
	Long: `Scans a directory for Git repositories and generates a corresponding repolift 
workspace configuration. It will then prompt you to either append it to your 
global config, save it to a new file, or exit.`,
	Args: cobra.ExactArgs(1), // Ensures exactly one argument (the path) is provided.
	Run: func(cmd *cobra.Command, args []string) {
		workspacePath := args[0]
		fmt.Printf("🔍 Scanning directory '%s' to generate configuration...\n", workspacePath)

		// 1. Discover repositories
		newWorkspace, err := importer.Discover(workspacePath)
		if err != nil {
			fmt.Printf("❌ Error during discovery: %v\n", err)
			os.Exit(1)
		}

		if len(newWorkspace.Repositories) == 0 {
			fmt.Println("✅ No Git repositories found in the specified directory.")
			return
		}

		// 2. Present the generated config
		fmt.Println("\n📄 Generated Workspace Configuration:")
		fmt.Println("------------------------------------")
		yamlBytes, _ := yaml.Marshal([]config.Workspace{*newWorkspace})
		fmt.Println(string(yamlBytes))
		fmt.Println("------------------------------------")

		// 3. Start interactive prompt
		handleGeneratedConfig(newWorkspace)
	},
}

func handleGeneratedConfig(ws *config.Workspace) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("What would you like to do? [ (a)ppend to global config | (s)ave to new file | (e)xit ]: ")
		response, _ := reader.ReadString('\n')
		response = strings.ToLower(strings.TrimSpace(response))

		switch response {
		case "a", "append":
			appendToGlobalConfig(ws)
			return
		case "s", "save":
			saveToNewFile(ws)
			return
		case "e", "exit":
			fmt.Println("Aborted.")
			return
		default:
			fmt.Println("Invalid option. Please choose 'a', 's', or 'e'.")
		}
	}
}

// appendToGlobalConfig now correctly uses the XDG standard.
func appendToGlobalConfig(ws *config.Workspace) {
	// 1. Get the XDG-compliant default path.
	configPath, err := config.GetDefaultConfigPath()
	if err != nil {
		fmt.Printf("❌ Error determining default config path: %v\n", err)
		return
	}

	// 2. Ensure the directory exists before trying to read or write.
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		fmt.Printf("❌ Error creating config directory '%s': %v\n", configDir, err)
		return
	}

	// 3. Load existing config (it's okay if it doesn't exist).
	cfg, err := config.Load(configPath)
	if err != nil && !os.IsNotExist(err) {
		fmt.Printf("❌ Error loading existing config '%s': %v\n", configPath, err)
		return
	}
	if cfg == nil {
		cfg = &config.Config{} // Create a new config if one doesn't exist
	}

	// 4. Append the new workspace and write back.
	cfg.Workspaces = append(cfg.Workspaces, *ws)
	yamlBytes, _ := yaml.Marshal(cfg)
	err = os.WriteFile(configPath, yamlBytes, 0644)
	if err != nil {
		fmt.Printf("❌ Error writing to global config '%s': %v\n", configPath, err)
		return
	}
	fmt.Printf("✅ Successfully appended to '%s'\n", configPath)
}

func saveToNewFile(ws *config.Workspace) {
	fileName := "repolift-generated.yaml"
	cfg := &config.Config{Workspaces: []config.Workspace{*ws}}
	yamlBytes, _ := yaml.Marshal(cfg)
	err := os.WriteFile(fileName, yamlBytes, 0644)
	if err != nil {
		fmt.Printf("❌ Error writing to new file '%s': %v\n", fileName, err)
		return
	}
	fmt.Printf("✅ Successfully saved to '%s'\n", fileName)
}

func init() {
	rootCmd.AddCommand(configFromWorkspaceCmd)
}
