[![Build Status](https://travis-ci.org/mobingilabs/mocli.svg?branch=master)](https://travis-ci.org/mobingilabs/mocli)
[![Build status](https://ci.appveyor.com/api/projects/status/hv1y1n3oku9frxye?svg=true)](https://ci.appveyor.com/project/flowerinthenight/mocli)

# mocli

Command line interface for Mobingi API.

### Build

This tool uses [`dep`](https://github.com/golang/dep) for dependency management. You can install `dep` via

```
$ go get -u github.com/golang/dep/cmd/dep
```

To build the tool, run

```
$ go build -v
```

### Usage

#### Login

This is the first command you need to run to use the other commands. To login, run

```
$ mocli login --client-id=value --client-secret=value
```

This will create a file `<home_folder/.mocli/credentials` that contains the API token to be used for your subsequent queries.

#### List stacks

To list your running stacks, run

```
$ mocli list
```

By default, this will list your task in a tabular form with minimal information. If you want to display more information, run

```
$ mocli list --fmt=text
```

or

```
$ mocli list --fmt=json
```

You can also save the output to a file by adding the option `--out=[full_path_to_file]`. If your path contains a whitespace, enclose the value with double-quotes.

#### Describe a stack

To describe a specific stack, run

```
$ mocli describe --id=value
```

You can get the stack id from the `list` subcommand. This command also supports `text` and `json` output formats via the `--fmt=[text|json]` option, as well as writing to a file via the `--out=[full_path_to_file]` option.
