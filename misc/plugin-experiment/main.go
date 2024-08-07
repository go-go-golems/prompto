package main

import (
	"bytes"
	"fmt"
	"github.com/go-go-golems/glazed/pkg/cmds/layers"
	"github.com/go-go-golems/glazed/pkg/cmds/parameters"
	"os/exec"
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

	// Create a parameter layer from the YAML output
	parameterLayer, err := layers.NewParameterLayerFromYAML(yamlOutput)
	if err != nil {
		fmt.Printf("Error creating parameter layer: %v\n", err)
		return
	}

	// Access the parameter layer properties
	fmt.Printf("Slug: %s\n", parameterLayer.Slug)
	fmt.Printf("Name: %s\n", parameterLayer.Name)
	fmt.Printf("Description: %s\n", parameterLayer.Description)

	fmt.Println("Flags:")
	parameterLayer.ParameterDefinitions.ForEach(func(flag *parameters.ParameterDefinition) {
		if !flag.IsArgument {
			fmt.Printf("  - Name: %s\n", flag.Name)
			fmt.Printf("    Type: %s\n", flag.Type)
			fmt.Printf("    Help: %s\n", flag.Help)
			if flag.Default != nil {
				fmt.Printf("    Default: %v\n", *flag.Default)
			}
		}
	})

	fmt.Println("Arguments:")
	parameterLayer.ParameterDefinitions.ForEach(func(flag *parameters.ParameterDefinition) {
		if flag.IsArgument {
			fmt.Printf("  - Name: %s\n", flag.Name)
			fmt.Printf("    Type: %s\n", flag.Type)
			fmt.Printf("    Help: %s\n", flag.Help)
			if flag.Default != nil {
				fmt.Printf("    Default: %v\n", *flag.Default)
			}
		}
	})
}
