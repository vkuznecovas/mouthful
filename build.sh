#!/bin/bash

set +x
rm -rf ./dist
# make dirs to accomodate files
mkdir -p ./dist
mkdir -p ./dist/static
mkdir -p ./dist/static/client
mkdir -p ./dist/static/admin

# bundle client
cd ./client
npm i
npm run build
mv ./build/bundle.js ../dist/static/client.js
cd ..

# bundle admin
cd ./admin
npm i
npm run build
mv ./build/* ../dist/static/admin
cd ..

# build binary
go build -o dist/mouthful main.go
chmod +x dist/mouthful

# copy over config
cp ./config.json dist/config.json
