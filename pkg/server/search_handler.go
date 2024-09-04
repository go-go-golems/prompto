package server

import (
	"github.com/go-go-golems/prompto/pkg"
	"html/template"
	"net/http"
	"strings"
)

func searchHandler(state *ServerState) http.HandlerFunc {
	tmpl := template.Must(template.New("searchResults").Parse(`
		{{range $group, $promptos := .}}
			<h2>{{$group}}</h2>
			<ul>
				{{range $promptos}}
					<li><a href="/prompts/{{.Name}}">{{.Name}}</a></li>
				{{end}}
			</ul>
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
		results := make(map[string][]pkg.Prompto)

		state.mu.RLock()
		for _, files := range state.Files {
			for _, file := range files {
				if strings.Contains(strings.ToLower(file.Name), strings.ToLower(query)) {
					group := strings.SplitN(file.Name, "/", 2)[0]
					results[group] = append(results[group], file)
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
