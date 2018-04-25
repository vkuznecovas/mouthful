[![Build Status](https://travis-ci.org/vkuznecovas/mouthful.svg?branch=master)](https://travis-ci.org/vkuznecovas/mouthful)
[![Go Report Card](https://goreportcard.com/badge/github.com/vkuznecovas/mouthful)](https://goreportcard.com/report/github.com/vkuznecovas/mouthful)
[![codecov](https://codecov.io/gh/vkuznecovas/mouthful/branch/master/graph/badge.svg)](https://codecov.io/gh/vkuznecovas/mouthful)
[![Documentation](https://godoc.org/github.com/vkuznecovas/mouthful?status.svg)](https://godoc.org/github.com/vkuznecovas/mouthful) 

# Mouthful is a self-hosted alternative to Disqus.

Mouthful is a lightweight commenting server written in GO and Preact. It's a self hosted alternative to disqus that's ad free.

There's a demo hosted at [mouthful.dizzy.zone](https://mouthful.dizzy.zone). Check it out!

* [Features](#features)
* [Get mouthful](#installation)
    * [The easy way](#the-easy-way)
    * [Build from source](#building-mouthful-yourself)
    * [Run it on Docker](#mouthful-on-docker)
* [Configure mouthful](#configuring-mouthful)
    * [Moderation](#moderation)
    * [Caching](#caching)
    * [Rate limiting](#rate-limiting)
    * [Styling](#styling)
    * [Cross-origin resource sharing](#cross-origin-resource-sharing)
    * [Data sources](#data-source)
    * [Config file from Docker image](#config-file-from-docker)
* [Contributing](#contributing)
* [Wish list](#wish-list)
* [Get in touch](#get-in-touch)
* [Who uses mouthful](#who-uses-mouthful?)

# Features

* Multiple database support(sqlite, mysql, postgres, dynamodb)
* Moderation with an administration panel
* Server side caching to prevent excessive database calls
* Rate limiting
* Honeypot feature, to prevent bots from posting comments
* Migrations from existing commenting engines(isso)
* Configuration - most of the features can be turned on or off, as well as customized to your preferences.

# Installation

## The easy way

### Backend

Head over to [release](https://github.com/vkuznecovas/mouthful/releases) page and download an archive for your OS. Extract, change the config.json you find in the archive according to your preferences. For more info on configuration, head to the [configuration section](#configuring-mouthful).

Run the binary and that's it! You now have the backend running. 

### Client 

Now, all that's left to do is include the following html in your blog/website on the pages you want mouthful to work on:

```
<div id="mouthful-comments" data-url="http://localhost:8080"></div>
<script src="http://localhost:8080/client.js"></script>
```

Once that is set up, you should be able to start using mouthful.

## Building mouthful yourself

To start using mouthful, you'll need:

* A working GO environment
* Dep
* Node with npm
* A server to put mouthful on

> If you do not have these tools set up, please refer to the [installing dependencies section](#installing-dependencies).

If you have all the dependencies, clone the mouthful repository. In the root of this repository run `build.sh`. Give it some time, this will install all the dependencies for both go and node and create a directory inside the root of this repository called `/dist`. Inside, you'll find all you need to run mouthful. That is:

* A config file
* A binary to start the mouthful backend
* A static directory, containing all the javascript and html needed to serve both the client and the admin panel(if enabled)

To configure your mouthful instance to your hearts content, please refer to the [configuration section](#configuring-mouthful).

Once you've done with the configuration, just copy over the `/dist` contents to your server and run the `/dist/mouthful` binary. Take note that the mouthful binary will look for a config.json file its directory.

### Installing dependencies

1) To install Go, please refer to the GO documentation found [here](https://golang.org/doc/install#install)
1) To install node and npm, please refer to the Node documentation found [here](https://nodejs.org/en/download/package-manager/)
1) To install Dep, please refer to Dep documentation found [here](https://github.com/golang/dep#installation)
1) Once you have all the tools installed, follow the [Installation guide](#installation)

## Mouthful on Docker

### Build the image

1. Clone the project

```sh
git clone https://github.com/vkuznecovas/mouthful.git
```

2. Get in the project folder then build the image
    
```sh
docker build -t mouthful .
```

The [`Dockerfile`](Dockerfile) is going to build on the master branch by default, you can specify a version

```sh
docker build --build-args "MOUTHFUL_VER=1.0.3" -t mouthful .
```

### Run the image

Once image is built, simply run

```sh
docker run -d \
    --name mouthful \
    -v $(pwd)/data:/app/data
    -p 8080:8080
    mouthful
```

Alternatively you can use the official image `vkuznecovas/mouthful`

```sh
docker run -d \
    --name mouthful \
    -v $(pwd)/data:/app/data
    -p 8080:8080
    vkuznecovas/mouthful
```

**Note:** `/app/data` needs to contain a valid `config.json` file, read the note in [moderation](#moderation). You can extract the config file from the docker image, see [getting config file from docker](#config-file-from-docker).

# Configuring mouthful

Nearly all the features of mouthful can be customized and turned on or off. All within the config.json file.

Here's a short overview:

## Moderation

Mouthful comes with moderation support out of the box. If moderation is enabled, it does not show the comments users post instantly, those will have to be approved first through the mouthful admin panel. This also allows for comment modification or deletion.

**Note:** You need to change the default password in [config.json](config.json#L5), else `mouthful` will fail to start.

## Caching

Mouthful can cache end results(full sets of comments for threads) for a given period of time. This allows for quicker responses, lower number of database queries at the cost of extra memory for the running mouthful binary.

## Rate limiting

Mouthful can limit the amount of posts a person can post within the same hour.

## Styling

Mouthful comes with a default style out of the box, but you can override it in a couple of ways:

1) Disable the default styling in config and add the required css to your webpage.
2) Fork the repo and change the style in `client/src/components/client/style.scss`.

## Paging of comments

Mouthful can either display all the comments on page load, or page them. The page size can be specified in config.

## Cross-Origin Resource Sharing

Mouthful can either allow all origins to access its backend from browser or limit that to a given list of domains.

## Data sources

Mouthful supports different data stores for different needs. Currently supported data store list is as follows:

* sqlite
* postgres
* mysql
* aws dynamodb

For a list of configuration options and config file examples, head over to [configuration documentation and examples](./examples/configs/README.md)

## Config file from Docker

You can get the default `config.json` by running

```sh
docker run --rm vkuznecovas/mouthful cat /app/data/config.json > config.json
```

This will create a file named `config.json` in your host machine, you can edit it as you please. Make sure it is present in the `data` folder before runnig the docker image, read the note in [run the image](#run-the-image).

# Contributing

Contributions are more than welcome. If you've found a bug, raise an issue. If you've got a feature request, open up an issue as well. I'll try and keep the api stable, as well as tag each release with a semantic version.

# Wish list

I'm a keen backender and not too sharp on the frontend part. If you're willing to contribute, front end(both client and admin) are not in the best of shapes, especially the admin panel. Frontend might require a refactor. Any addition of tests would be great as well. Migrations from other commenting engines would be encouraged as well. If someone could send me a disqus dump, I'd make a migration for that.

# Get in touch

If you'd like to get in touch with me, you can drop me a message on [my twitter](https://twitter.com/DizzyZoneBlog).
# Who uses mouthful?

* [Mouthful authors' blog, dizzy.zone](https://dizzy.zone)

> Feel free to do a PR and include yourself if you end up using mouthful.
