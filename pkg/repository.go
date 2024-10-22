package pkg

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/go-go-golems/clay/pkg/watcher"
	"github.com/rs/zerolog/log"
)

type Repository struct {
	Path     string
	Promptos map[string]Prompto
}

func NewRepository(path string) *Repository {
	return &Repository{
		Path:     path,
		Promptos: make(map[string]Prompto),
	}
}

func (r *Repository) LoadPromptos() error {
	promptoDir := filepath.Join(r.Path, "prompto")

	if _, err := os.Stat(promptoDir); os.IsNotExist(err) {
		r.Promptos = make(map[string]Prompto)
		return nil
	}

	newPromptos := make(map[string]Prompto)

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

		prompto, err := r.loadSinglePrompto(path, info)
		if err != nil {
			return err
		}
		if prompto != nil {
			newPromptos[prompto.FilePath] = *prompto
		}

		return nil
	})

	if err != nil {
		return err
	}

	r.Promptos = newPromptos
	return nil
}

func (r *Repository) GetPromptos() []Prompto {
	var promptos []Prompto
	for _, prompto := range r.Promptos {
		promptos = append(promptos, prompto)
	}
	return promptos
}

func (r *Repository) loadSinglePrompto(path string, info os.FileInfo) (*Prompto, error) {
	if info.IsDir() {
		return nil, nil
	}

	relpath := strings.TrimPrefix(path, r.Path)
	relpath = strings.TrimPrefix(relpath, "/")
	name := strings.TrimPrefix(relpath, "prompto/")
	group := strings.SplitN(name, "/", 2)[0]

	prompto := &Prompto{
		Name:       name,
		Group:      group,
		Type:       Plain,
		FilePath:   path,
		Repository: r.Path,
	}

	ext := strings.ToLower(filepath.Ext(name))
	isYAMLfile := ext == ".yaml" || ext == ".yml"

	if info.Mode()&os.ModeSymlink != 0 {
		info_, err := os.Stat(path)
		if err != nil {
			log.Warn().Err(err).Str("path", path).Msg("Failed to stat symlink target")
		} else {
			if (info_.Mode() & 0111) != 0 {
				prompto.Type = Executable
			}
		}
		prompto.Group = group
	} else {
		if (info.Mode() & 0111) != 0 {
			prompto.Type = Executable
		}
		prompto.Group = group
	}

	if prompto.Type != Executable && isYAMLfile {
		if cmd, ok := LoadTemplateCommand(path); ok {
			prompto.Name = name[:len(name)-len(ext)]
			prompto.Type = TemplateCommand
			prompto.Command = cmd
		}
	}

	return prompto, nil
}

func (r *Repository) AddPrompto(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	prompto, err := r.loadSinglePrompto(path, info)
	if err != nil {
		return err
	}

	if prompto != nil {
		r.Promptos[prompto.FilePath] = *prompto
	}

	return nil
}

func (r *Repository) RemovePrompto(path string) error {
	delete(r.Promptos, path)
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
