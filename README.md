# Storm v0.2.0

[![CircleCI](https://circleci.com/gh/medtune/storm.svg?style=svg)](https://circleci.com/gh/medtune/storm)


## Description

Storm is a dataset generator tool based on google search engine. 

Supported search data types: 
 - Images (PNG, JPEG)

Supported output data types:
 - Tensorflow records (Protocol buffer)
 - Simple rendering (JPEG / PNG under floders)

## Install storm


Simple way: using `go get` command

`go get -u github.com/medtune/storm/...`

From binaries

```shell
#TODO
```

## Usage

*Requirements :*

- Google search engine ID
- Google Cloud Platform credentials

make sure to export GCP Credentials file by 

```shell
# GCP Creds
export GOOGLE_APPLICATION_CREDENTIALS=path/to/my/crendentials/file.json
# Search engine ID
export GOOGLE_CSE_ID=MY_SEARCH_ENGINE_ID
```

Basic example:

```shell
storm -q "cute cats" \     #search query
      -e $ENGINE_ID \      #engine ID 
      -n 30 \              #number of queries to make
      -r linear:300x300 \  #image resize before save
      -o cats              #output file name
```
## Command specs

`#TODO`
