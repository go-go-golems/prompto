package pkg

import (
	"github.com/go-go-golems/glazed/pkg/cmds"
	"os"
	"path/filepath"
	"strings"
)

type FileType int

const (
	Plain FileType = iota
	Executable
	TemplateCommand
)

type FileInfo struct {
	Name    string
	Type    FileType
	Command *cmds.TemplateCommand
}

func LoadTemplateCommand(path string) (*cmds.TemplateCommand, bool) {
	f, err := os.Open(path)
	if err != nil {
		return nil, false
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	tcl := &cmds.TemplateCommandLoader{}
	commands, err := tcl.LoadCommandFromYAML(f)
	if err != nil {
		return nil, false
	}
	if len(commands) != 1 {
		return nil, false
	}
	return commands[0].(*cmds.TemplateCommand), err == nil
}

func GetFilesFromRepo(repo string) ([]FileInfo, error) {
	promptoDir := filepath.Join(repo, "prompto")

	// check if repo/prompto exists, else return empty list
	if _, err := os.Stat(promptoDir); os.IsNotExist(err) {
		return []FileInfo{}, nil
	}
	var files []FileInfo

	err := filepath.Walk(promptoDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// skip . files and directories
		if strings.HasPrefix(info.Name(), ".") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// strip the repo path from the file name
		relpath := path[len(repo):]
		name := strings.TrimPrefix(relpath, "/prompto/")

		if !info.IsDir() {
			file := FileInfo{Name: name, Type: Plain}
			ext := strings.ToLower(filepath.Ext(name))

			isYAMLfile := ext == ".yaml" || ext == ".yml"

			// check for symbolic link
			if info.Mode()&os.ModeSymlink != 0 {
				// check if the link source is executable
				info_, err := os.Stat(path)
				if err != nil {
					return err
				}
				if (info_.Mode() & 0111) != 0 {
					file.Type = Executable
				}
			} else {
				if (info.Mode() & 0111) != 0 {
					file.Type = Executable
				}
			}

			if file.Type != Executable {
				if isYAMLfile {
					if cmd, ok := LoadTemplateCommand(path); ok {
						file.Name = name[:len(name)-len(ext)]
						file.Type = TemplateCommand
						file.Command = cmd
					}
				}
			}

			files = append(files, file)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}
