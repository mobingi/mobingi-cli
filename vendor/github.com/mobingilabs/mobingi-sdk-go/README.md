[![logo](./logos/mobingi-205x119.png)](https://mobingi.co.jp/)

[![Build Status](https://travis-ci.org/mobingilabs/mobingi-sdk-go.svg?branch=master)](https://travis-ci.org/mobingilabs/mobingi-sdk-go)
[![Build status](https://ci.appveyor.com/api/projects/status/7085b5hnw6ehbdh9/branch/master?svg=true)](https://ci.appveyor.com/project/flowerinthenight/mobingi-sdk-go/branch/master)

This sdk uses [dep](https://github.com/golang/dep) as its vendor manager. To install dependencies and run tests:

```bash
# install dep
$ go get -u -v github.com/golang/dep/...

# install (update) dependencies
$ dep ensure -v && dep ensure -update -v

# run tests
$ go test ./... -cover
```
