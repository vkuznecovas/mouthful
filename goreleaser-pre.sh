#!/bin/bash
set +x
rm -rf ./static
# make dirs to accomodate files
mkdir -p ./static
mkdir -p ./data

# bundle client
cd ./client
cp ./build/bundle.js ../static/client.js
cd ..

# bundle admin
cd ./admin
cp -r ./build/* ../static/
cd ..

# copy config
cp ./config.json ./data/config.json
