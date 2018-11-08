#!/usr/bin/env bash
# Based on: https://www.digitalocean.com/community/tutorials/how-to-build-go-executables-for-multiple-platforms-on-ubuntu-16-04
set -e

# Check usage
if ! { [[ $1 =~ [0-9]\.[0-9]\.[0-9] ]] && { [ "$#" == "1" ] || { [ "$#" == "2" ] && [ "$2" == "--linux-static" ] ;} ;} ;}
then
  echo "usage: $0 (release number) <--linux-static>"
  echo "release number format: [0-9].[0-9].[0-9]"
  exit 1
fi

RELEASE_NUMBER=$1
if [ "$#" == "2" ] && [ "$2" == "--linux-static" ]
then
  STATIC=true
else
  STATIC=false
fi

platforms=( "linux/amd64" "darwin/amd64" "windows/amd64" )

rm -rf dist && mkdir -p dist

for platform in "${platforms[@]}"
do
  (
    platform_split=( ${platform//\// } )
    GOOS="${platform_split[0]}"
    GOARCH="${platform_split[1]}"
    if [ "$STATIC" == "true" ] && [ "$GOOS" == "linux" ]
    then
      GOOS="$GOOS" GOARCH="$GOARCH" CC=/usr/local/musl/bin/musl-gcc go build --ldflags '-linkmode external -extldflags "-static"' -o dist/woleet-cli
    else
      GOOS="$GOOS" GOARCH="$GOARCH" go build -o dist/woleet-cli
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
