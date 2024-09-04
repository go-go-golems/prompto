package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"

	"github.com/go-go-golems/glazed/pkg/helpers/templating"
	"github.com/go-go-golems/prompto/pkg"
	"github.com/spf13/viper"
)

type ServerState struct {
	Repositories []string
	Files        map[string][]pkg.Prompto
	mu           sync.RWMutex
}

func NewServerState() *ServerState {
	return &ServerState{
		Files: make(map[string][]pkg.Prompto),
	}
}

func (s *ServerState) LoadRepositories() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Repositories = viper.GetStringSlice("repositories")
	for _, repoPath := range s.Repositories {
		repo := pkg.NewRepository(repoPath)
		err := repo.LoadPromptos()
		if err != nil {
			return fmt.Errorf("error loading files from repository %s: %w", repoPath, err)
		}
		s.Files[repoPath] = repo.Promptos
	}
	return nil
}

func (s *ServerState) CreateTemplateWithFuncs(name, tmpl string, state *ServerState) (*template.Template, error) {
	funcMap := template.FuncMap{
		"PromptosByGroup": func(group string) []pkg.Prompto {
			return state.GetPromptosByGroup(group)
		},
		"AllRepositories": func() []string {
			return state.GetAllRepositories()
		},
		"AllPromptos": func() []pkg.Prompto {
			return state.GetAllPromptos()
		},
		"AllGroups": func() []string {
			return state.GetAllGroups()
		},
		"PromptosByRepository": func(repo string) []pkg.Prompto {
			return state.GetPromptosByRepository(repo)
		},
		"GroupsByRepository": func(repo string) []string {
			return state.GetGroupsByRepository(repo)
		},
		"PromptosForRepositoryAndGroup": func(repo, group string) []pkg.Prompto {
			return state.GetPromptosForRepositoryAndGroup(repo, group)
		},
	}

	return templating.CreateHTMLTemplate(name).
		Funcs(funcMap).
		Parse(tmpl)
}

// New methods to return data useful for rendering templates
func (s *ServerState) GetAllRepositories() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Repositories
}

func (s *ServerState) GetAllPromptos() []pkg.Prompto {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var allPromptos []pkg.Prompto
	for _, promptos := range s.Files {
		allPromptos = append(allPromptos, promptos...)
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
	for _, promptos := range s.Files {
		for _, prompto := range promptos {
			group := strings.SplitN(prompto.Name, "/", 2)[0]
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
	for _, promptos := range s.Files {
		for _, prompto := range promptos {
			if strings.HasPrefix(prompto.Name, group+"/") {
				groupPromptos = append(groupPromptos, prompto)
			}
		}
	}

	sort.Slice(groupPromptos, func(i, j int) bool {
		return groupPromptos[i].Name < groupPromptos[j].Name
	})

	return groupPromptos
}

func (s *ServerState) GetPromptosByRepository(repo string) []pkg.Prompto {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Files[repo]
}

func (s *ServerState) GetGroupsByRepository(repo string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	groupSet := make(map[string]struct{})
	for _, prompto := range s.Files[repo] {
		group := strings.SplitN(prompto.Name, "/", 2)[0]
		groupSet[group] = struct{}{}
	}

	var groups []string
	for group := range groupSet {
		groups = append(groups, group)
	}

	sort.Strings(groups)
	return groups
}

// New function to get promptos for a repository and group
func (s *ServerState) GetPromptosForRepositoryAndGroup(repo, group string) []pkg.Prompto {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var groupPromptos []pkg.Prompto
	for _, prompto := range s.Files[repo] {
		if strings.HasPrefix(prompto.Name, group+"/") {
			groupPromptos = append(groupPromptos, prompto)
		}
	}

	sort.Slice(groupPromptos, func(i, j int) bool {
		return groupPromptos[i].Name < groupPromptos[j].Name
	})

	return groupPromptos
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
