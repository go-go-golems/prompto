package handlers

import (
	"net/http"
	"strings"

	"github.com/go-go-golems/prompto/pkg"
	"github.com/go-go-golems/prompto/pkg/server/templates/components"
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
		component.Render(r.Context(), w)
	}
}

func (h *Handlers) Search() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		if query == "" {
			return
		}

		var matchingPrompts []pkg.Prompto
		for _, repo := range h.state.Repos {
			for _, prompt := range repo.GetPromptos() {
				if strings.Contains(strings.ToLower(prompt.Name), strings.ToLower(query)) ||
					strings.Contains(strings.ToLower(prompt.Group), strings.ToLower(query)) {
					matchingPrompts = append(matchingPrompts, prompt)
				}
			}
		}

		component := components.PromptList(matchingPrompts)
		component.Render(r.Context(), w)
	}
}
