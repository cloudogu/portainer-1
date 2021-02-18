#!/usr/bin/env bash

PLATFORM=$1
ARCH=$2

export GOPATH="/tmp/go"

binary="portainer"

mkdir -p dist
mkdir -p ${GOPATH}/src/github.com/cloudogu/portainer-ce

cp -R api ${GOPATH}/src/github.com/cloudogu/portainer-ce/api

cd 'api/cmd/portainer'

go get -t -d -v ./...
GOOS=${PLATFORM} GOARCH=${ARCH} CGO_ENABLED=0 go build -a --installsuffix cgo --ldflags '-s'

if [ "${PLATFORM}" == 'windows' ]; then
  mv "$BUILD_SOURCESDIRECTORY/api/cmd/portainer/${binary}.exe" "$BUILD_SOURCESDIRECTORY/dist/portainer.exe"
else  
  mv "$BUILD_SOURCESDIRECTORY/api/cmd/portainer/$binary" "$BUILD_SOURCESDIRECTORY/dist/portainer"
fi
