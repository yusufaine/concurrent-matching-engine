#!/bin/bash

# This script runs all tests in the "tests" directory and fast fails if any test fails.

if [[ $# -gt 1 ]]; then
  echo "Usage: $0 [thread|address]"
  exit 1
fi

# check if the argument passed is "thread" or "address"
IS_THREAD=0
IS_ADDRESS=0
if [[ $# -eq 0 ]]; then
  echo "Building without any sanitizer..."
else
  if [[ $1 == "thread" ]]; then
    echo "Building with thread sanitizer..."
    IS_THREAD=1
  elif [[ $1 == "address" ]]; then
    echo "Building with address sanitizer..."
    IS_ADDRESS=1
  else
    echo "Usage: $0 [thread|address]"
    exit 1
  fi
fi

if [[ -f Makefile ]]; then
  if [[ $IS_THREAD -eq 1 ]]; then
    make thread > /dev/null
  elif [[ $IS_ADDRESS -eq 1 ]]; then
    make address > /dev/null
  else
    make all > /dev/null
  fi
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
  out=$(./grader engine < "$file" 2>&1 | tail -1 | tr -d '\0')
  if [ "$out" != "test passed." ]; then
    echo "Test for '$file' ❌:"
    echo "  -- $out"
    echo "  -- Run the following for the full output: './grader engine < $file'"
    exit 1
  else 
    echo "Test case '$file' ✅"
  fi
done

echo "Finished processing all files."
