#!/usr/bin/env bash
set -e 

DOCKERFILE=$(printf 'FROM golang:1.13-alpine\nRUN apk add --no-cache build-base bash zip\nENTRYPOINT ["bash"]')

echo "$DOCKERFILE" | docker build -t woleet-cli-builder -

docker run -it --rm -v "$PWD:/woleet-cli" -w "/woleet-cli" woleet-cli-builder -c "./buildMultiArch.sh --linux-static-alpine && chown -R $(id -u "$(whoami)"):$(id -g "$(whoami)") dist"
