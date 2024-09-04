package pkg

import (
	"context"
	"github.com/go-go-golems/clay/pkg/watcher"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Repository struct {
	Path     string
	Promptos []Prompto
}

func NewRepository(path string) *Repository {
	return &Repository{Path: path}
}

func (r *Repository) LoadPromptos() error {
	promptoDir := filepath.Join(r.Path, "prompto")

	if _, err := os.Stat(promptoDir); os.IsNotExist(err) {
		r.Promptos = []Prompto{}
		return nil
	}

	var promptos []Prompto

	err := filepath.Walk(promptoDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasPrefix(info.Name(), ".") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		relpath := strings.TrimPrefix(path, r.Path)
		relpath = strings.TrimPrefix(relpath, "/")
		name := strings.TrimPrefix(relpath, "prompto/")
		group := strings.SplitN(name, "/", 2)[0]

		if !info.IsDir() {
			prompto := Prompto{
				Name:       name,
				Group:      group,
				Type:       Plain,
				FilePath:   path,   // Store absolute file path
				Repository: r.Path, // Store repository
			}
			ext := strings.ToLower(filepath.Ext(name))

			isYAMLfile := ext == ".yaml" || ext == ".yml"

			if info.Mode()&os.ModeSymlink != 0 {
				info_, err := os.Stat(path)
				if err != nil {
					return err
				}
				if (info_.Mode() & 0111) != 0 {
					prompto.Type = Executable
				}
				prompto.Group = group
			} else {
				if (info.Mode() & 0111) != 0 {
					prompto.Type = Executable
				}
				prompto.Group = group
			}

			if prompto.Type != Executable {
				if isYAMLfile {
					if cmd, ok := LoadTemplateCommand(path); ok {
						prompto.Name = name[:len(name)-len(ext)]
						prompto.Type = TemplateCommand
						prompto.Command = cmd
					}
				}
			}

			promptos = append(promptos, prompto)
		}

		return nil
	})

	if err != nil {
		return err
	}

	r.Promptos = promptos
	return nil
}

func (r *Repository) Refresh() error {
	return r.LoadPromptos()
}

func (r *Repository) GroupPromptos() map[string][]Prompto {
	grouped := make(map[string][]Prompto)
	for _, prompto := range r.Promptos {
		group := strings.SplitN(prompto.Name, "/", 2)[0]
		grouped[group] = append(grouped[group], prompto)
	}
	return grouped
}

func (r *Repository) GetGroups() []string {
	grouped := r.GroupPromptos()
	var groups []string
	for group := range grouped {
		groups = append(groups, group)
	}
	sort.Strings(groups)
	return groups
}

func (r *Repository) GetPromptosByGroup(group string) []Prompto {
	grouped := r.GroupPromptos()
	promptos := grouped[group]
	sort.Slice(promptos, func(i, j int) bool {
		return promptos[i].Name < promptos[j].Name
	})
	return promptos
}

func (r *Repository) Watch(ctx context.Context, options ...watcher.Option) error {
	options = append(options,
		watcher.WithWriteCallback(func(path string) error {
			log.Debug().Msgf("Loading %s", path)
			return r.LoadPromptos()
		}),
		watcher.WithRemoveCallback(func(path string) error {
			log.Debug().Msgf("Removing %s", path)
			return r.LoadPromptos()
		}),
		watcher.WithPaths(r.Path),
	)

	w := watcher.NewWatcher(options...)
	err := w.Run(ctx)
	if err != nil {
		return err
	}
	return nil
}
