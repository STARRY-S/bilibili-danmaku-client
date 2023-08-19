#!/usr/bin/env bash

set -euo pipefail

cd $(dirname $0)/../
WORKINGDIR=$(pwd)

ARCH=('amd64')
OS=('windows' 'linux')
BUILD_FLAG="-s -w"

GITCOMMIT=$(git rev-parse HEAD || echo "")
if [[ -z "${GITCOMMIT}" ]]; then
    if [[ -n "${GITHUB_SHA:-}" ]]; then
        echo "GITHUB_SHA: ${GITHUB_SHA}"
        GITCOMMIT=${GITHUB_SHA}
    else
        GITCOMMIT="UNKNOW"
    fi
fi
VERSION=$(git describe --tags 2>/dev/null || echo "")
if [[ -z "${VERSION}" ]]; then
    if [ "${GITHUB_REF_TYPE:-}" = "tag" ] && [ -n "${GITHUB_REF:-}" ]; then
        echo "GITHUB_REF: ${GITHUB_REF}"
        VERSION=${GITHUB_REF#refs/tags/}
    else
        VERSION="HEAD-${GITCOMMIT:0:8}"
    fi
fi
echo "Build version: ${VERSION}"

mkdir -p $WORKINGDIR/build
cd $WORKINGDIR/build

for os in ${OS[@]}
do
    for arch in ${ARCH[@]}
    do
        OUTPUT="bilibili-danmaku-client-$os-$arch-$VERSION"
        if [[ $os = "windows" ]]; then
            OUTPUT=$OUTPUT.exe
        fi
        CGO_ENABLED=1 GOOS=$os GOARCH=$arch go build -ldflags "${BUILD_FLAG}" -o $OUTPUT ..
        echo $(pwd)/$OUTPUT
    done
done
