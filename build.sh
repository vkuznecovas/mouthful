#!/bin/bash

set +x
rm -rf ./dist
# make dirs to accomodate files
mkdir -p ./dist
mkdir -p ./dist/static

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

# install deps
dep ensure

# build binary
go build -o dist/mouthful main.go
chmod +x dist/mouthful

# copy over config
cp ./config.json dist/config.json
