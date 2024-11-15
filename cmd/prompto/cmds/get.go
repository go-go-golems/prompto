package cmds

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	files2 "github.com/go-go-golems/glazed/pkg/helpers/files"
	"github.com/go-go-golems/prompto/pkg"
	"github.com/spf13/cobra"
)

type GetCommand struct {
	repositories []string
}

func NewGetCommand(options *CommandOptions) *cobra.Command {
	getCmd := &GetCommand{
		repositories: options.Repositories,
	}

	cmd := &cobra.Command{
		Use:   "get prompt [args]",
		Short: "Get and execute or print a file",
		RunE:  getCmd.run,
	}

	cmd.Flags().Bool("print-path", false, "Print the path of the prompt")

	return cmd
}

func (g *GetCommand) interactiveGet() error {
	var allFiles []pkg.Prompto
	var selectedPrompt string
	var searchTerm string

	// Gather all files from repositories
	for _, repoPath := range g.repositories {
		repo := pkg.NewRepository(repoPath)
		err := repo.LoadPromptos()
		if err != nil {
			return err
		}
		allFiles = append(allFiles, repo.GetPromptos()...)
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Search for a prompt").
				Value(&searchTerm),
			huh.NewSelect[string]().
				Title("Select a prompt").
				Value(&selectedPrompt).
				OptionsFunc(func() []huh.Option[string] {
					var filteredOptions []huh.Option[string]
					for _, file := range allFiles {
						if strings.Contains(strings.ToLower(file.Name), strings.ToLower(searchTerm)) {
							filteredOptions = append(filteredOptions, huh.NewOption(file.Name, file.Name))
						}
					}
					return filteredOptions
				}, &searchTerm),
		),
	)

	err := form.Run()
	if err != nil {
		return err
	}

	if selectedPrompt == "" {
		return fmt.Errorf("no prompt selected")
	}

	// Find the selected prompt in the already gathered files
	for _, file := range allFiles {
		if file.Name == selectedPrompt {
			// Find the repository for this file
			for _, repoPath := range g.repositories {
				s, err := file.Render(repoPath, []string{})
				if err == nil {
					fmt.Println(s)
					return nil
				}
			}
			return fmt.Errorf("error rendering selected prompt")
		}
	}

	return fmt.Errorf("selected prompt not found")
}

func (g *GetCommand) run(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return g.interactiveGet()
	}

	prompt := args[0]
	restArgs := args[1:]

	printPath, _ := cmd.Flags().GetBool("print-path")

	for _, repoPath := range g.repositories {
		repo := pkg.NewRepository(repoPath)
		err := repo.LoadPromptos()
		if err != nil {
			return err
		}

		for _, file := range repo.Promptos {
			if file.Name == prompt {
				if printPath {
					if file.Type == pkg.Plain {
						fmt.Println(filepath.Join(repoPath, "prompto", file.Name))
						continue
					}

					s, err := file.Render(repoPath, restArgs)
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
				s, err := file.Render(repoPath, restArgs)
				if err != nil {
					return err
				}
				fmt.Println(s)
			}
		}
	}

	return nil
}
