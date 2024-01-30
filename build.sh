#!/bin/sh

platforms="linux darwin windows"
arch="amd64 arm64"

for a in $arch; do
  for p in $platforms; do 
    printf "\033[1;32mINFO:\033[0;0m Building for $p...\n"
    env GOOS=$p GOARCH=$a go build -o ./bin/seam-$p-$a
  done
done