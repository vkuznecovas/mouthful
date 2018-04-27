# Migration tool from disqus xml to mouthful sqlite

This tool will migrate all the threads and comments from disqus xml dump to mouthful sqlite.

## Usage

Simply run the main.go, providing 1 argument: path to disqus xml dump file

`go run main.go ./disqus.xml`

This will create a mouthful.db file in the current directory as output

> special thanks to Reddit user [doenietzomoeilijk](https://www.reddit.com/user/doenietzomoeilijk) for providing the dump to make this possible
