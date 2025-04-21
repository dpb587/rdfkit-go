#!/bin/bash

set -euo pipefail

cd "$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

rm -fr testdata/

mkdir testdata/

wget -O testdata/manifest.jsonld https://w3c.github.io/json-ld-api/tests/toRdf-manifest.jsonld

iter() {
  while read -r p; do
    mkdir -p "$( dirname "testdata/${p}" )"
    wget -O "testdata/${p}" "https://w3c.github.io/json-ld-api/tests/${p}"
  done 
}

iter < <(
  jq -r '.sequence[].expect | select(.)' testdata/manifest.jsonld
  jq -r '.sequence[].input | select(.)' testdata/manifest.jsonld
  jq -r '.sequence[].option.expandContext | select(.)' < testdata/manifest.jsonld
)

iter < <(
  for d in expand toRdf ; do
    find ./testdata/$d -name '*-in.jsonld' -exec cat {} \; \
      | jq -r \
        '..
          | select(. | type == "object") 
          | [ .["@context"], .["@import"] ][]
          | select(.) 
          | if type != "array" then [.] else . end 
          | map(select(type == "string"))[]
        ' \
      | grep '\.jsonld' \
      | sed "s#^#$d/#" \
      || true
  done
)

iter < <(
  for d in expand toRdf ; do
    find ./testdata/$d -name '*-context*.jsonld' -exec cat {} \; \
      | jq -r \
        '..
          | select(. | type == "object") 
          | [ .["@context"], .["@import"] ][]
          | select(.) 
          | if type != "array" then [.] else . end 
          | map(select(type == "string"))[]
        ' \
      | grep '\.jsonld' \
      | sed "s#^#$d/#" \
      | sed 's#/../#/#' \
      || true
  done
)

cd testdata/

GZIP=-9 tar -czf ../testdata.tar.gz ./

cd ../

rm -fr testdata/
