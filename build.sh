#!/bin/bash

set -euxo
rm -rf ./dist
# make dirs to accomodate files
mkdir -p ./dist
mkdir -p ./dist/static
mkdir -p ./dist/data

# bundle client
cd ./client
npm i
npm run build
cp ./build/bundle.js ../dist/static/client.js
cd ..

# bundle admin
cd ./admin
npm i
npm run build
cp -r ./build/* ../dist/static/
cd ..

# build binary
go build -tags=jsoniter -a -ldflags="-s -w" -installsuffix cgo -o dist/mouthful main.go
chmod +x dist/mouthful

# copy over config
cp ./config.json dist/data/config.json
