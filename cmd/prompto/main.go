package main

import (
	"fmt"
	clay "github.com/go-go-golems/clay/pkg"
	"os"

	"github.com/go-go-golems/glazed/pkg/help"
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
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// reinitialize the logger because we can now parse --log-level and co
		// from the command line flag
		err := clay.InitLogger()
		cobra.CheckErr(err)
	},
}

func initRootCmd() (*help.HelpSystem, error) {
	helpSystem := help.NewHelpSystem()
	err := doc.AddDocToHelpSystem(helpSystem)
	cobra.CheckErr(err)

	helpSystem.SetupCobraRootCommand(rootCmd)

	err = clay.InitViper("prompto", rootCmd)
	cobra.CheckErr(err)
	err = clay.InitLogger()
	cobra.CheckErr(err)

	return helpSystem, nil
}

func main() {
	helpSystem, err := initRootCmd()
	cobra.CheckErr(err)

	err = doc.AddDocToHelpSystem(helpSystem)
	cobra.CheckErr(err)

	rootCmd.AddCommand(cmds.NewGetCommand())
	rootCmd.AddCommand(cmds.NewListCommand())
	command, err := cmds.NewConfigGroupCommand(helpSystem)
	cobra.CheckErr(err)
	rootCmd.AddCommand(command)
	rootCmd.AddCommand(cmds.NewServeCommand())

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
