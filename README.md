[![Build Status](https://travis-ci.org/mobingilabs/mocli.svg?branch=master)](https://travis-ci.org/mobingilabs/mocli)
[![Build status](https://ci.appveyor.com/api/projects/status/hv1y1n3oku9frxye?svg=true)](https://ci.appveyor.com/project/flowerinthenight/mocli)

# mocli

Command line interface for Mobingi API.

# Installation

To install `mocli`, run

```
$ go get github.com/mobingilabs/mocli
```

# Build

This tool uses [`dep`](https://github.com/golang/dep) for dependency management. You can install `dep` via

```
$ go get -u github.com/golang/dep/cmd/dep
```

If new vendor libraries are added, run

```
$ dep ensure -update -v
```

to add them to the `vendor` folder. To build the tool, run

```
$ go build -v
```

# Usage

## Login

To display usage information, run

```
$ mocli login --help
```

This is the first command you need to run to use the other commands. To login, run

```
$ mocli login --client-id=value --client-secret=value
```

This will create a file `credentials` under `{home}/.mocli/` folder that will contain the API token to be used for your subsequent commands.

## Stack operations

### List stacks

To diplay usage information, run

```
$ mocli stack list --help
```

To list your stacks, run

```
$ mocli stack list
```

By default, this will list your task in a tabular form with minimal information (`--fmt=min`). If you want to display more information, run

```
$ mocli stack list --fmt=[text|json]
```

You can also save the output to a file by adding the option `--out=[full_path_to_file]`. If your path contains whitespaces, enclose `full_path_to_file` with double-quotes.

### Describe a stack

To display usage information, run

```
$ mocli stack describe --help
```

To describe a specific stack, run

```
$ mocli stack describe --id=stack_id
```

You can get the stack id from the `stack list` command. This command supports `text`, `json`, `raw`, and `min` output formats via the `--fmt=[text|json|raw|min]` option, as well as writing to a file via the `--out=[full_path_to_file]` option.

### Delete a stack

To display usage information, run

```
$ mocli stack delete --help
```

To delete a stack, run

```
$ mocli stack delete --id=stack_id
```

You can get the stack id from the `stack list` command.

## Server config operations

### Show server config

To display usage information, run

```
$ mocli svrconf show --help
```

To show server config of a specific stack, run

```
$ mocli svrconf show --id=stack_id
```

You can get the stack id from the `stack list` command. This command supports `json`, and `raw` output formats via the `--fmt=[json|raw]` option, as well as writing to a file via the `--out=[full_path_to_file]` option.

### Update server config

To display usage information, run

```
$ mocli svrconf update --help
```

To update server config environment variables, run

```
$ mocli svrconf update --id=stack_id --env=KEY1:value1,KEY2:value2,KEYx:valuex...
```

or if you have whitespaces in the input, enclose it with double quotes.

```
$ mocli svrconf update --id=stack_id --env="KEY1: value1, KEY2: value2, KEYx: valuex, ..."
```

To clear all environment variables, set `--env=null` option.

```
$ mocli svrconf update --id=stack_id --env=null
```

To update server config filepath, run

```
$ mocli svrconf update --id=stack_id --filepath=new_file_path
```

Note that when you provide update options simultaneously (for example, you provide `--env=value` and `--filepath=value` at the same time), the tool will send each option as a separate request.

## Mobingi Docker registry

### Get token for Docker Registry API

To display usage information, run

```
$ mocli registry token --help
```

To get token for Docker Registry API access, run

```
$ mocli registry token \
      --username=foo \
      --password=bar \
      --account=foo \
      --service="Mobingi Docker Registry" \
      --scope="repository:foo/container:*"
```

where `username` is a subuser under your Mobingi account. You can also remove `--username` and/or `--password`.

```
$ mocli registry token --scope="repository:foo/container:*"
Username:
Password:
```

You can then use the token when using Docker Registry API. For example,

```
$ curl -H "Authorization: Bearer token" \
      -H "Accept application/vnd.docker.distribution.manifest.v2+json" \
      https://registry.mobingi.com/v2/foo/container/manifests/latest
```
