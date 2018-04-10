#!/bin/bash
set +x
rm -rf ./static
mkdir ./static
mkdir ./static/admin
mkdir ./static/client

# bundle client
cd ./client
npm run build
cp ./build/bundle.js ../static/client.js
cd ..

# bundle admin
cd ./admin
npm run build
cp -R ./build/* ../static/
cd ..
