#!/bin/sh

rm -r output

mkdir -p output

./sync_icon.sh

cp -r ./static/icons ./output/icons
cp ./static/workflow/* ./output/

GOOS=darwin GOARCH=amd64 go build -o output/devdocs_alfred_amd64
GOOS=darwin GOARCH=arm64 go build -o output/devdocs_alfred_arm64

