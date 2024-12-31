package cmds

import (
	"github.com/go-go-golems/prompto/pkg/server"
	"github.com/spf13/cobra"
)

type ServeCommand struct {
	repositories []string
}

func NewServeCommand(options *CommandOptions) *cobra.Command {
	serveCmd := &ServeCommand{
		repositories: options.Repositories,
	}

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start a web server to serve prompts",
		RunE:  serveCmd.run,
	}

	cmd.Flags().Int("port", 8080, "Port to run the server on")
	cmd.Flags().Bool("watching", true, "Watch for changes to the repositories")

	return cmd
}

func (s *ServeCommand) run(cmd *cobra.Command, args []string) error {
	port, _ := cmd.Flags().GetInt("port")
	watching, _ := cmd.Flags().GetBool("watching")

	return server.Serve(port, watching, s.repositories)
}
