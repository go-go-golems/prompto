---
Title: Getting Started with Prompto
Slug: getting-started
Short: An overview of Prompto, its installation, and basic usage
Topics:
- introduction
- installation
- usage
Commands:
- prompto config repositories add
- prompto config repositories remove
- prompto config repositories get
- prompto list
- prompto get
- prompto serve
Flags:
IsTopLevel: true
IsTemplate: false
ShowPerDefault: true
SectionType: GeneralTopic
---

# Getting Started with Prompto

## Overview

`prompto` is a command-line tool designed to generate context for prompting Large Language Models (LLMs). It bridges the gap between the limited context size of LLMs and the necessity for substantial API and documentation during code generation.

## Why Prompto?

Prompting LLMs with an exhaustive context derived from API and other technical documents aids in generating precise code. Manually copying and pasting the appropriate context from source files can be laborious. Though tools like [oak](https://github.com/go-go-golems/oak) facilitate by using tree-sitter queries on codebases, there's often a need to run diverse sets of queries on distinct parts of the codebase. Using custom shell scripts to compile this information became a cumbersome chore, especially when extracting contexts from multiple repositories. Enter `prompto`, crafted to alleviate these concerns.

## How Prompto Works

`prompto` functions by scanning a predefined list of repositories mentioned in `~/.prompto/config.yaml`. For each repository, the tool searches for a `prompto/` directory. All files within this directory are treated as prompts, retrievable using the `prompto get` command. In cases where the file is executable, the `prompto get` command will run it. Consequently, this allows users to present both static files (documentation, example data) and dynamic prompts (live commands).

**Note:** As `prompto` doesn't intrinsically verify the safety of the commands it executes, users are strongly cautioned against employing `prompto` on repositories not under their direct control.

## Installation

You can install Prompto using various package managers:

### Homebrew

```bash
brew tap go-go-golems/go-go-go
brew install go-go-golems/go-go-go/prompto
```

### apt-get

```bash
echo "deb [trusted=yes] https://apt.fury.io/go-go-golems/ /" >> /etc/apt/sources.list.d/fury.list
apt-get update
apt-get install prompto
```

### yum

```bash
echo "
[fury]
name=Gemfury Private Repo
baseurl=https://yum.fury.io/go-go-golems/
enabled=1
gpgcheck=0
" >> /etc/yum.repos.d/fury.repo
yum install prompto
```

### Go

```bash
go get -u github.com/go-go-golems/prompto/cmd/prompto
```

You can also install by downloading the binaries straight from [GitHub](https://github.com/go-go-golems/prompto/releases).

## Basic Usage

### Configuring Repositories

To configure the repositories for scanning, use the following commands:

```bash
# Add repositories
prompto config repositories add /path/to/repo1 /path/to/repo2

# Remove repositories
prompto config repositories remove /path/to/repo1

# List configured repositories
prompto config repositories get
```

These commands will update the `~/.prompto/config.yaml` file automatically.

### Creating Prompts in a Repository

For `prompto` to recognize prompts within a repository, create a `prompto/` directory at the root of the repository. Place any desired files (documentation, example data) or executable scripts within this directory.

### Listing Available Prompts

To view all available prompts, use:

```bash
prompto list
```

### Getting a Prompt Context

To retrieve the context from a specific prompt, use:

```bash
prompto get [prompt-name]
```

### Web Interface

Prompto also provides a web interface for easier access to your prompts. To start the web server, use:

```bash
prompto serve
```

This will start a local web server, typically on `http://localhost:8080`.

#### Using with Cursor IDE

Run `prompto help cursor-setup` to get started.

## Safety Considerations

Always ensure that repositories added to `prompto` are safe and trusted. As `prompto` doesn't inherently verify the safety of executed commands, it's vital to be cautious and avoid using repositories that may contain malicious content.

## Examples

Here are some examples of using Prompto:

```bash
# List available prompts
❯ prompto list
/home/manuel/code/wesen/corporate-headquarters/glazed glazed/definitions
/home/manuel/code/wesen/corporate-headquarters/common-sense cms/data/example-flags.md
/home/manuel/code/wesen/corporate-headquarters/common-sense cms/data/form-dsl.yaml
/home/manuel/code/wesen/corporate-headquarters/common-sense cms/data/glazed-types.md

# Retrieve specific contexts
❯ prompto get cms/sql
// Schema is a struct that represents the schema of a CMS object.
// This contains all the necessary tables, as well as the main table for the object
// to which all other tables are joined on its `id` using the `parent_id` field.
//
// A CMS object is represents by a main table and multiple additional tables used
...

❯ prompto get glazed/definitions
File: pkg/cmds/cmds.go
// CommandDescription contains the necessary information for registering
// a command with cobra. Because a command gets registered in a verb tree,
// a full list of Parents all the way to the root needs to be provided.
type CommandDescription struct {
...

❯ prompto get cms/data/plant-dsl.yaml
# This is the main table that represents the plant object
#
#plant:
  help: This table represents a plant object in the CMS
  fields:
...
```

With this guide, you should be ready to start using Prompto to enhance your LLM prompting workflow!