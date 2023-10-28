package pkg

import (
	"bytes"
	"context"
	"fmt"
	"github.com/go-go-golems/glazed/pkg/cmds/layers"
	"github.com/go-go-golems/glazed/pkg/cmds/parameters"
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
		parsedLayers := map[string]*layers.ParsedParameterLayer{}
		ps, args_, err := parameters.GatherFlagsFromStringList(
			restArgs, file.Command.Flags,
			false, false,
			"")
		if err != nil {
			return "", err
		}
		arguments, err := parameters.GatherArguments(args_, file.Command.Arguments, false, false)
		if err != nil {
			return "", err
		}
		for p := arguments.Oldest(); p != nil; p = p.Next() {
			k, v := p.Key, p.Value
			ps[k] = v
		}
		err = file.Command.RunIntoWriter(context.Background(), parsedLayers, ps, buf)
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
