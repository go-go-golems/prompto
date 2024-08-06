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

	for _, repo := range repositories {
		files, err := pkg.GetFilesFromRepo(repo)
		if err != nil {
			return err
		}

		for _, file := range files {
			fmt.Println(repo, file.Name)
		}
	}

	return nil
}
