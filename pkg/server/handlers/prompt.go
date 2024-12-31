package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-go-golems/prompto/pkg"
	"github.com/go-go-golems/prompto/pkg/server/templates/components"
	"github.com/rs/zerolog/log"
)

func (h *Handlers) PromptList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 4 {
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}

		repo := parts[2]
		group := parts[3]

		repository, ok := h.state.Repos[repo]
		if !ok {
			http.Error(w, "Repository not found", http.StatusNotFound)
			return
		}

		prompts := repository.GetPromptosByGroup(group)
		component := components.PromptList(prompts)
		err := component.Render(r.Context(), w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (h *Handlers) PromptContent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := log.With().Str("handler", "PromptContent").Logger()

		name := r.PathValue("name")
		logger.Debug().Str("path", name).Msg("handling prompt request")

		if name == "" {
			logger.Debug().Msg("invalid path: empty name")
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}

		// Handle root directory listing
		if name == "" {
			logger.Debug().Msg("rendering root directory listing")
			var allPrompts []pkg.Prompto
			for _, repo := range h.state.Repos {
				allPrompts = append(allPrompts, repo.GetPromptos()...)
			}
			component := components.PromptList(allPrompts)
			err := component.Render(r.Context(), w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		// Split the path into group and prompt path
		parts := strings.SplitN(name, "/", 2)

		// Handle group listing
		if len(parts) == 1 {
			group := parts[0]
			logger = logger.With().Str("group", group).Logger()
			logger.Debug().Msg("rendering group listing")

			// Get all prompts from this group across all repositories
			var prompts []pkg.Prompto
			for _, repo := range h.state.Repos {
				prompts = append(prompts, repo.GetPromptosByGroup(group)...)
			}

			component := components.PromptList(prompts)
			err := component.Render(r.Context(), w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		group := parts[0]
		promptPath := name
		logger = logger.With().
			Str("group", group).
			Str("promptPath", promptPath).
			Logger()

		logger.Debug().Msg("looking up prompt")

		// Get all prompts for this group and find the matching one
		files := h.state.GetPromptosByGroup(group)
		var foundFile pkg.Prompto
		for _, file := range files {
			if file.Name == promptPath {
				foundFile = file
				break
			}
		}

		if foundFile.Name == "" {
			logger.Debug().Msg("prompt not found")
			http.Error(w, "Prompt not found", http.StatusNotFound)
			return
		}

		logger = logger.With().
			Str("repository", foundFile.Repository).
			Str("prompt", foundFile.Name).
			Logger()
		logger.Debug().Msg("found prompt, rendering with args")

		// Extract URL parameters
		queryParams := r.URL.Query()
		var restArgs []string
		for key, values := range queryParams {
			for _, value := range values {
				if value == "" {
					// pass non-keyword arguments as a straight string
					restArgs = append(restArgs, key)
				} else {
					restArgs = append(restArgs, fmt.Sprintf("--%s", key), value)
				}
			}
		}

		content, err := foundFile.Render(foundFile.Repository, restArgs)
		if err != nil {
			logger.Debug().Err(err).Msg("error rendering prompt")
			http.Error(w, "Error rendering prompt", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/markdown")
		_, _ = w.Write([]byte(content))
	}
}

func (h *Handlers) Search() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		var matchingPrompts []pkg.Prompto

		if query == "" {
			// Show all prompts when query is empty
			repoNames := h.state.GetAllRepositories()
			for _, repoName := range repoNames {
				repo := h.state.Repos[repoName]
				matchingPrompts = append(matchingPrompts, repo.GetPromptos()...)
			}
		} else {
			// Search for matching prompts
			for _, repo := range h.state.Repos {
				for _, prompt := range repo.GetPromptos() {
					if strings.Contains(strings.ToLower(prompt.Name), strings.ToLower(query)) ||
						strings.Contains(strings.ToLower(prompt.Group), strings.ToLower(query)) {
						matchingPrompts = append(matchingPrompts, prompt)
					}
				}
			}
		}

		component := components.PromptList(matchingPrompts)
		err := component.Render(r.Context(), w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
