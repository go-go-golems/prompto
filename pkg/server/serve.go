package server

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/go-go-golems/prompto/pkg"
	"github.com/spf13/viper"
)

type ServerState struct {
	Repositories []string
	Files        map[string][]pkg.FileInfo
	mu           sync.RWMutex
}

func NewServerState() *ServerState {
	return &ServerState{
		Files: make(map[string][]pkg.FileInfo),
	}
}

func (s *ServerState) LoadRepositories() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Repositories = viper.GetStringSlice("repositories")
	for _, repo := range s.Repositories {
		files, err := pkg.GetFilesFromRepo(repo)
		if err != nil {
			return fmt.Errorf("error loading files from repository %s: %w", repo, err)
		}
		s.Files[repo] = files
	}
	return nil
}

func Serve(port int) error {
	state := NewServerState()
	if err := state.LoadRepositories(); err != nil {
		return fmt.Errorf("error loading repositories: %w", err)
	}

	http.HandleFunc("/", logHandler(rootHandler(state)))
	http.HandleFunc("/prompts/", logHandler(promptHandler(state)))
	http.HandleFunc("/search", logHandler(searchHandler(state)))
	http.HandleFunc("/refresh", logHandler(refreshHandler(state)))
	http.HandleFunc("/repositories", logHandler(repositoriesHandler(state)))

	fmt.Printf("Server is running on http://localhost:%d\n", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func logHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next(w, r)
	}
}
