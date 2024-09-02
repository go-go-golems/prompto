---
Title: How to use prompto's API
Slug: prompto-api
Short: How to use prompto's API to build application rendering prompts.
Topics:
- prompto
- development
- api
Commands:
- GetFilesFromRepo
- LoadTemplateCommand
- RenderFile
IsTopLevel: true
IsTemplate: false
ShowPerDefault: true
SectionType: GeneralTopic
---

# Prompto: Context Generation for LLM Prompts

Prompto is a Go package designed to generate context for prompting Large Language Models (LLMs). It helps developers bridge the gap between LLMs' limited context size and the need for comprehensive API and documentation during code generation.

## Installation

Install Prompto using Go:

```bash
go get github.com/go-go-golems/prompto
```

## Usage

### Configuration

Create a `~/.prompto/config.yaml` file to specify repositories:

```yaml
repositories:
  - /path/to/repo1
  - /path/to/repo2
```

### Creating Prompts

In each repository, create a `prompto/` directory and add files or executable scripts.

## Core Functions

Prompto provides three main functions:

1. List available prompts (`GetFilesFromRepo`)
2. Load and execute prompts (`LoadTemplateCommand`)
3. Render files (`RenderFile`)

### GetFilesFromRepo

This function retrieves all prompt files from a specified repository.

```go
func GetFilesFromRepo(repo string) ([]FileInfo, error) {
    // Implementation details...
}
```

### LoadTemplateCommand

This function loads a template command from a YAML file.

```go
func LoadTemplateCommand(path string) (*cmds.TemplateCommand, bool) {
    // Implementation details...
}
```

### RenderFile

This function renders a prompt file with given arguments.

```go
func RenderFile(repo string, file FileInfo, args []string) (string, error) {
    // Implementation details...
}
```

## Example Usage

Here's an example of how to use Prompto's core functions:

```go
import (
    "fmt"
    "github.com/go-go-golems/prompto/pkg"
)

func renderPrompt(repo, promptName string, args []string) error {
    files, err := pkg.GetFilesFromRepo(repo)
    if err != nil {
        return err
    }

    for _, file := range files {
        if file.Name == promptName {
            s, err := pkg.RenderFile(repo, file, args)
            if err != nil {
                return err
            }
            fmt.Println(s)
            return nil
        }
    }

    return fmt.Errorf("prompt not found")
}
```

This example demonstrates how to:
1. Get files from a repository using `pkg.GetFilesFromRepo`
2. Find a specific prompt file
3. Render the file using `pkg.RenderFile`
4. Print the rendered content
