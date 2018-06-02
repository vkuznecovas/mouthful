# Mouthful export tool

This tool is used to connect to the database you're running, and dump all the comments and threads to a dump file. 

These can then be archive and reimported.

## Usage 

Simply run the main.go providing the path to mouthful config as a command line argument like so: `go run main.go ./config.json`. This will create a `mouthful.dmp` file in the current directory with a dump of all the threads and comments.
