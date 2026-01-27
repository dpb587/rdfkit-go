#!/bin/bash

set -euo pipefail

assets_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )/../../tmp/build-assets"

echolog () { echo $( date -u +[%Y-%m-%dT%H:%M:%SZ] ) "$@" ; }

assets_build_tag=$( git tag --points-at HEAD | grep . || echo "dev" )
assets_build_commit=$( git rev-parse HEAD | cut -c-10 )
assets_build_clean=true
assets_build_time=$( date -u +%Y-%m-%dT%H:%M:%SZ )
assets_version=$( sed 's/^v//' <<< "${assets_build_tag}" )

if [[ $( git clean -dnx | wc -l ) -gt 0 ]] ; then
  if [[ "${assets_build_tag}" != "dev" ]]; then
    echo "ERROR: building an official version requires a clean repository"
    git clean -dnx

    exit 1
  fi

  assets_build_clean=false
  assets_version=dev
fi

echolog "properties/build tag=${assets_build_tag} commit=${assets_build_commit} clean=${assets_build_clean} time=${assets_build_time}"

export CGO_ENABLED=0

function assets_build_go () {
  bundle="${1}" ; shift
  package="${1}" ; shift
  output_name="${1}" ; shift

  pushd "${package}" > /dev/null

  echolog "build-go bundle=${bundle} package=${package} version=${assets_version}"

  for target in "${@}"
  do
    target_os="$( cut -d- -f1 <<< "${target}" )"
    target_arch="$( cut -d- -f2 <<< "${target}" )"

    echolog "properties/runtime os=${target_os} arch=${target_arch}"

    bundle_dir="${assets_dir}/${bundle}-${target_os}-${target_arch}"
    mkdir -p "${bundle_dir}"

    output_file="${output_name}"

    if [ "${target_os}" == "windows" ]
    then
      output_file="${output_file}.exe"
    fi

    GOOS="${target_os}" GOARCH="${target_arch}" go build \
      -ldflags "
        -s -w
        -X main.Version=${assets_version}
        -X main.BuildTag=${assets_build_tag}
        -X main.BuildCommit=${assets_build_commit}
        -X main.BuildClean=${assets_build_clean}
        -X main.BuildTime=${assets_build_time}
      " \
      -o "${bundle_dir}/${output_file}" \
      .
  done

  popd > /dev/null
}

function assets_package_zip () {
  echolog "package-zip name=${2}"

  pushd "${assets_dir}" > /dev/null

  zip -9r "${assets_dir}/${2}.zip" "${1}"
  rm -fr "${1}"

  popd > /dev/null
}

function assets_package_tar_gz () {
  echolog "package-tar-gz name=${2}"

  pushd "${assets_dir}" > /dev/null

  tar -vcf "${assets_dir}/${2}.tar" "${1}"
  rm -fr "${1}"

  gzip -9 "${assets_dir}/${2}.tar"

  popd > /dev/null
}
