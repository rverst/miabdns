#!/usr/bin/env sh

version=$1
os=$2
arch=$3
ldflags=$4
filename=$5

cd /go/src || exit 9
dir="/go/src/build/${version}"

for GOOS in $os; do
  for GOARCH in $arch; do
    ext=""
    if [ "$GOOS" = "windows" ]; then
      ext=".exe"
    fi

    OSARCH_BUILDDIR="${dir}/${filename}_${GOOS}-${GOARCH}"
    FILE="${OSARCH_BUILDDIR}/${filename}${ext}"
    CGO_ENABLED=0
    export GOOS GOARCH CGO_ENABLED
    mkdir -p "$OSARCH_BUILDDIR"

    go build -ldflags "$ldflags" -o "$FILE" ./cmd/miabdns-server/main.go

  done
done
