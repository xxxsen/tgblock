#!/bin/bash

VERSION=v0.0.1

docker buildx build -t xxxsen/tgblock:${VERSION} --platform=linux/amd64,linux/arm64 . --push