#!/bin/bash

set -euo pipefail

cd examples

go mod tidy

go run ./rdf-to-dot -i https://www.w3.org/2000/01/rdf-schema.ttl \
  | tee rdf-to-dot/readme-output.dot \
  | dot -Tsvg \
  > rdf-to-dot/readme-output.svg

go run ./html-extract https://microsoft.com \
  > html-extract/readme-output.ttl
