package main

import (
	"fmt"
	"github.com/go-go-golems/glazed/pkg/cmds"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type FileType int

const (
	Plain FileType = iota
	Executable
	TemplateCommand
)

type FileInfo struct {
	Name string
	Type FileType
}

func isTemplateCommand(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	tcl := &cmds.TemplateCommandLoader{}
	_, err = tcl.LoadCommandFromYAML(f)
	return err == nil
}

func getFilesFromRepo(repo string) ([]FileInfo, error) {
	var files []FileInfo

	err := filepath.Walk(filepath.Join(repo, "prompto"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// strip the repo path from the file name
		relpath := path[len(repo):]
		name := strings.TrimPrefix(relpath, "/prompto/")

		if !info.IsDir() {
			file := FileInfo{Name: name}

			if (info.Mode() & 0111) != 0 {
				file.Type = Executable
			} else if isTemplateCommand(path) {
				file.Type = TemplateCommand
			} else {
				file.Type = Plain
			}

			files = append(files, file)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "prompto",
		Short: "prompto generates prompts from a list of repositories",
		Long: `This program loads a list of repositories from a yaml config file
and looks for a file that matches the prompt.`,
	}

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

	repositories := viper.GetStringSlice("repositories")

	for _, repo := range repositories {
		files, err := getFilesFromRepo(repo)
		if err != nil {
			return err
		}

		for _, file := range files {
			if file.Name == prompt {
				path := filepath.Join(repo, "prompto", file.Name)
				if file.Type == Executable {
					// The file is executable, execute it
					c := exec.Command(path, restArgs...)
					c.Dir = repo
					c.Stdout = os.Stdout
					c.Stderr = os.Stderr
					return c.Run()
				} else {
					// The file is not executable, print its content
					b, err := os.ReadFile(path)
					if err != nil {
						return err
					}
					fmt.Println(string(b))
				}
			}
		}
	}

	return nil
}

func list(cmd *cobra.Command, args []string) error {
	repositories := viper.GetStringSlice("repositories")

	for _, repo := range repositories {
		files, err := getFilesFromRepo(repo)
		if err != nil {
			return err
		}

		for _, file := range files {
			fmt.Println(repo, file.Name)
		}
	}

	return nil
}
