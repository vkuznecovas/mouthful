#!/bin/bash
set +x
rm -rf ./static
# make dirs to accomodate files
mkdir -p ./static


# bundle client
cd ./client
npm run build
cp ./build/bundle.js ../static/client.js
cd ..

# bundle admin
cd ./admin
npm run build
cp -r ./build/* ../static/
cd ..
