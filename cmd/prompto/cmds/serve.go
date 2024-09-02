package cmds

import (
	"github.com/go-go-golems/prompto/pkg/server"
	"github.com/spf13/cobra"
)

func NewServeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start a web server to serve prompts",
		RunE:  serve,
	}

	cmd.Flags().Int("port", 8080, "Port to run the server on")

	return cmd
}

func serve(cmd *cobra.Command, args []string) error {
	port, _ := cmd.Flags().GetInt("port")

	return server.Serve(port)
}
