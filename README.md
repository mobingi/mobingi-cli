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

To build the tool, run

```
$ go build -v
```

# Usage

### Login

This is the first command you need to run to use the other commands. To login, run

```
$ mocli login --client-id=[value] --client-secret=[value]
```

This will create a file `credentials` under `[home_folder]/.mocli/` folder that will contain the API token to be used for your subsequent commands.

## Stack operations

### List stacks

To list your running stacks, run

```
$ mocli stack list
```

By default, this will list your task in a tabular form with minimal information (`--fmt=min`). If you want to display more information, run

```
$ mocli stack list --fmt=text
```

or

```
$ mocli stack list --fmt=json
```

You can also save the output to a file by adding the option `--out=[full_path_to_file]`. If your path contains whitespaces, enclose the value with double-quotes.

### Describe a stack

To describe a specific stack, run

```
$ mocli stack describe --id=[value]
```

where `[value]` is the stack id. You can get the stack id from the `stack list` command. This command supports `text`, `json`, `raw`, and `min` output formats via the `--fmt=[text|json|raw|min]` option, as well as writing to a file via the `--out=[full_path_to_file]` option.

### Delete a stack

To delete a stack, run

```
$ mocli stack delete --id=[value]
```

where `[value]` is the stack id. You can get the stack id from the `stack list` command.

## Server config operations

### Show server config

To show server config of a specific stack, run

```
$ mocli svrconf show --id=[value]
```

where `[value]` is the stack id. You can get the stack id from the `stack list` subcommand. This command supports `json`, and `raw` output formats via the `--fmt=[json|raw]` option, as well as writing to a file via the `--out=[full_path_to_file]` option.

### Update server config environment variables

To update server config environment variables, run

```
$ mocli svrconf update --id=[value] --env=KEY1:value1,KEY2:value2,KEYx:valuex...
```

or if you have whitespaces in the input, enclose it with double quotes.

```
$ mocli svrconf update --id=[value] --env="KEY1: value1, KEY2: value2, KEYx: valuex, ..."
```

To clear all environment variables, set `--env=null` option.

```
$ mocli svrconf update --id=[value] --env=null
```

The `--id=[value]` option should be the target stack id.

### Update server config filepath

To update server config filepath, run

```
$ mocli svrconf update --id=[value] --filepath=[value]
```

Note that when you provide update options simultaneously (for example, you provide `--env=value` and `--filepath=value` at the same time), the tool will send each option as a separate request.
