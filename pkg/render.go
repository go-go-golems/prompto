package pkg

import (
	"bytes"
	"context"
	"fmt"
	"github.com/go-go-golems/glazed/pkg/cmds/layers"
	parameters2 "github.com/go-go-golems/glazed/pkg/cmds/parameters"
	"gopkg.in/errgo.v2/fmt/errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func RenderFile(repo string, file FileInfo, restArgs []string) (string, error) {
	// Define the path to the file
	path := filepath.Join(repo, "prompto", file.Name)

	switch file.Type {
	case Executable:
		// The file is executable, execute it
		c := exec.Command(path, restArgs...)
		c.Dir = repo
		var out bytes.Buffer
		c.Stdout = &out

		// Get the current working directory
		currentDir, err := os.Getwd()
		if err != nil {
			return "", err
		}

		// Set the PROMPTO_PARENT_PWD environment variable to the current directory
		c.Env = append(os.Environ(), "PROMPTO_PARENT_PWD="+currentDir)

		err = c.Run()
		if err != nil {
			return "", err
		}
		return out.String(), nil
	case TemplateCommand:
		// The file is a glazed TemplateCommand, execute it by passing the arguments
		buf := &strings.Builder{}

		parsedLayers := layers.NewParsedLayers()
		err := file.Command.Layers.ForEachE(func(_ string, l layers.ParameterLayer) error {
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
					return errors.Newf("layer %s does not accept arguments (only default layer can)", l.GetSlug())
				}

				argumentDefinitions := file.Command.GetDefaultArguments()
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

		err = file.Command.RunIntoWriter(context.Background(), parsedLayers, buf)
		if err != nil {
			return "", err
		}
		return buf.String(), nil
	case Plain:
		// The file is not executable, return its content
		b, err := os.ReadFile(path)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}

	return "", fmt.Errorf("unsupported file type: %v", file.Type)
}
