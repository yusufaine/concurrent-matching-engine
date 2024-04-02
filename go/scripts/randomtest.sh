#!/bin/bash

# This script is meant to be used as the merging criteria and not part of the active CI/CD.

# This script generates random tests using `randomtest.go` and pipes the output to 
# `./grader engine` to check if the test passes. If the test fails, the generated 
# input is saved to a file in the `tests` directory for further inspection and should
# be used in the CI/CD pipeline to ensure that the issues are fixed.

# Usage:
#   ./scripts/randomtest.sh [count]
#     count: Number of random tests to run. Default is 10.

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


GRADER_BIN="./grader"

# if machine is arm64, use the arm64 binary
if [ "$(uname -m)" == "arm64" ]; then
  GRADER_BIN="./grader_arm64-apple-darwin"
fi

if [ -z "$1" ]; then
  count=10
else
  count=$1
fi

for i in $(seq 1 $count); do
  timestamp=$(date +%Y%m%d%H%M%S)
  go run ./scripts/randomtest.go > tests/test_${timestamp}.in
  resp=$($GRADER_BIN engine < tests/test_${timestamp}.in 2>&1 | tail -1 | tr -d '\0')
  if [ "$resp" != "test passed." ]; then
    echo "Random test ${i} failed ❌"
    echo "  Summary -- ${resp}"
    echo "  Run the following for the full output: '$GRADER_BIN engine < tests/test_${timestamp}.in'"
    exit 1
  else
    echo "Random test ${i} passed ✅"
    rm tests/test_${timestamp}.in
  fi
done

