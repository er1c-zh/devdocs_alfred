#!/bin/sh

rm -r output

mkdir -p output

./sync_icon.sh

cp -r ./static/icons ./output/icons
cp ./static/workflow/* ./output/

go build -o output/devdocs_alfred
