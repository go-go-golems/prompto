package cmds

import (
	"fmt"
	"path/filepath"

	"github.com/go-go-golems/prompto/pkg"
	"github.com/spf13/cobra"
)

type WhichCommand struct {
	repositories []string
}

func NewWhichCommand(options *CommandOptions) *cobra.Command {
	whichCmd := &WhichCommand{
		repositories: options.Repositories,
	}

	return &cobra.Command{
		Use:   "which prompt-name",
		Short: "Show the source location of a prompt",
		Args:  cobra.ExactArgs(1),
		RunE:  whichCmd.run,
	}
}

func (w *WhichCommand) run(cmd *cobra.Command, args []string) error {
	promptName := args[0]

	for _, repoPath := range w.repositories {
		repo := pkg.NewRepository(repoPath)
		err := repo.LoadPromptos()
		if err != nil {
			return err
		}

		for _, file := range repo.Promptos {
			if file.Name == promptName {
				fullPath := filepath.Join(repoPath, "prompto", file.Name)
				fmt.Printf("Repository: %s\nFile: %s\nFull path: %s\n", repoPath, file.Name, fullPath)
				return nil
			}
		}
	}

	return fmt.Errorf("prompt '%s' not found in any repository", promptName)
}
