package workspace

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/repolift/repolift/config"
)

type Manager struct {
	config *config.Config
}

func NewManager(cfg *config.Config) *Manager {
	return &Manager{config: cfg}
}

func (m *Manager) Apply() error {
	for _, ws := range m.config.Workspaces {
		expandedPath, err := expandPath(ws.Path)
		if err != nil {
			return fmt.Errorf("failed to expand path %s: %w", ws.Path, err)
		}

		fmt.Printf("Processing workspace: %s\n", expandedPath)

		// Create workspace directory if it doesn't exist
		if err := os.MkdirAll(expandedPath, 0755); err != nil {
			return fmt.Errorf("failed to create workspace directory %s: %w", expandedPath, err)
		}

		for _, repo := range ws.Repositories {
			repoPath := filepath.Join(expandedPath, repo.Dir)

			if err := m.ensureRepository(repo.URL, repoPath); err != nil {
				fmt.Printf("❌ Failed to setup repository %s: %v\n", repo.Dir, err)
				continue
			}
		}
	}

	return nil
}

func (m *Manager) ensureRepository(url, path string) error {
	// Check if directory already exists
	if _, err := os.Stat(path); err == nil {
		// Directory exists, check if it's a git repository
		if _, err := os.Stat(filepath.Join(path, ".git")); err == nil {
			fmt.Printf("✅ Repository already exists: %s\n", path)
			// Later in v3 we can add logic to git pull here
			return nil
		}
		return fmt.Errorf("directory exists but is not a git repository")
	}

	fmt.Printf("⏳ Cloning %s into %s...\n", url, path)

	cmd := exec.Command("git", "clone", url, path)
	// Bind stdout and stderr so the user can see git's output
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git clone failed: %w", err)
	}

	fmt.Printf("✅ Successfully cloned %s\n", path)
	return nil
}

// expandPath expands the tilde (~) to the user's home directory
func expandPath(path string) (string, error) {
	if !strings.HasPrefix(path, "~") {
		return filepath.Abs(path)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot get user home directory: %w", err)
	}

	// Remove the tilde and prepend the home directory
	return filepath.Join(homeDir, strings.TrimPrefix(path, "~")), nil
}
