#!/bin/bash -e

# Computes the overall test coverage of the project.
COVERAGE_THRESHOLD=0
COVFILE=build/coverage.out
COVHTMLFILE=build/coverage.html
MODE=count

# Ensure build/ directory is present
mkdir -p build/

# Ensure mocks are generated
mockery --all --keeptree

# Build the coverage profile
go test -v -coverpkg=./db/...,./grpc/...,./metadata/...,./zeus/...,./config/...,./util/...,./authz/... -covermode=$MODE -coverprofile=${COVFILE} ./...

# Get coverage report in HTML
go tool cover -html=${COVFILE} -o ${COVHTMLFILE}

# Output the percent
COVERAGE_PERCENT=$(go tool cover -func=${COVFILE} | \
  tail -n 1 | \
  awk '{ print $3 }' | \
  sed -e 's/^\([0-9]*\).*$/\1/g')

echo "Coverage of the project is $COVERAGE_PERCENT"
echo "Coverage Threshold of the project is $COVERAGE_THRESHOLD"

if [ "$COVERAGE_PERCENT" -lt $COVERAGE_THRESHOLD ] ; then
  echo "Coverage check failed. Please add more test cases"
  exit 1
else
  echo "Coverage check passed"
  exit 0
fi
