package planner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rkathriner/repolift/config"
)

// ActionType defines the kind of operation to be performed.
type ActionType string

const (
	ActionCreateWorkspaceDir ActionType = "CREATE_DIR"
	ActionCloneRepository    ActionType = "CLONE"
	ActionNoOp               ActionType = "NO_OP" // Already exists and is correct
)

// Action represents a single operation in the execution plan.
type Action struct {
	Type        ActionType
	Path        string
	URL         string // Only for CLONE actions
	Description string
}

// Plan is a collection of actions to be executed.
type Plan struct {
	Actions []Action
}

// Planner generates an execution plan based on the desired state.
type Planner struct {
	config *config.Config
}

func New(cfg *config.Config) *Planner {
	return &Planner{config: cfg}
}

// GeneratePlan creates a plan of actions to match the desired state.
func (p *Planner) GeneratePlan() (*Plan, error) {
	plan := &Plan{}

	for _, ws := range p.config.Workspaces {
		expandedPath, err := p.expandPath(ws.Path)
		if err != nil {
			return nil, fmt.Errorf("failed to expand path %s: %w", ws.Path, err)
		}

		// 1. Check if workspace directory exists
		if _, err := os.Stat(expandedPath); os.IsNotExist(err) {
			plan.Actions = append(plan.Actions, Action{
				Type:        ActionCreateWorkspaceDir,
				Path:        expandedPath,
				Description: fmt.Sprintf("Create workspace directory '%s'", expandedPath),
			})
		}

		for _, repo := range ws.Repositories {
			repoPath := filepath.Join(expandedPath, repo.Dir)

			// 2. Check if repository directory exists
			if _, err := os.Stat(repoPath); os.IsNotExist(err) {
				plan.Actions = append(plan.Actions, Action{
					Type:        ActionCloneRepository,
					Path:        repoPath,
					URL:         repo.URL,
					Description: fmt.Sprintf("Clone repository '%s' into '%s'", repo.URL, repoPath),
				})
			} else if err == nil {
				if _, gitErr := os.Stat(filepath.Join(repoPath, ".git")); gitErr == nil {
					plan.Actions = append(plan.Actions, Action{
						Type:        ActionNoOp,
						Path:        repoPath,
						Description: fmt.Sprintf("Repository '%s' already exists", repoPath),
					})
				} else {
					return nil, fmt.Errorf("conflict: path '%s' exists but is not a git repository", repoPath)
				}
			} else {
				return nil, fmt.Errorf("could not stat path '%s': %w", repoPath, err)
			}
		}
	}

	return plan, nil
}

// expandPath handles path expansion with respect to the default_start_path.
func (p *Planner) expandPath(path string) (string, error) {
	// Highest priority: Tilde expansion
	if strings.HasPrefix(path, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("cannot get user home directory: %w", err)
		}
		return filepath.Join(homeDir, strings.TrimPrefix(path, "~")), nil
	}

	// Second priority: Absolute paths
	if filepath.IsAbs(path) {
		return path, nil
	}

	// Third priority: Relative paths are joined with default_start_path
	basePath := p.config.DefaultStartPath
	if basePath == "" {
		// If no default_start_path, use the current working directory.
		cwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("could not get current working directory: %w", err)
		}
		basePath = cwd
	}

	// The base path itself might need expansion (e.g., if it's "~/...").
	// We can use a temporary planner for this to avoid recursion issues.
	// A simpler way for now is to assume default_start_path is clean or absolute.
	// Let's expand it just in case.
	if strings.HasPrefix(basePath, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("cannot get user home directory for default_start_path: %w", err)
		}
		basePath = filepath.Join(homeDir, strings.TrimPrefix(basePath, "~"))
	}


	return filepath.Join(basePath, path), nil
}
