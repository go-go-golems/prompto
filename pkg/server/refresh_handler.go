package server

import (
	"net/http"
)

func refreshHandler(state *ServerState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if err := state.LoadRepositories(); err != nil {
			http.Error(w, "Error refreshing repositories", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Repositories refreshed successfully"))
	}
}