package server

import (
	_ "embed"
	"html/template"
	"net/http"
	"strings"

	"github.com/go-go-golems/prompto/pkg"
	"github.com/go-go-golems/prompto/pkg/server/state"
)

//go:embed static/js/favorites.js
var favoritesJS string

//go:embed static/templates/root.html
var rootTemplate string

//go:embed static/templates/repoList.html
var repoListTemplate string

func rootHandler(state_ *state.ServerState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		tmpl, err := state_.CreateTemplateWithFuncs("root", rootTemplate+repoListTemplate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := struct {
			Groups      []string
			FavoritesJS template.JS
		}{
			Groups:      state_.GetAllGroups(),
			FavoritesJS: template.JS(favoritesJS),
		}

		w.Header().Set("Content-Type", "text/html")
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func searchHandler(state_ *state.ServerState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		query := r.FormValue("search")
		results := make(map[string][]pkg.Prompto)

		for _, file := range state_.GetAllPromptos() {
			if strings.Contains(strings.ToLower(file.Name), strings.ToLower(query)) {
				group := strings.SplitN(file.Name, "/", 2)[0]
				results[group] = append(results[group], file)
			}
		}

		groups := make([]string, 0)
		for group := range results {
			groups = append(groups, group)
		}

		funcMap := template.FuncMap{
			"PromptosByGroup": func(group string) []pkg.Prompto {
				return results[group]
			},
		}

		tmpl, err := state_.CreateTemplateWithFuncs("repoList", repoListTemplate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl = tmpl.Funcs(funcMap)

		data := struct {
			Groups []string
		}{
			Groups: groups,
		}

		w.Header().Set("Content-Type", "text/html")
		err = tmpl.ExecuteTemplate(w, "repoList", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
