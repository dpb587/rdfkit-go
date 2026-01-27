#!/bin/bash

set -euo pipefail

cd "$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

rm -fr testdata

git clone --depth=1 --no-checkout https://github.com/rdfa/rdfa.github.io.git testdata

cd testdata

git reset origin/HEAD

git archive HEAD | tar -xzf- $( git ls-tree --name-only -r HEAD | grep \
  -e 'LICENSE' \
  -e '^contexts/' \
  -e '^test-suite/' \
  -e '^vocabs/'
)

cat > .git.sh <<EOF
git init
git remote add origin $( git remote get-url origin )
git fetch origin
git reset $( git rev-parse HEAD )
EOF

HEAD_TIME=$( git log -1 --format=%ci HEAD )

rm -fr .git

gtar -cf- \
  --sort=name \
  --format=posix \
  --pax-option='exthdr.name=%d/PaxHeaders/%f' \
  --pax-option='delete=atime,delete=ctime' \
  --clamp-mtime \
  --mtime="${HEAD_TIME}" \
  --numeric-owner \
  --owner=0 \
  --group=0 \
  --mode='go+u,go-w' \
  . \
  | gzip \
    --no-name \
    --best \
    > ../testdata.tar.gz

cd ..

rm -fr testdata
