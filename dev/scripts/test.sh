#!/bin/bash

set -euo pipefail

exitcode=0

root="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )/../.."

cd "${root}"

echo
echo "==="
echo "=== go mod"
echo "==="
echo

while read -r file
do
  dir="$( dirname "${file}" )"

  echo ""${dir}""
  pushd "${root}/${dir}" > /dev/null

  go mod tidy
  go fmt ./...

  popd > /dev/null
done < <(
  find . -name go.mod
)

echo
echo "==="
echo "=== go build"
echo "==="
echo

while read -r file
do
  dir="$( dirname "${file}" )"

  echo ""${dir}""
  pushd "${root}/${dir}" > /dev/null

  go build -o /dev/null .

  popd > /dev/null
done < <(
  find . -name main.go
)

echo
echo "==="
echo "=== go test"
echo "==="
echo

rm -fr tmp/earl-reports
find . -type d -name testoutput -exec rm -fr {} +

while read -r dir
do
  echo ""${dir}""
  pushd "${root}/${dir}" > /dev/null

  mkdir testoutput

  export TESTING_EARL_OUTPUT="testoutput/earl.ttl"
  export TESTING_DEBUG_RDFIO_OUTPUT="testoutput/rdfio.txt"

  set +e

  go test -race -shuffle=on .

  if [ $? -ne 0 ]; then
    exitcode=1
  fi

  set -e

  popd > /dev/null
done < <(
  find . -type f -name testsuite_test.go -exec dirname {} \;
)

mkdir -p tmp/earl-reports

while read -r file
do
  testpackage="$( dirname "$(dirname "${file}")" )"
  dest="${root}/tmp/earl-reports/${testpackage}.ttl"

  mkdir -p "$( dirname "${dest}" )"
  cp "${file}" "${dest}"
done < <(
  find . -path '*/testoutput/earl.ttl'
)

cd tmp

find earl-reports -type f -print0 \
  | tar -czf earl-reports.tar.gz --null -T -

echo
echo "==="
echo "=== exit ${exitcode}"
echo "==="
echo

[[ $exitcode -eq 0 ]] || exit $exitcode

exec git diff --exit-code
