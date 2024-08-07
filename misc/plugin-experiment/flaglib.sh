#!/bin/bash

# Function to output commands, their flags, and arguments as YAML
output_yaml() {
  echo "commands:"
  for cmd in "${commands[@]}"; do
    IFS='|' read -r slug name description flags_var args_var <<< "$cmd"
    echo "  - slug: $slug"
    echo "    name: $name"
    echo "    description: $description"

    if [[ -n "$flags_var" ]]; then
      echo "    flags:"
      flags_array_name="${flags_var}"
      for flag in "${!flags_array_name[@]}"; do
        IFS='|' read -r flag_name type help default <<< "$flag"
        echo "      - name: $flag_name"
        echo "        type: $type"
        echo "        help: $help"
        if [[ -n "$default" ]]; then
          if [[ "$type" == *"List" ]]; then
            IFS=',' read -r -a default_array <<< "$default"
            echo "        default:"
            for item in "${default_array[@]}"; do
              echo "          - $item"
            done
          else
            echo "        default: $default"
          fi
        else
          echo "        default: []"
        fi
      done
    fi

    if [[ -n "$args_var" ]]; then
      echo "    arguments:"
      args_array_name="${args_var}"
      for arg in "${!args_array_name[@]}"; do
        IFS='|' read -r arg_name type help <<< "$arg"
        echo "      - name: $arg_name"
        echo "        type: $type"
        echo "        help: $help"
      done
    fi
  done
}