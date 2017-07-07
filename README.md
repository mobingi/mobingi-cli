[![Build Status](https://travis-ci.org/mobingilabs/mocli.svg?branch=master)](https://travis-ci.org/mobingilabs/mocli)
[![Build status](https://ci.appveyor.com/api/projects/status/hv1y1n3oku9frxye?svg=true)](https://ci.appveyor.com/project/flowerinthenight/mocli)

# mocli

Command line interface for Mobingi API.

### Build

This tool uses [`dep`](https://github.com/golang/dep) for dependency management. Install `dep` via

```
$ go get -u github.com/golang/dep/cmd/dep
```

To build the tool, run

```
$ dep ensure -v
$ go build -v
```
