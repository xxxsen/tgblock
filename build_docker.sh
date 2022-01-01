#!/bin/bash

VERSION=v0.0.2

docker buildx build -t xxxsen/tgblock:${VERSION} -t xxxsen/tgblock:latest --platform=linux/amd64,linux/arm64 . --push