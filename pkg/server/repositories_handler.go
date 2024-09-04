package server

import (
	"encoding/json"
	"net/http"
)

func repositoriesHandler(state *ServerState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state.mu.RLock()
		repos := state.Repositories
		state.mu.RUnlock()

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(repos)
		if err != nil {
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
			return
		}
	}
}
