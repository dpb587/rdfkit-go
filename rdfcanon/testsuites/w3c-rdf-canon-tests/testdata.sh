#!/bin/bash

set -euo pipefail

cd "$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

rm -fr testdata/

git clone --depth 1 https://github.com/w3c/rdf-canon testdata/

cd testdata/

find . -mindepth 1 -maxdepth 1 \
  ! -name 'tests' \
  ! -name 'LICENSE*' \
  ! -name '.git' \
  -exec rm -rf {} +
rm -rf .git

GZIP=-9 tar -czf ../testdata.tar.gz ./

cd ../

rm -fr testdata/
