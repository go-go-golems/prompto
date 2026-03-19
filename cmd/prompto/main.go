package main

import (
	"fmt"
	"os"
	"strings"

	clay "github.com/go-go-golems/clay/pkg"
	clay_repositories "github.com/go-go-golems/clay/pkg/cmds/repositories"
	"github.com/go-go-golems/glazed/pkg/cmds/logging"
	"github.com/go-go-golems/glazed/pkg/help"
	helpCmd "github.com/go-go-golems/glazed/pkg/help/cmd"
	"github.com/go-go-golems/prompto/cmd/prompto/cmds"
	"github.com/go-go-golems/prompto/pkg/doc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "prompto",
	Short: "prompto generates prompts from a list of repositories",
	Long: `This program loads a list of repositories from a yaml config file
and looks for a file that matches the prompt.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return logging.InitLoggerFromCobra(cmd)
	},
}

func initRootCmd() (*help.HelpSystem, error) {
	helpSystem := help.NewHelpSystem()
	err := doc.AddDocToHelpSystem(helpSystem)
	cobra.CheckErr(err)

	helpCmd.SetupCobraRootCommand(helpSystem, rootCmd)

	err = clay.InitGlazed("prompto", rootCmd)
	cobra.CheckErr(err)

	rootCmd.AddCommand(clay_repositories.NewRepositoriesGroupCommand())

	return helpSystem, nil
}

func main() {
	helpSystem, err := initRootCmd()
	cobra.CheckErr(err)

	err = doc.AddDocToHelpSystem(helpSystem)
	cobra.CheckErr(err)

	viper.SetConfigName("config")
	viper.AddConfigPath(os.ExpandEnv("$HOME/.prompto"))
	viper.SetEnvPrefix("prompto")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	repositories := viper.GetStringSlice("repositories")
	options := cmds.NewCommandOptions(repositories)

	for _, cmd := range cmds.NewCommands(options) {
		rootCmd.AddCommand(cmd)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
