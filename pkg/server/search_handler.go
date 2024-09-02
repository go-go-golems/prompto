package server

import (
	"html/template"
	"net/http"
	"strings"
)

func searchHandler(state *ServerState) http.HandlerFunc {
	tmpl := template.Must(template.New("searchResults").Parse(`
		{{range .}}
			<li><a href="/prompts/{{.Repo}}/{{.Name}}">{{.Repo}}/{{.Name}}</a></li>
		{{else}}
			<li>No results found</li>
		{{end}}
	`))

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		query := r.FormValue("search")
		results := []struct {
			Repo string
			Name string
		}{}

		state.mu.RLock()
		for repo, files := range state.Files {
			for _, file := range files {
				if strings.Contains(strings.ToLower(file.Name), strings.ToLower(query)) {
					results = append(results, struct {
						Repo string
						Name string
					}{
						Repo: repo,
						Name: file.Name,
					})
				}
			}
		}
		state.mu.RUnlock()

		w.Header().Set("Content-Type", "text/html")
		err := tmpl.Execute(w, results)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
