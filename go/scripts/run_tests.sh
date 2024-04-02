#!/bin/bash

# This script runs all tests in the "tests" directory and fast fails if any test fails.

GRADER_BIN="./grader"

# if machine is arm64, use the arm64 binary
if [ "$(uname -m)" == "arm64" ]; then
  GRADER_BIN="./grader_arm64-apple-darwin"
fi

if [[ $# -gt 1 ]]; then
  echo "Usage: $0"
  exit 1
fi

if [[ -f Makefile ]]; then
  make engine > /dev/null
  if [[ $? -ne 0 ]]; then 
    echo "An error occurred while building the project. Exiting..."
    exit 1
  fi
else
  echo "Run this script from the root of the repository (cwd contains Makefile)"
  exit 1
fi


# Set the directory containing the ".in" files
DIRECTORY="tests"

# Find all files in the directory ending with ".in"
FILES=$(find "$DIRECTORY" -maxdepth 1 -name "*.in")

# Check if any files were found
if [[ -z "$FILES" ]]; then
  echo "No files found in $DIRECTORY ending with .in"
  exit 0
fi

echo "Running tests..."
# Loop through each file and pipe it to "./grader engine"
for file in $FILES; do
  out=$($GRADER_BIN engine < "$file" 2>&1 | tail -1 | tr -d '\0')
  if [ "$out" != "test passed." ]; then
    echo "Test for '$file' ❌:"
    echo "  -- $out"
    echo "  -- Run the following for the full output: '$GRADER_BIN engine < $file'"
    exit 1
  else 
    echo "Test case '$file' ✅"
  fi
done

echo "Finished processing all files."
