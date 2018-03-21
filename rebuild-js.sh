#!/bin/bash
set +x
rm -rf ./static
mkdir ./static
# create client config
go run cmd/util/transformConfig.go ./config.json
cp config.front.json client/src/components/client/config.json
mv config.front.json admin/src/routes/panel/config.json

# bundle client
cd ./client
npm run build
mv ./build/bundle.js ../static/client.js
cd ..

# bundle admin
cd ./admin
npm run build
mv ./build/* ../static
cd ..
