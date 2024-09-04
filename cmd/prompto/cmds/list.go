package cmds

import (
	"fmt"
	"github.com/go-go-golems/prompto/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all prompts in the repositories",
		RunE:  list,
	}
}

func list(cmd *cobra.Command, args []string) error {
	repositories := viper.GetStringSlice("repositories")

	for _, repoPath := range repositories {
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
