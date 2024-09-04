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

## Configuration

Create a `~/.prompto/config.yaml` file to specify repositories:

```yaml
repositories:
  - /path/to/repo1
  - /path/to/repo2
```

## Core Structures

### Repository

The `Repository` struct represents a repository containing prompts:

```go
type Repository struct {
    Path     string
    Promptos []Prompto
}
```

### Prompto

The `Prompto` struct represents an individual prompt:

```go
type Prompto struct {
    Name       string
    Group      string
    Type       FileType
    Command    *cmds.TemplateCommand
    FilePath   string
    Repository string
}
```

## Core Functions

### Repository Methods

1. `NewRepository(path string) *Repository`
   - Creates a new Repository instance.

2. `(r *Repository) LoadPromptos() error`
   - Loads all prompts from the repository.

3. `(r *Repository) Refresh() error`
   - Reloads all prompts from the repository.

4. `(r *Repository) GroupPromptos() map[string][]Prompto`
   - Groups prompts by their group name.

5. `(r *Repository) GetGroups() []string`
   - Returns a sorted list of all group names in the repository.

6. `(r *Repository) GetPromptosByGroup(group string) []Prompto`
   - Returns a sorted list of prompts for a specific group.

7. `(r *Repository) Watch(ctx context.Context, options ...watcher.Option) error`
   - Sets up a file watcher for the repository to automatically reload prompts on changes.

### Prompto Methods

1. `(p *Prompto) Render(repo string, restArgs []string) (string, error)`
   - Renders the prompt with the given arguments.

### Utility Functions

1. `LoadTemplateCommand(path string) (*cmds.TemplateCommand, bool)`
   - Loads a template command from a YAML file.

## Server Usage

Prompto includes a server component for serving prompts via HTTP:

```go
func Serve(port int, watching bool) error
```

This function sets up an HTTP server with the following endpoints:

- `/`: Root handler
- `/prompts/`: Prompt handler
- `/search`: Search handler
- `/refresh`: Refresh handler
- `/repositories`: Repositories handler

The server uses a `ServerState` struct to manage the state of repositories and prompts:

```go
type ServerState struct {
    Repositories []string
    Repos        map[string]*pkg.Repository
    mu           sync.RWMutex
    Watching     bool
}
```

### ServerState Methods

1. `NewServerState(watching bool) *ServerState`
   - Creates a new ServerState instance.

2. `(s *ServerState) LoadRepositories() error`
   - Loads all repositories specified in the configuration.

3. `(s *ServerState) CreateTemplateWithFuncs(name, tmpl string) (*template.Template, error)`
   - Creates an HTML template with custom functions for rendering prompts and repositories.

4. `(s *ServerState) GetAllRepositories() []string`
   - Returns a list of all repository paths.

5. `(s *ServerState) GetAllPromptos() []pkg.Prompto`
   - Returns a sorted list of all prompts across all repositories.

6. `(s *ServerState) GetAllGroups() []string`
   - Returns a sorted list of all unique group names across all repositories.

7. `(s *ServerState) GetPromptosByGroup(group string) []pkg.Prompto`
   - Returns a sorted list of prompts for a specific group across all repositories.

8. `(s *ServerState) GetPromptosByRepository(repo string) []pkg.Prompto`
   - Returns all prompts for a specific repository.

9. `(s *ServerState) GetGroupsByRepository(repo string) []string`
   - Returns all groups for a specific repository.

10. `(s *ServerState) GetPromptosForRepositoryAndGroup(repo, group string) []pkg.Prompto`
    - Returns all prompts for a specific repository and group.

11. `(s *ServerState) WatchRepositories(ctx context.Context) error`
    - Sets up file watchers for all repositories if watching is enabled.

## Example Usage

Here's an example of how to use Prompto to render a prompt:

```go
import (
    "fmt"
    "github.com/go-go-golems/prompto/pkg"
)

func main() {
    repo := pkg.NewRepository("/path/to/repository")
    err := repo.LoadPromptos()
    if err != nil {
        fmt.Printf("Error loading prompts: %v\n", err)
        return
    }

    promptos := repo.GetPromptosByGroup("example-group")
    if len(promptos) > 0 {
        result, err := promptos[0].Render(repo.Path, []string{"arg1", "arg2"})
        if err != nil {
            fmt.Printf("Error rendering prompt: %v\n", err)
            return
        }
        fmt.Println(result)
    }
}
```

This example demonstrates how to:
1. Create a new Repository
2. Load prompts from the repository
3. Get prompts for a specific group
4. Render a prompt with arguments

For server usage, you can start the Prompto server like this:

```go
package main

import (
    "fmt"
    "github.com/go-go-golems/prompto/pkg/server"
)

func main() {
    port := 8080
    watching := true

    err := server.Serve(port, watching)
    if err != nil {
        fmt.Printf("Error starting server: %v\n", err)
    }
}
```

This will start a server on port 8080, watching for changes in the repositories if `watching` is set to `true`.