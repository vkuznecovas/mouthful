# Mouthful import tool

This tool is used to connect to the database you're running, and import a previous dump to it. It does a stupid insert, so not a good fit for merging with existing data.

## Usage 

Simply run the main.go providing the path to mouthful config and the path to the existing dump as a command line argument like so: `go run main.go ./config.json ./mouthful.dmp`. This will then insert all the threads and comments to the database pointed at by the `config.json`.
