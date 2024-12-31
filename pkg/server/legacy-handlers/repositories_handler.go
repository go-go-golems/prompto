package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-go-golems/prompto/pkg/server/state"
)

func repositoriesHandler(state_ *state.ServerState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		repos := state_.Repositories

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(repos)
		if err != nil {
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
			return
		}
	}
}
