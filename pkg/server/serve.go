package server

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"time"

	"github.com/go-go-golems/prompto/pkg/server/handlers"
	"github.com/go-go-golems/prompto/pkg/server/state"
)

//go:embed static
var staticFiles embed.FS

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

	// Serve embedded static files
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return fmt.Errorf("error setting up static file system: %w", err)
	}
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	// Create server with timeouts
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           mux,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}

	fmt.Printf("Server is running on http://localhost:%d\n", port)
	return srv.ListenAndServe()
}
