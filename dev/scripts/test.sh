#!/bin/bash

set -euo pipefail

root="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )/../.."

cd "${root}"

echo
echo "==="
echo "=== go modules"
echo "==="
echo

while read -r file
do
  dir="$( dirname "${file}" )"

  echo
  echo "=== ${dir}"
  echo

  pushd "${root}/${dir}" > /dev/null

  go mod tidy
  go fmt ./...
  go test -count 1 -race -shuffle=on ./...

  popd > /dev/null
done < <(
  find . -name go.mod | sort
)

echo
echo "==="
echo "=== go build"
echo "==="
echo

while read -r file
do
  dir="$( dirname "${file}" )"

  echo
  echo "=== ${dir}"
  echo

  pushd "${root}/${dir}" > /dev/null

  go build -o /dev/null .

  popd > /dev/null
done < <(
  find . -name main.go | sort
)

echo
echo "==="
echo "=== testsuites"
echo "==="
echo

rm -fr tmp/earl-reports
find . -type d -name testoutput -exec rm -fr {} +

while read -r dir
do
  echo
  echo "=== ${dir}"
  echo

  pushd "${root}/${dir}" > /dev/null

  mkdir testoutput

  export TESTING_EARL_OUTPUT="testoutput/earl.ttl"
  export TESTING_DEBUG_RDFIO_OUTPUT="testoutput/rdfio.txt"
  export TESTING_DEV_EARL_SUMMARY_OUTPUT="testoutput/earl-summary.txt"

  go test -count 1 .

  popd > /dev/null
done < <(
  find . -type f -name testsuite_test.go -exec dirname {} \; | sort
)

mkdir -p tmp/earl-reports

while read -r file
do
  testpackage="$( dirname "$(dirname "${file}")" )"
  dest="${root}/tmp/earl-reports/${testpackage}.ttl"

  mkdir -p "$( dirname "${dest}" )"
  cp "${file}" "${dest}"
done < <(
  find . -path '*/testoutput/earl.ttl' | sort
)

cd tmp

find earl-reports -type f -print0 \
  | sort -z \
  | tar -czf earl-reports.tar.gz --null -T -

echo
echo "==="
echo "=== git diff"
echo "==="
echo

git diff --exit-code
