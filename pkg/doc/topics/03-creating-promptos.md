---
Title: Creating Promptos
Slug: creating-promptos
Short: Learn how to create and use different types of promptos in the prompto system
Topics:
- prompto
- promptos
- development
Commands:
- prompto get
Flags:
- --name
- --language
IsTopLevel: true
IsTemplate: false
ShowPerDefault: true
SectionType: Tutorial
---

# Creating Promptos

Promptos are the core components of the `prompto` system. They provide context, generate dynamic content, or create parameterized templates for Large Language Model prompts. This guide will walk you through creating and using different types of promptos.

## Types of Promptos

There are three main types of promptos:

1. Static File Prompto
2. Executable Script Prompto
3. Template Command Prompto

## Creating Promptos

All promptos should be placed in the `prompto/` directory of your repository. Ensure that your repository is added to your `prompto` configuration. 

Promptos are placed under a subdirectory of the `prompto/` directory. For example:

```
prompto/
├── test/
│   ├── foobar.md
│   └── foobar.sh
│   └── greeting.yaml
```

### 1. Static File Prompto

Static file promptos are simple text files containing information or context for your prompts.

To create a static file prompto:

1. Create a new file in your `prompto/` directory, e.g., `prompto/test/foobar.md`:

```markdown
# Foobar Documentation

This is a sample static file prompto.

- Point 1
- Point 2
- Point 3
```

2. Retrieve the prompto using:

```bash
prompto get test/foobar.md
```

### 2. Executable Script Prompto

Executable script promptos are scripts that generate dynamic content when executed.

To create an executable script prompto:

1. Create a new file in your `prompto/` directory, e.g., `prompto/test/foobar.sh`:

```bash
#!/bin/bash

echo "This is a dynamic prompto generated at $(date)"
echo "Current working directory: $PWD"
echo "Parent working directory: $PROMPTO_PARENT_PWD"
```

2. Make the script executable:

```bash
chmod +x prompto/test/foobar.sh
```

3. Retrieve the prompto using:

```bash
prompto get test/foobar.sh
```

### 3. Template Command Prompto

Template command promptos are YAML files that define a command with parameters and a template for generating content.

To create a template command prompto:

1. Create a new YAML file in your `prompto/` directory, e.g., `prompto/test/greeting.yaml`:

```yaml
name: greeting
short: Generate a greeting
parameters:
  - name: name
    type: string
    help: Name to greet
  - name: language
    type: string
    help: Language for the greeting
    default: English
template: |
  {{- if eq .language "English" -}}
  Hello, {{ .name }}!
  {{- else if eq .language "Spanish" -}}
  ¡Hola, {{ .name }}!
  {{- else if eq .language "French" -}}
  Bonjour, {{ .name }}!
  {{- else -}}
  Greeting in {{ .language }} is not supported.
  {{- end -}}
```

2. Retrieve the prompto with parameters using:

```bash
prompto get test/greeting.yaml --name Alice --language Spanish
```

## Using Promptos

You can use these promptos in your development workflow by retrieving them with the `prompto get` command:

```bash
# Get a static file prompto
prompto get test/foobar.md

# Get an executable script prompto
prompto get test/foobar.sh

# Get a template command prompto with parameters
prompto get test/greeting.yaml --name Bob --language French
```

These promptos can be used to provide context, generate dynamic content, or create parameterized templates for your Large Language Model prompts.

Remember to place all your promptos in the `prompto/` directory of your repository, and make sure the repository is added to your `prompto` configuration.

## Best Practices

1. Use descriptive names for your prompto files.
2. Organize promptos into subdirectories based on their purpose or related functionality.
3. Keep static file promptos focused and concise.
4. For executable script promptos, consider error handling and input validation.
5. For template command promptos, provide clear parameter descriptions and default values when appropriate.

By following these guidelines, you can create effective promptos that enhance your development workflow and improve your interactions with Large Language Models.