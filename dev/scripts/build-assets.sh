#!/bin/bash

set -euo pipefail

cd "$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )/../.."

source dev/scripts/.assets.sh

rm -fr "${assets_dir}"

assets_build_go rdfkit cmd/rdfkit rdfkit darwin-amd64 darwin-arm64 linux-amd64 linux-arm64 windows-amd64

assets_package_zip rdfkit-darwin-amd64 "rdfkit-${assets_version}-darwin-amd64"
assets_package_zip rdfkit-darwin-arm64 "rdfkit-${assets_version}-darwin-arm64"
assets_package_tar_gz rdfkit-linux-arm64 "rdfkit-${assets_version}-linux-arm64"
assets_package_tar_gz rdfkit-linux-amd64 "rdfkit-${assets_version}-linux-amd64"
assets_package_zip rdfkit-windows-amd64 "rdfkit-${assets_version}-windows-amd64"
