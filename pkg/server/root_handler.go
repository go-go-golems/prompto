package server

import (
	"net/http"
)

func rootHandler(state *ServerState) http.HandlerFunc {
	tmplStr := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Prompto Repositories</title>
			<script src="https://unpkg.com/htmx.org@1.9.6"></script>
			<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/milligram/1.4.1/milligram.min.css">
			<style>
				body { padding: 20px; }
				.htmx-indicator { display: none; }
				.htmx-request .htmx-indicator { display: inline; }
				#repo-list { margin-top: 20px; }
			</style>
		</head>
		<body>
			<div class="container">
				<h1>Prompto Repositories</h1>
				<input type="text" 
					   placeholder="Search repositories..." 
					   name="search" 
					   hx-post="/search"
					   hx-trigger="keyup changed delay:200ms, search"
					   hx-target="#repo-list"
					   hx-indicator=".htmx-indicator">
				<span class="htmx-indicator">Searching...</span>
				<div id="repo-list">
					{{range $group := .Groups}}
						<h2>{{.}}</h2>
						<ul>
							{{ range (PromptosByGroup $group) }}
								<li><a href="/prompts/{{.Name}}">{{.Name}}</a></li>
							{{end}}
						</ul>
					{{end}}
				</div>
			</div>
		</body>
		</html>
	`

	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		// Use the CreateTemplateWithFuncs method from ServerState
		tmpl, err := state.CreateTemplateWithFuncs("root", tmplStr, state)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := struct {
			Groups []string
		}{
			Groups: state.GetAllGroups(),
		}

		w.Header().Set("Content-Type", "text/html")
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
