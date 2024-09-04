package server

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/go-go-golems/prompto/pkg"
)

func promptHandler(state *ServerState) http.HandlerFunc {
	listTmpl := template.Must(template.New("promptList").Parse(`
		<ul id="content-list">
			{{range .}}
				<li><a href="/prompts/{{.Name}}">{{.Name}}</a></li>
			{{end}}
		</ul>
	`))

	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/prompts/"), "/")
		parts := strings.SplitN(path, "/", 2)

		if len(parts) == 1 {
			// Directory listing
			group := parts[0]
			state.mu.RLock()
			files := state.GetPromptosByGroup(group)
			state.mu.RUnlock()

			w.Header().Set("Content-Type", "text/html")
			err := listTmpl.Execute(w, files)
			if err != nil {
				http.Error(w, "Error rendering template", http.StatusInternalServerError)
				return
			}
			return
		}

		if len(parts) != 2 {
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}

		group := parts[0]
		promptName := path

		state.mu.RLock()
		files := state.GetPromptosByGroup(group)
		state.mu.RUnlock()

		var foundFile pkg.Prompto
		for _, file := range files {
			if file.Name == promptName {
				foundFile = file
				break
			}
		}

		if foundFile.Name == "" {
			http.Error(w, "Prompt not found", http.StatusNotFound)
			return
		}

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
			http.Error(w, "Error rendering prompt", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/markdown")
		_, _ = w.Write([]byte(content))
	}
}
