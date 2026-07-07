package executor

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/sass1997/repolift/internal/planner"
)

// ExecutePlan runs the actions defined in a plan.
func ExecutePlan(plan *planner.Plan) error {
	for _, action := range plan.Actions {
		switch action.Type {
		case planner.ActionCreateWorkspaceDir:
			fmt.Printf("📂 Creating directory: %s\n", action.Path)
			if err := os.MkdirAll(action.Path, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", action.Path, err)
			}
		case planner.ActionCloneRepository:
			fmt.Printf("⏳ Cloning %s into %s...\n", action.URL, action.Path)
			cmd := exec.Command("git", "clone", action.URL, action.Path)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("git clone failed for %s: %w", action.URL, err)
			}
			fmt.Printf("✅ Successfully cloned %s\n", action.Path)
		case planner.ActionNoOp:
			// Do nothing for no-op actions, they are just for information.
			fmt.Printf("✅ No operation needed for: %s\n", action.Path)
		default:
			return fmt.Errorf("unknown action type: %s", action.Type)
		}
	}
	return nil
}
