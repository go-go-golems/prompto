package main

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

func main() {
	var rootCmd = &cobra.Command{
		Use:   "prompto",
		Short: "prompto generates prompts from a list of repositories",
		Long: `This program loads a list of repositories from a yaml config file
and looks for a file that matches the prompt.`,
	}

	getCmd.Flags().Bool("print-path", false, "Print the path of the prompt")

	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(listCmd)

	viper.SetConfigName("config")
	viper.AddConfigPath(os.ExpandEnv("$HOME/.prompto"))
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var getCmd = &cobra.Command{
	Use:   "get prompt [args]",
	Short: "Get and execute or print a file",
	Args:  cobra.MinimumNArgs(1),
	RunE:  get,
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all prompts in the repositories",
	RunE:  list,
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

					//_, _ = fmt.Fprintf(os.Stderr, "Deleted %d files\n", len(deletedFiles))
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
