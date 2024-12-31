package handlers

import (
	"encoding/json"
	"net/http"
)

func (h *Handlers) Repositories() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		repos := h.state.GetAllRepositories()

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(repos)
		if err != nil {
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
			return
		}
	}
}
