[![logo](./logos/mobinginewlogo.png)](https://mobingi.co.jp/)

[![Build Status](https://travis-ci.org/mobingi/mobingi-cli.svg?branch=master)](https://travis-ci.org/mobingi/mobingi-cli)
[![Build status](https://ci.appveyor.com/api/projects/status/k7tmyr3l4dws4usd/branch/master?svg=true)](https://ci.appveyor.com/project/flowerinthenight/mobingi-cli/branch/master)

# mobingi-cli

mobingi-cli is the official command line interface for Mobingi API and services. 

See the documentation on https://learn.mobingi.com/enterprise/cli.

Documentation is written in markdown and can be located [here](https://github.com/mobingi/mobingi/blob/docs/docs/markdown/enterprise/doc-cli.md).

# Getting started

### Getting mobingi-cli

The easiest way to get mobingi-cli is to use one of the [pre-built release binaries](https://github.com/mobingi/mobingi-cli/releases) which are available for OSX, Linux, and Windows.

If you want to try the latest version, you can build mobingi-cli from the master branch. You need to have [Go](https://golang.org/) installed (version 1.8+ is required). Note that the master branch may be in an unstable or even broken state during development.

### Building mobingi-cli

```
$ git clone https://github.com/mobingi/mobingi-cli
$ cd mobingi-cli
$ go build -v
$ ./mobingi-cli
```

You can also install the binary to your `$GOPATH/bin` folder (`$GOPATH/bin` should be added to your `$PATH` environment variable). 

```
$ go get -u github.com/mobingi/mobingi-cli
$ mobingi-cli
```
