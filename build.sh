set +x
rm -rf ./dist
# make dirs to accomodate files
mkdir -p ./dist
mkdir -p ./dist/static

# create client config
go run scripts/transformConfig.go ./config.json
mv config.front.json client/src/components/client/config.json

# bundle client
cd ./client
npm i
npm run build
mv ./build/bundle.js ../dist/static/config.js
cd ..

# bundle admin
cd ./admin
npm i
npm run build
mv ./build/* ../dist/static
cd ..

# build binary
go build -o dist/mouthful main.go

# copy over config
mv ./config.json dist/config.json