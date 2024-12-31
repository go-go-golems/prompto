package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-go-golems/prompto/pkg/server/handlers"
	"github.com/go-go-golems/prompto/pkg/server/state"
)

func Serve(port int, watching bool, repositories []string) error {
	state := state.NewServerState(watching)
	state.Repositories = repositories
	if err := state.LoadRepositories(); err != nil {
		return fmt.Errorf("error loading repositories: %w", err)
	}

	// Start watching repositories if watching is enabled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if watching {
		if err := state.WatchRepositories(ctx); err != nil {
			return fmt.Errorf("error watching repositories: %w", err)
		}
	}

	h := handlers.NewHandlers(state)

	// Set up routes
	mux := http.NewServeMux()
	mux.Handle("/", h.Index())
	mux.Handle("/prompts/{name...}", h.PromptContent())
	mux.Handle("/search", h.Search())
	mux.Handle("/refresh", h.Refresh())
	mux.Handle("/repositories", h.Repositories())

	// Serve static files
	fs := http.FileServer(http.Dir("pkg/server/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Printf("Server is running on http://localhost:%d\n", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}
