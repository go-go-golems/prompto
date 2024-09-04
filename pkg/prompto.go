package pkg

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/layers"
	parameters2 "github.com/go-go-golems/glazed/pkg/cmds/parameters"
	"github.com/pkg/errors"
)

type FileType int

const (
	Plain FileType = iota
	Executable
	TemplateCommand
)

type Prompto struct {
	Name       string
	Group      string
	Type       FileType
	Command    *cmds.TemplateCommand
	FilePath   string // New field to store absolute file path
	Repository string // New field to store repository
}

func (p *Prompto) Render(repo string, restArgs []string) (string, error) {
	path := filepath.Join(repo, "prompto", p.Name)

	switch p.Type {
	case Executable:
		c := exec.Command(path, restArgs...)
		c.Dir = repo
		var out bytes.Buffer
		c.Stdout = &out

		currentDir, err := os.Getwd()
		if err != nil {
			return "", err
		}

		c.Env = append(os.Environ(), "PROMPTO_PARENT_PWD="+currentDir)

		err = c.Run()
		if err != nil {
			return "", err
		}
		return out.String(), nil
	case TemplateCommand:
		buf := &strings.Builder{}

		parsedLayers := layers.NewParsedLayers()
		err := p.Command.Layers.ForEachE(func(_ string, l layers.ParameterLayer) error {
			parameters := l.GetParameterDefinitions()

			parsedFlags, args_, err := parameters.GatherFlagsFromStringList(
				restArgs, false, false, "",
				parameters2.WithParseStepSource("command-line"),
			)
			if err != nil {
				return err
			}

			arguments := parameters2.NewParsedParameters()
			if len(args_) > 0 {
				if l.GetSlug() != layers.DefaultSlug {
					return errors.Errorf("layer %s does not accept arguments (only default layer can)", l.GetSlug())
				}

				argumentDefinitions := p.Command.GetDefaultArguments()
				arguments, err = argumentDefinitions.GatherArguments(
					args_, false, false,
					parameters2.WithParseStepSource("cli"),
				)
				if err != nil {
					return err
				}
			}

			parsedLayer, err := layers.NewParsedLayer(l,
				layers.WithParsedParameters(parsedFlags),
				layers.WithParsedParameters(arguments),
			)
			if err != nil {
				return err
			}

			parsedLayers.Set(l.GetSlug(), parsedLayer)

			return nil
		})
		if err != nil {
			return "", err
		}

		err = p.Command.RunIntoWriter(context.Background(), parsedLayers, buf)
		if err != nil {
			return "", err
		}
		return buf.String(), nil
	case Plain:
		b, err := os.ReadFile(path)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}

	return "", errors.Errorf("unsupported file type: %v", p.Type)
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
