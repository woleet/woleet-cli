#!/usr/bin/env bash
# Based on: https://www.digitalocean.com/community/tutorials/how-to-build-go-executables-for-multiple-platforms-on-ubuntu-16-04
set -e

# Check usage
if !( [ $# == 0 ] || ( [ $# == 1 ] && [ $1 == "--compress" ] ) )
then
  echo "usage: $0 <--compress>"
  exit 1
fi

platforms=( "linux/amd64" "darwin/amd64" "windows/amd64")

rm -rf dist && mkdir -p dist

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name=woleet-cli'-'$GOOS'-'$GOARCH
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi
    env GOOS=$GOOS GOARCH=$GOARCH go build -o dist/$output_name
    if [ $# == 1 ]
    then
      upx --best --brute dist/$output_name
    fi
done

