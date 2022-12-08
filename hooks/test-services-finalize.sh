#!/bin/bash -e

# Computes the overall test coverage of the project.
COVERAGE_THRESHOLD=0
COVFILE=coverage.out
MODE=count

# ensure mocks are generated
mockery --all --keeptree
# Build the coverage profile
go test -v -coverpkg=./db/...,./grpc/... -covermode=$MODE -coverprofile=${COVFILE} ./...

# get coverage report in HTML
#go tool cover -html=${COVFILE} -o coverage.html

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