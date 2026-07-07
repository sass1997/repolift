package importer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sass1997/repolift/config"
)

// Discover scans a root path and generates a Workspace configuration
// by finding all git repositories within it.
func Discover(rootPath string) (*config.Workspace, error) {
	absRootPath, err := filepath.Abs(rootPath)
	if err != nil {
		return nil, fmt.Errorf("could not get absolute path for '%s': %w", rootPath, err)
	}

	workspace := &config.Workspace{
		Path:         absRootPath, // We use the absolute path for the new config
		Repositories: []config.Repository{},
	}

	err = filepath.WalkDir(absRootPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// If we find a .git directory, we've found a repository.
		if d.IsDir() && d.Name() == ".git" {
			repoDir := filepath.Dir(path)

			// Get the remote URL from the git config
			remoteURL, gitErr := getGitRemoteURL(repoDir)
			if gitErr != nil {
				fmt.Printf("⚠️  Skipping '%s': could not get remote URL (%v)\n", repoDir, gitErr)
				return filepath.SkipDir // Skip this directory
			}

			// The 'dir' in the config is the relative path from the workspace root
			relativeDir, _ := filepath.Rel(absRootPath, repoDir)

			repo := config.Repository{
				URL: remoteURL,
				Dir: relativeDir,
			}
			workspace.Repositories = append(workspace.Repositories, repo)

			// We don't need to look for git repos inside another git repo.
			return filepath.SkipDir
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking directory '%s': %w", rootPath, err)
	}

	return workspace, nil
}

// getGitRemoteURL executes 'git config' to find the URL of the 'origin' remote.
func getGitRemoteURL(repoPath string) (string, error) {
	cmd := exec.Command("git", "-C", repoPath, "config", "--get", "remote.origin.url")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
