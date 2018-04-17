# Migration tool from mouthful sqlite to mouthful dynamodb

This tool will migrate all the threads and comments from mouthful sqlite to mouthful dynamodb.

## Usage

Simply run the main.go, providing 2 arguments: path to sqlite mouthful database and path to mouthful config.json containing all the required settings for dynamodb.

`go run main.go ./mouthful.db ./config.json`

For example configs, please refer to [the configuration guide](../../../examples/configs/README.md)