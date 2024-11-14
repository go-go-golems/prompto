package cmds

import (
	"fmt"
	"github.com/go-go-golems/prompto/pkg"
	"github.com/spf13/cobra"
)

type ListCommand struct {
	repositories []string
}

func NewListCommand(options *CommandOptions) *cobra.Command {
	listCmd := &ListCommand{
		repositories: options.Repositories,
	}

	return &cobra.Command{
		Use:   "list",
		Short: "List all prompts in the repositories",
		RunE:  listCmd.run,
	}
}

func (l *ListCommand) run(cmd *cobra.Command, args []string) error {
	for _, repoPath := range l.repositories {
		repo := pkg.NewRepository(repoPath)
		err := repo.LoadPromptos()
		if err != nil {
			return err
		}

		for _, file := range repo.Promptos {
			fmt.Println(repoPath, file.Name)
		}
	}

	return nil
}
