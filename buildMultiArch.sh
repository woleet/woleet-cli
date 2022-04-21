#!/usr/bin/env bash
# Based on: https://www.digitalocean.com/community/tutorials/how-to-build-go-executables-for-multiple-platforms-on-ubuntu-16-04
set -e

# Check usage
if ! { [ "$#" == "0" ] || { [ "$#" == "1" ] && { [ "$1" == "--linux-static*" ] || [ "$1" == "--linux-static-alpine" ] ;} ;} ;}
then
  echo "usage: $0 <--linux-static || --linux-static-alpine>"
  exit 1
fi

RELEASE_NUMBER=$(cat cmd/root.go | grep 'Version:' | grep -oE '[0-9]+\.[0-9]+\.[0-9]+' ) 

platforms=( "linux/amd64" "darwin/amd64" "windows/amd64" )

rm -rf dist && mkdir -p dist

for platform in "${platforms[@]}"
do
  (
    platform_split=( ${platform//\// } )
    GOOS="${platform_split[0]}"
    GOARCH="${platform_split[1]}"
    if [ "$#" == "1" ] && [ "$GOOS" == "linux" ]
    then
      if [ "$1" == "--linux-static" ]
      then
        CC='/usr/local/musl/bin/musl-gcc'
      elif [ "$1" == "--linux-static-alpine" ]
      then
        CC='/usr/bin/x86_64-alpine-linux-musl-gcc'
      fi
      GOOS="$GOOS" GOARCH="$GOARCH" CC="$CC" go build -buildvcs=false --ldflags '-linkmode external -extldflags "-static"' -o dist/woleet-cli
    else
      GOOS="$GOOS" GOARCH="$GOARCH" go build -buildvcs=false -o dist/woleet-cli
    fi
    cd dist || exit
    if [ "$GOOS" != "windows" ]
    then
      tar -czf "woleet-cli_${RELEASE_NUMBER}_${GOOS}_x86_64.tar.gz" woleet-cli
    else
      mv woleet-cli woleet-cli.exe
      zip -q "woleet-cli_${RELEASE_NUMBER}_${GOOS}_x86_64.zip" woleet-cli.exe
    fi
    rm -f woleet-cli woleet-cli.exe
  )
done
