package handlers

import (
	"net/http"

	"github.com/go-go-golems/prompto/pkg/server/templates/pages"
)

func (h *Handlers) Index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		component := pages.Index(h.state.GetAllRepositories(), h.state.Repos)
		component.Render(r.Context(), w)
	}
}
