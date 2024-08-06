package cmds

import (
	"fmt"
	files2 "github.com/go-go-golems/glazed/pkg/helpers/files"
	"github.com/go-go-golems/prompto/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

func NewGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get prompt [args]",
		Short: "Get and execute or print a file",
		Args:  cobra.MinimumNArgs(1),
		RunE:  get,
	}

	cmd.Flags().Bool("print-path", false, "Print the path of the prompt")

	return cmd
}

func get(cmd *cobra.Command, args []string) error {
	prompt := args[0]
	restArgs := args[1:]

	printPath, _ := cmd.Flags().GetBool("print-path")

	repositories := viper.GetStringSlice("repositories")

	for _, repo := range repositories {
		files, err := pkg.GetFilesFromRepo(repo)
		if err != nil {
			return err
		}

		for _, file := range files {
			if file.Name == prompt {
				if printPath {
					if file.Type == pkg.Plain {
						fmt.Println(filepath.Join(repo, "prompto", file.Name))
						continue
					}

					s, err := pkg.RenderFile(repo, file, restArgs)
					if err != nil {
						return err
					}

					// transform prompt into a safe filename by removing / and \
					promptString := strings.ReplaceAll(prompt, "/", "-")
					promptString = strings.ReplaceAll(promptString, "\\", "-")

					deletedFiles, err := files2.GarbageCollectTemporaryFiles(os.TempDir(), "prompto-*", 20)
					if err != nil {
						return err
					}

					_ = deletedFiles

					f, err := os.CreateTemp(os.TempDir(), fmt.Sprintf("prompto-%s-*", promptString))
					if err != nil {
						return err
					}

					_, err = f.Write([]byte(s))
					if err != nil {
						_ = f.Close()
						return err
					}
					_ = f.Close()

					fmt.Println(f.Name())

					continue
				}
				s, err := pkg.RenderFile(repo, file, restArgs)
				if err != nil {
					return err
				}
				fmt.Println(s)
			}
		}
	}

	return nil
}
