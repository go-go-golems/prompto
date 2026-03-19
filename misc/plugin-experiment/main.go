package main

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/go-go-golems/glazed/pkg/cmds/fields"
	"github.com/go-go-golems/glazed/pkg/cmds/schema"
)

func main() {
	// Execute the shell script and capture the YAML output
	cmd := exec.Command("bash", "-c", "./misc/plugin-experiment/experiment.sh")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error executing shell script: %v\n", err)
		return
	}

	// Read the YAML output
	yamlOutput := out.Bytes()

	fmt.Printf("YAML Output:\n%s\n", yamlOutput)
	fmt.Println()

	// Create a section from the YAML output
	section, err := schema.NewSectionFromYAML(yamlOutput)
	if err != nil {
		fmt.Printf("Error creating section: %v\n", err)
		return
	}

	// Access the section properties
	fmt.Printf("Slug: %s\n", section.Slug)
	fmt.Printf("Name: %s\n", section.Name)
	fmt.Printf("Description: %s\n", section.Description)

	definitions := section.GetDefinitions()

	fmt.Println("Flags:")
	definitions.ForEach(func(def *fields.Definition) {
		if !def.IsArgument {
			fmt.Printf("  - Name: %s\n", def.Name)
			fmt.Printf("    Type: %s\n", def.Type)
			fmt.Printf("    Help: %s\n", def.Help)
			if def.Default != nil {
				fmt.Printf("    Default: %v\n", *def.Default)
			}
		}
	})

	fmt.Println("Arguments:")
	definitions.ForEach(func(def *fields.Definition) {
		if def.IsArgument {
			fmt.Printf("  - Name: %s\n", def.Name)
			fmt.Printf("    Type: %s\n", def.Type)
			fmt.Printf("    Help: %s\n", def.Help)
			if def.Default != nil {
				fmt.Printf("    Default: %v\n", *def.Default)
			}
		}
	})
}
