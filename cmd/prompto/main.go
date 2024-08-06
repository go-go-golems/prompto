package main

import (
	"fmt"
	"github.com/go-go-golems/prompto/cmd/prompto/cmds"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "prompto",
		Short: "prompto generates prompts from a list of repositories",
		Long: `This program loads a list of repositories from a yaml config file
and looks for a file that matches the prompt.`,
	}

	rootCmd.AddCommand(cmds.NewGetCommand())
	rootCmd.AddCommand(cmds.NewListCommand())

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
