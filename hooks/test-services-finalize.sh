#!/bin/bash -e

# Computes the overall test coverage of the project.
COVERAGE_THRESHOLD=0
COVFILE=build/coverage.out
COVHTMLFILE=build/coverage.html
MODE=count

# Needed to make gotestsum work on circleci
export PATH=$PATH:/home/circleci/project/.workspace/build-tools/golang-go1.19/go/bin
export PATH=$PATH:/home/circleci/go/bin

# Ensure build/ directory is present
mkdir -p build/

# Ensure mocks are generated
echo "--------------------------------"
printf "Generating Mocks\n"
echo "--------------------------------"
mockery --all --keeptree

# Build the coverage profile
echo "--------------------------------"
printf "Running Unit Tests\n"
echo "--------------------------------"
gotestsum --format-hide-empty-pkg --packages="./..." -- -covermode=${MODE} -coverprofile=${COVFILE}

# Get coverage report in HTML
echo "--------------------------------"
printf "Generate HTML Coverage Report\n"
echo "--------------------------------"
go tool cover -html=${COVFILE} -o ${COVHTMLFILE}

# Output the percent
COVERAGE_PERCENT=$(go tool cover -func=${COVFILE} | \
  tail -n 1 | \
  awk '{ print $3 }' | \
  sed -e 's/^\([0-9]*\).*$/\1/g')

echo "Coverage of the project is $COVERAGE_PERCENT%"
echo "Coverage Threshold of the project is $COVERAGE_THRESHOLD%"

if [ "$COVERAGE_PERCENT" -lt $COVERAGE_THRESHOLD ] ; then
  echo "Coverage check failed. Please add more test cases"
  exit 1
else
  echo "Coverage check passed"
  exit 0
fi
