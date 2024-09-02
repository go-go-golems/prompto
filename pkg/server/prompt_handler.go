package server

import (
	"github.com/go-go-golems/prompto/pkg"
	"html/template"
	"net/http"
	"strings"
)

func promptHandler(state *ServerState) http.HandlerFunc {
	listTmpl := template.Must(template.New("promptList").Parse(`
		<ul id="content-list">
			{{range .}}
				<li><a href="/prompts/{{$.Repo}}/{{.Name}}">{{.Name}}</a></li>
			{{end}}
		</ul>
	`))

	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/prompts/")
		parts := strings.SplitN(path, "/", 2)

		if len(parts) == 1 {
			// Directory listing
			repo := parts[0]
			state.mu.RLock()
			files, ok := state.Files[repo]
			state.mu.RUnlock()

			if !ok {
				http.Error(w, "Repository not found", http.StatusNotFound)
				return
			}

			w.Header().Set("Content-Type", "text/html")
			err := listTmpl.Execute(w, struct {
				Repo  string
				Files []pkg.FileInfo
			}{
				Repo:  repo,
				Files: files,
			})
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

		repo := parts[0]
		promptName := parts[1]

		state.mu.RLock()
		files, ok := state.Files[repo]
		state.mu.RUnlock()

		if !ok {
			http.Error(w, "Repository not found", http.StatusNotFound)
			return
		}

		var foundFile pkg.FileInfo
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

		content, err := pkg.RenderFile(repo, foundFile, []string{})
		if err != nil {
			http.Error(w, "Error rendering prompt", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/markdown")
		_, _ = w.Write([]byte(content))
	}
}
