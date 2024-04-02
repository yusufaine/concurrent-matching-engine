#!/bin/bash

# This script is meant to be used as the merging criteria and not part of the active CI/CD.

# This script generates random tests using `randomtest.py` and pipes the output to 
# `./grader engine` to check if the test passes. If the test fails, the generated 
# input is saved to a file in the `tests` directory for further inspection and should
# be used in the CI/CD pipeline to ensure that the issues are fixed.

# Usage:
#   ./scripts/randomtest.sh [count]
#     count: Number of random tests to run. Default is 10.

if [ -z "$1" ]; then
  count=10
else
  count=$1
fi

for i in $(seq 1 $count); do
  timestamp=$(date +%Y%m%d%H%M%S)
  python3 scripts/randomtest.py > tests/test_${timestamp}.in
  resp=$(./grader engine < tests/test_${timestamp}.in 2>&1 | tail -1 | tr -d '\0')
  if [ "$resp" != "test passed." ]; then
    echo "Random test ${i} failed ❌"
    echo "  Summary -- ${resp}"
    echo "  Run the following for the full output: './grader engine < tests/test_${timestamp}.in'"
    exit 1
  else
    echo "Random test ${i} passed ✅"
    rm tests/test_${timestamp}.in
  fi
done

