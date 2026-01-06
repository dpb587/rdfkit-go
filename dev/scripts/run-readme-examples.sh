#!/bin/bash

set -euo pipefail

cd "$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )/../../examples"

go mod tidy

go run ./html-extract https://microsoft.com \
  > html-extract/readme-output.ttl

cd ../cmd/rdfkit

go run . --help \
  > ../../dev/artifacts/readme-rdfkit-help.txt 2>&1

go run . pipe --help \
  > ../../dev/artifacts/readme-rdfkit-pipe-help.txt 2>&1


go run . export-dot \
  -i https://www.w3.org/2000/01/rdf-schema.ttl \
  | tee ../../dev/artifacts/readme-rdf-ontology.dot \
  | dot -Tsvg \
  > ../../dev/artifacts/readme-rdf-ontology.svg
