package main

import (
	"context"
	"fmt"
	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/layers"
	"github.com/go-go-golems/glazed/pkg/cmds/parameters"
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
	Name    string
	Type    FileType
	Command *cmds.TemplateCommand
}

func loadTemplateCommand(path string) (*cmds.TemplateCommand, bool) {
	f, err := os.Open(path)
	if err != nil {
		return nil, false
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	tcl := &cmds.TemplateCommandLoader{}
	commands, err := tcl.LoadCommandFromYAML(f)
	if err != nil {
		return nil, false
	}
	if len(commands) != 1 {
		return nil, false
	}
	return commands[0].(*cmds.TemplateCommand), err == nil
}

func getFilesFromRepo(repo string) ([]FileInfo, error) {
	promptoDir := filepath.Join(repo, "prompto")

	// check if repo/prompto exists, else return empty list
	if _, err := os.Stat(promptoDir); os.IsNotExist(err) {
		return []FileInfo{}, nil
	}
	var files []FileInfo

	err := filepath.Walk(promptoDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// skip . files and directories
		if strings.HasPrefix(info.Name(), ".") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// strip the repo path from the file name
		relpath := path[len(repo):]
		name := strings.TrimPrefix(relpath, "/prompto/")

		if !info.IsDir() {
			file := FileInfo{Name: name}
			ext := strings.ToLower(filepath.Ext(name))

			isYAMLfile := ext == ".yaml" || ext == ".yml"

			if (info.Mode() & 0111) != 0 {
				file.Type = Executable
			} else if isYAMLfile {
				if cmd, ok := loadTemplateCommand(path); ok {
					file.Name = name[:len(name)-len(ext)]
					file.Type = TemplateCommand
					file.Command = cmd
				} else {
					file.Type = Plain
				}
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
				switch file.Type {
				case Executable:
					// The file is executable, execute it
					c := exec.Command(path, restArgs...)
					c.Dir = repo
					c.Stdout = os.Stdout
					c.Stderr = os.Stderr
					return c.Run()
				case TemplateCommand:
					// The file is a glazed TemplateCommand, execute it by passing the arguments
					buf := &strings.Builder{}
					parsedLayers := map[string]*layers.ParsedParameterLayer{}
					ps, args_, err := parameters.GatherFlagsFromStringList(
						restArgs, file.Command.Flags,
						false, false,
						"")
					if err != nil {
						return err
					}
					arguments, err := parameters.GatherArguments(args_, file.Command.Arguments, false, false)
					if err != nil {
						return err
					}
					for p := arguments.Oldest(); p != nil; p = p.Next() {
						k, v := p.Key, p.Value
						ps[k] = v
					}
					err = file.Command.RunIntoWriter(context.Background(), parsedLayers, ps, buf)
					if err != nil {
						return err
					}
					fmt.Println(buf.String())
				case Plain:
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
