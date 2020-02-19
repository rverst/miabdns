#!/usr/bin/env bash
# This script builds release versions of the miabdns-server using the official golang docker image.

go=1.13
os="linux darwin windows"
arch="amd64 386"
builddate=$(date -u -Iseconds)
version=""
comhash="unknown"
slim=0

cd "$(dirname "${BASH_SOURCE[0]}")" || exit 1
cd .. || exit 1

while getopts ":o:a:v:x" opt; do
  case $opt in
    o) os="$OPTARG"
    ;;
    a) arch="$OPTARG"
    ;;
    v) version="$OPTARG"
    ;;
    x) slim=1
    ;;
    \?) echo "Invalid option -$OPTARG" >&2
    ;;
  esac
done

if [[ "$version" ]]; then
  version=$(./tools/check_semver.sh -v "$version")
  exitcode=$?
  [ $exitcode != 0 ] && { printf '%s\n' "check_semver exited with non-zero ($exitcode)"; exit 2; }
fi

if [[ $(command -v git) ]]; then
  comhash=$(git rev-parse HEAD)
  if [[ ! "$version" ]]; then
    tags=$(git tag -l 'v*' --points-at "$comhash")
    version=$(./tools/check_semver.sh -v "$tags")
    exitcode=$?
    if [[ ! $version ]]; then
        printf "Can't find a valid version tag on HEAD, please provide the version as parameter.\n"
        printf  "Either there is no valid tag (e.g. v1.0.0) or there is more than one.\n"
        if [[ $tags ]]; then
          printf  "Tags found:\n"
          printf  "%s\n\n" "$tags"
        fi
        exit 2
      fi
  fi
fi

ldflags="-w -s -X 'main.Version=$version' -X 'main.CommitHash=$comhash' -X 'main.BuildDate=$builddate'"
filename="miabdns"
builddir="$(pwd)/build/$version"
license="$(pwd)/LICENSE.txt"
readme="$(pwd)/README.md"
config="$(pwd)/example_config.json"


if [[ -d "$builddir" ]]; then
  rm -r "$builddir"
fi
mkdir -p "$builddir"

docker run --rm -it -v "$(pwd)":/go/src golang:"$go" /go/src/tools/build.sh "$version" "$os" "$arch" "$ldflags" "$filename"
#docker run --rm -it -v "$(pwd)":/go/src golang:"$go" /bin/bash; exit

if [[ $slim -eq 1 ]]; then
  echo "$version"
  exit 0
fi

chksum="$builddir/miabdns_${version}_checksums.txt"

touch "$chksum"
echo "sha256 checksums" > "$chksum"

for d in "$builddir"/*; do

  if [[ ! -d $d ]]; then
    continue
  fi

  cp "$license" "$d"
  cp "$readme" "$d"
  cp "$config" "$d"

  cd "$d" || continue
  if [[ -f "$filename.exe" ]]; then
    archive="$d.zip"
    zip -r "$archive" "./"
  else
    archive="$d.tar.gz"
    tar -czvf "$archive" "./"
  fi

  archive="$(basename -- "$archive")"
  cd "$builddir" || continue
  sha256sum "$archive" >>"$chksum"

done
exit 0