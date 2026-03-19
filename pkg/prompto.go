package pkg

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/fields"
	"github.com/go-go-golems/glazed/pkg/cmds/schema"
	"github.com/go-go-golems/glazed/pkg/cmds/values"
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

		parsedValues := values.New()
		err := p.Command.Schema.ForEachE(func(_ string, section schema.Section) error {
			definitions := section.GetDefinitions()

			parsedFlags, args_, err := definitions.GatherFlagsFromStringList(
				restArgs, false, false, "",
				fields.WithSource("command-line"),
			)
			if err != nil {
				return err
			}

			arguments := fields.NewFieldValues()
			if len(args_) > 0 {
				if section.GetSlug() != schema.DefaultSlug {
					return errors.Errorf("section %s does not accept arguments (only default section can)", section.GetSlug())
				}

				argumentDefinitions := p.Command.GetDefaultArguments()
				arguments, err = argumentDefinitions.GatherArguments(
					args_, false, false,
					fields.WithSource("cli"),
				)
				if err != nil {
					return err
				}
			}

			sectionValues, err := values.NewSectionValues(section,
				values.WithFields(parsedFlags),
				values.WithFields(arguments),
			)
			if err != nil {
				return err
			}

			parsedValues.Set(section.GetSlug(), sectionValues)

			return nil
		})
		if err != nil {
			return "", err
		}

		err = p.Command.RunIntoWriter(context.Background(), parsedValues, buf)
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
