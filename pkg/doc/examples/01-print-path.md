---
Title: Using the --print-path Flag in Prompto
Slug: using-print-path-flag
Short: |
  The --print-path flag in Prompto allows users to print the path of a prompt.
  This can be useful when working with various GO GO GOLEMS commands.
Topics:
- Command Line Usage
- Flags
Commands:
- prompto get
- prompto list
Flags:
- --print-path
IsTopLevel: false
ShowPerDefault: false
SectionType: CommandUsage
---

Prompto is a command line tool that generates prompts from a list of repositories.
One of its features is the `--print-path` flag, which allows users to print the path of a prompt.
This can be particularly useful when working with various GO GO GOLEMS commands which accept a list of files.

## Usage

The `--print-path` flag can be used with the `get` command in Prompto. Here is an example of its usage:

```bash
prompto get prompt --print-path -- [flags and arguments]
```

In this command, `prompt` is the name of the prompt you want to get, 
and `[flags and arguments]` represents any additional arguments you want to pass to the command.

## How it Works

When the `--print-path` flag is used, Prompto will print the path of the prompt.
If the file type is `Plain`, it will simply print the path.

If the file type is not `Plain`, it will render the file with the provided arguments,
create a temporary file with the rendered content,
and print the path of this temporary file. In that case, only 20 previously created temporary files will be kept,
the rest will be removed

## Example

Here is an example of using the `--print-path` flag with the `get` command:

```bash
prompto get my-prompt --print-path
```

This command will print the path of the `my-prompt` prompt.
