package cmds

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/go-go-golems/prompto/pkg"
	"github.com/spf13/cobra"
)

type EditCommand struct {
	repositories []string
}

func NewEditCommand(options *CommandOptions) *cobra.Command {
	editCmd := &EditCommand{
		repositories: options.Repositories,
	}

	return &cobra.Command{
		Use:   "edit prompt-name",
		Short: "Edit a prompt in your default editor",
		Args:  cobra.ExactArgs(1),
		RunE:  editCmd.run,
	}
}

func openInEditor(filepath string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim" // Default to vim
	}

	cmd := exec.Command(editor, filepath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (e *EditCommand) run(cmd *cobra.Command, args []string) error {
	promptName := args[0]

	for _, repoPath := range e.repositories {
		repo := pkg.NewRepository(repoPath)
		err := repo.LoadPromptos()
		if err != nil {
			return err
		}

		for _, file := range repo.Promptos {
			if file.Name == promptName {
				fullPath := filepath.Join(repoPath, "prompto", file.Name)
				err := openInEditor(fullPath)
				if err != nil {
					return fmt.Errorf("failed to open editor: %w", err)
				}
				return nil
			}
		}
	}

	return fmt.Errorf("prompt '%s' not found in any repository", promptName)
}
