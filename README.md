[![logo](./logos/mobinginewlogo.png)](https://mobingi.co.jp/)

[![Build Status](https://travis-ci.org/mobingilabs/mocli.svg?branch=master)](https://travis-ci.org/mobingilabs/mocli)
[![Build status](https://ci.appveyor.com/api/projects/status/hv1y1n3oku9frxye?svg=true)](https://ci.appveyor.com/project/flowerinthenight/mocli)

# mocli

```
Command line interface for Mobingi API and services.

Endpoints based on '--runenv' option:

  dev
    - https://apidev.mobingi.com
    - https://dockereg2.labs.mobingi.com

  qa
    - https://apiqa.mobingi.com
    - https://registry.mobingi.com

  prod
    - https://api.mobingi.com
    - https://registry.mobingi.com

Usage:
  mocli [command]

Available Commands:
  creds       manage your credentials
  help        help about any command
  login       login to Mobingi API
  registry    manage your docker registry
  stack       manage your stack
  svrconf     manage your server config file
  version     print the version

Flags:
      --apiver string   API version (default "v2")
      --debug           debug mode when error
  -f, --fmt string      output format (values depends on command)
  -h, --help            help for mocli
  -n, --indent int      indent padding when fmt is 'text' or 'json' (default 4)
  -o, --out string      full file path to write the output
      --runenv string   run in environment (dev, qa, prod) (default "dev")
      --rurl string     base url for Docker Registry
      --timeout int     timeout in seconds (default 120)
      --token string    access token
      --url string      base url for API
      --verbose         verbose output

Use "mocli [command] --help" for more information about a command.
```

# Getting started

### Getting mocli

The easiest way to get mocli is to use one of the [pre-built release binaries](https://github.com/mobingilabs/mocli/releases) which are available for OSX, Linux, and Windows.

If you want to try the latest version, you can build mocli from the master branch. You need to have [Go](https://golang.org/) installed (version 1.8+ is required). Note that the master branch may be in an unstable or even broken state during development.

### Building mocli

```
$ git clone https://github.com/mobingilabs/mocli
$ cd mocli
$ go build -v
$ ./mocli
```

You can also install the binary to your `$GOPATH/bin` folder (`$GOPATH/bin` should be added to your `$PATH` environment variable). 

```
$ go get -u github.com/mobingilabs/mocli
$ mocli
```

# Usage

## Login

This is the first command you need to run to use the other commands. To login, run

```
$ mocli login --client-id=foo --client-secret=bar
```

This will create a file `credentials` under `$HOME/.mocli/` folder that will contain the API token to be used for your subsequent commands.

## Stack operations

### List stacks

Examples:

```
$ mocli stack list
$ mocli stack list --fmt=text
$ mocli stack list --fmt=json
$ mocli stack list --fmt=raw --out=`echo $HOME`/out.txt
```

Enclose with double quotes if absolute file path has whitespace(s) in it.

### Describe a stack

Examples:

```
$ mocli stack describe --id=foo
$ mocli stack describe --id=foo --fmt=raw --out=/home/bar/out.txt
```

You can get the stack id from the `stack list` command.

### Delete a stack

Example:

```
$ mocli stack delete --id=foo
```

## Server config operations

### Show server config

Example:

```
$ mocli svrconf show --id=foo
```

You can get the stack id from the `stack list` command.

### Update server config

Examples:

```
$ mocli svrconf update --id=foo --env=KEY1:value1,KEY2:value2,KEYx:valuex
```

If you have whitespaces in the input, enclose it with double quotes

```
$ mocli svrconf update --id=foo --env="KEY1: value1, KEY2: value2, KEYx: valuex"
```

To clear all environment variables, set `--env=null` option

```
$ mocli svrconf update --id=foo --env=null
```

To update server config filepath, run

```
$ mocli svrconf update --id=foo --filepath=git://github.com/mobingilabs/default
```

Note that when you provide update options simultaneously (for example, you provide `--env=FOO:bar` and `--filepath=test` at the same time), the tool will send each option as a separate request.

## Vendor credentials

### View vendor credentials

Examples:

```
$ mocli creds list
$ mocli creds list --fmt=json
$ mocli creds list --fmt=raw
```

## Mobingi Docker registry

### Get token for Docker Registry API

To get token for Docker Registry API access, run

```
$ mocli registry token \
      --username=foo \
      --password=bar \
      --service="Mobingi Docker Registry" \
      --scope="repository:foo/container:*"
```

where `username` is a subuser under your Mobingi account. You can also remove `--service`, `--username` and/or `--password`.

```
$ mocli registry token --scope="repository:foo/container:*"
Username:
Password:
```

By default, it will only print the token value. To print the raw JSON output, append the `--fmt=raw` option.
 
This is useful when you want to access the registry directly using other tools. For example, you can use the token when using Docker Registry API via `curl`.

```
$ curl -H "Authorization: Bearer token" \
      -H "Accept application/vnd.docker.distribution.manifest.v2+json" \
      https://registry.mobingi.com/v2/foo/container/manifests/latest
```

### List image tags

To list image tags, run

```
$ mocli registry tags [--username=foo] [--password=bar] --image=image_name
```

### List registry catalog

To print the catalog, run

```
$ mocli registry catalog [--username=foo] [--password=bar]
```

Note that this command is inherently slow.

### Print tag manifest

To print a tag manifest, run

```
$ mocli registry manifest [--username=foo] [--password=bar] --image=hello:latest
```

You can also write the output to a file via the `--fmt=full_path_to_file` option.

### Delete a tag

To delete a tag, run

```
$ mocli registry delete [--username=foo] [--password=bar] --image=hello:latest
```

## Verbose output

You can use the global `--verbose` option if you want to see more information during the command execution.
