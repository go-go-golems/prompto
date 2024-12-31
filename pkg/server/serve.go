package server

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"sort"
	"sync"

	"github.com/go-go-golems/clay/pkg/watcher"
	"github.com/go-go-golems/glazed/pkg/helpers/templating"
	"github.com/go-go-golems/prompto/pkg"
	"github.com/rs/zerolog/log"
)

type ServerState struct {
	Repositories []string
	Repos        map[string]*pkg.Repository
	mu           sync.RWMutex
	Watching     bool
}

func NewServerState(watching bool) *ServerState {
	return &ServerState{
		Repos:    make(map[string]*pkg.Repository),
		Watching: watching,
	}
}

func (s *ServerState) LoadRepositories() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, repoPath := range s.Repositories {
		repo := pkg.NewRepository(repoPath)
		err := repo.LoadPromptos()
		if err != nil {
			return fmt.Errorf("error loading files from repository %s: %w", repoPath, err)
		}
		s.Repos[repoPath] = repo
	}
	return nil
}

func (s *ServerState) CreateTemplateWithFuncs(name, tmpl string) (*template.Template, error) {
	funcMap := template.FuncMap{
		"PromptosByGroup": func(group string) []pkg.Prompto {
			return s.GetPromptosByGroup(group)
		},
		"AllRepositories": func() []string {
			return s.GetAllRepositories()
		},
		"AllPromptos": func() []pkg.Prompto {
			return s.GetAllPromptos()
		},
		"AllGroups": func() []string {
			return s.GetAllGroups()
		},
		"PromptosByRepository": func(repo string) []pkg.Prompto {
			return s.Repos[repo].GetPromptos()
		},
		"GroupsByRepository": func(repo string) []string {
			return s.GetGroupsByRepository(repo)
		},
		"PromptosForRepositoryAndGroup": func(repo, group string) []pkg.Prompto {
			return s.GetPromptosForRepositoryAndGroup(repo, group)
		},
	}

	return templating.CreateHTMLTemplate(name).
		Funcs(funcMap).
		Parse(tmpl)
}

func (s *ServerState) GetAllRepositories() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Repositories
}

func (s *ServerState) GetAllPromptos() []pkg.Prompto {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var allPromptos []pkg.Prompto
	for _, repo := range s.Repos {
		allPromptos = append(allPromptos, repo.GetPromptos()...)
	}

	sort.Slice(allPromptos, func(i, j int) bool {
		return allPromptos[i].Name < allPromptos[j].Name
	})

	return allPromptos
}

func (s *ServerState) GetAllGroups() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	groupSet := make(map[string]struct{})
	for _, repo := range s.Repos {
		for _, group := range repo.GetGroups() {
			groupSet[group] = struct{}{}
		}
	}

	var groups []string
	for group := range groupSet {
		groups = append(groups, group)
	}

	sort.Strings(groups)
	return groups
}

func (s *ServerState) GetPromptosByGroup(group string) []pkg.Prompto {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var groupPromptos []pkg.Prompto
	for _, repo := range s.Repos {
		groupPromptos = append(groupPromptos, repo.GetPromptosByGroup(group)...)
	}

	sort.Slice(groupPromptos, func(i, j int) bool {
		return groupPromptos[i].Name < groupPromptos[j].Name
	})

	return groupPromptos
}

func (s *ServerState) GetGroupsByRepository(repo string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Repos[repo].GetGroups()
}

func (s *ServerState) GetPromptosForRepositoryAndGroup(repo, group string) []pkg.Prompto {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Repos[repo].GetPromptosByGroup(group)
}

func (s *ServerState) WatchRepositories(ctx context.Context) error {
	if !s.Watching {
		return nil
	}

	for repoPath, repo := range s.Repos {
		options := []watcher.Option{
			watcher.WithWriteCallback(func(path string) error {
				log.Info().Msgf("File %s changed, reloading...", path)
				s.mu.Lock()
				defer s.mu.Unlock()
				return repo.AddPrompto(path)
			}),
			watcher.WithRemoveCallback(func(path string) error {
				log.Info().Msgf("File %s removed, removing from repository...", path)
				s.mu.Lock()
				defer s.mu.Unlock()
				return repo.RemovePrompto(path)
			}),
			watcher.WithPaths(filepath.Join(repoPath, "prompto")),
		}

		w := watcher.NewWatcher(options...)
		go func() {
			if err := w.Run(ctx); err != nil {
				log.Error().Err(err).Msg("Watcher error")
			}
		}()
	}

	return nil
}

func Serve(port int, watching bool, repositories []string) error {
	state := NewServerState(watching)
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
