#!/bin/bash

# Default values
branch="origin/main"
exclude_files=()
context_size="-U3"
include_paths=""
exclude_paths=""
exclude_package=false

# Function to display usage information
usage() {
  echo "Usage: $0 [options]"
  echo "Options:"
  echo "  -b, --branch BRANCH   Specify a branch (default: origin/main)"
  echo "  -e, --exclude FILES   Exclude specific files (comma-separated list)"
  echo "  -s, --short            Reduce diff context size to 5 lines"
  echo "  -o, --only PATHS       Include specific paths only (comma-separated list)"
  echo "  --no-package           Exclude common package manager files (go.mod, go.sum, package.json, package-lock.json, etc.)"
  exit 1
}

# Parse command-line arguments
while [[ $# -gt 0 ]]; do
  key="$1"
  case $key in
    -b|--branch)
      branch="$2"
      shift
      shift
      ;;
    -e|--exclude)
      exclude_files+=($(echo "$2" | tr ',' ' '))
      shift
      shift
      ;;
    -l|--long)
      context_size="-U10"
      shift
      ;;
    -s|--short)
      context_size="-U1"
      shift
      ;;
    -o|--only)
      include_paths="$2"
      shift
      shift
      ;;
    --no-tests)
      exclude_files+=("*.test.js" "*.test.ts" "*.spec.js" "*.spec.ts" "*_test.go")
      shift
      ;;
    --no-package)
      exclude_package=true
      shift
      ;;
    -h|--help)
      usage
      ;;
    *)
      echo "Unknown option: $1"
      usage
      ;;
  esac
done

# Exclude common package manager files
if [ "$exclude_package" = true ]; then
	exclude_files+=("go.mod" "go.sum" "package.json" "package-lock.json" "yarn.lock" "composer.json" "composer.lock" "yarn.lock", "sorbet")
fi

# Construct the exclusion patterns
exclude_patterns=""
for file in "${exclude_files[@]}"; do
  exclude_patterns+=" :!$file"
  # while we're printing, add quotes
  # exclude_patterns+=" ':!$file'"
done

# Construct the inclusion patterns
include_patterns=""
if [ -n "$include_paths" ]; then
  IFS=',' read -ra paths <<< "$include_paths"
  for path in "${paths[@]}"; do
    include_patterns+=" :$path"
  done
fi

cd "$PROMPTO_PARENT_PWD" 
# Print git diff summary with --stat
echo git diff --stat "$branch" -- . $exclude_patterns $include_patterns
git diff --stat "$branch" -- . $exclude_patterns $include_patterns

# Run git diff command with diff-filter=d to exclude deleted files
echo git diff -w "$context_size" "$branch" --diff-filter=d -- . $exclude_patterns $include_patterns
cd "$PROMPTO_PARENT_PWD" && pwd && git diff -w "$context_size" "$branch" --diff-filter=d -- . $exclude_patterns $include_patterns
