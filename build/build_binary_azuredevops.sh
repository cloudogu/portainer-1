#!/usr/bin/env bash
set -x

PLATFORM=$1
ARCH=$2

export GOPATH="/tmp/go"

binary="portainer"

mkdir -p dist
mkdir -p ${GOPATH}/src/github.com/portainer/portainer

cp -R api ${GOPATH}/src/github.com/cloudogu/portainer-ce/api

cp -r "./mustache-templates" "./dist"

cd 'api/cmd/portainer'

go get -t -d -v ./...
GOOS=${PLATFORM} GOARCH=${ARCH} CGO_ENABLED=0 go build -a -trimpath --installsuffix cgo --ldflags "-s \
-X 'github.com/cloudogu/portainer-ce/api/build.BuildNumber=${BUILDNUMBER}' \
-X 'github.com/cloudogu/portainer-ce/api/build.ImageTag=${CONTAINER_IMAGE_TAG}' \
-X 'github.com/cloudogu/portainer-ce/api/build.NodejsVersion=${NODE_VERSION}' \
-X 'github.com/cloudogu/portainer-ce/api/build.YarnVersion=${YARN_VERSION}' \
-X 'github.com/cloudogu/portainer-ce/api/build.WebpackVersion=${WEBPACK_VERSION}' \
-X 'github.com/cloudogu/portainer-ce/api/build.GoVersion=${GO_VERSION}'"

if [ "${PLATFORM}" == 'windows' ]; then
  mv "$BUILD_SOURCESDIRECTORY/api/cmd/portainer/${binary}.exe" "$BUILD_SOURCESDIRECTORY/dist/portainer.exe"
else
  mv "$BUILD_SOURCESDIRECTORY/api/cmd/portainer/$binary" "$BUILD_SOURCESDIRECTORY/dist/portainer"
fi

