#!/bin/bash

set -euo pipefail

cd "$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

curl -Lo testdata.tar.gz https://www.w3.org/2013/TrigTests/TESTS.tar.gz
