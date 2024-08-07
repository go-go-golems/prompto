#!/bin/bash


# Define the commands and their attributes
commands=(
  "default|Default Command|The default command with no verb.|default_flags|default_args"
  "fields|Fields Command|A command to configure fields.|fields_flags|fields_args"
  "filters|Filters Command|A command to configure filters.|filters_flags|filters_args"
)

# Define the flag lists for each command
default_flags=(
  "sort-columns|bool|Sort columns alphabetically|false"
)

fields_flags=(
  "fields|stringList|Fields to include in the output|all"
  "sort-columns|bool|Sort columns alphabetically|false"
  "reorder-columns|stringList|List of columns to reorder|"
)

filters_flags=(
  "filter|stringList|Fields to remove from output|"
  "remove-nulls|bool|Remove columns with null values|false"
  "remove-duplicates|stringList|List of columns to consider for duplicate removal (requires results to be sorted)|"
)

# Define the argument lists for each command
default_args=(
  "input_file|string|Path to the input file"
)

fields_args=(
  "input_file|string|Path to the input file"
  "output_file|string|Path to the output file"
)

filters_args=(
  "input_file|string|Path to the input file"
  "output_file|string|Path to the output file"
)

source "$(dirname "${BASH_SOURCE[0]}")/flaglib.sh"

# Output commands and their flags as YAML
output_yaml