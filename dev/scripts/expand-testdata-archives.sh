#!/bin/bash

set -euo pipefail

cd "$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )/../.."

while read -r dir; do
  echo "${dir}"

  pushd "$( dirname "${dir}" )" > /dev/null

  rm -fr testdata
  mkdir testdata
  cd testdata

  tar -xzf ../testdata.tar.gz

  popd > /dev/null
done < <(
  find . -type f -name testdata.tar.gz | sort
)
