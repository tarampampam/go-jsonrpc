<p align="center">
  <img src="https://hsto.org/webt/lj/s8/ev/ljs8evshzjvuhkmj_325uqycvu8.png" width="128" alt="logo"/>
</p>

# `go-jsonrpc`

![Release version][badge_release_version]
![Project language][badge_language]
[![Build Status][badge_build]][link_build]
[![Coverage][badge_coverage]][link_coverage]
[![Go Report][badge_goreport]][link_goreport]
[![License][badge_license]][link_license]

This package provides JsonRPC 2.0 implementation based on interfaces with easy components extending.

<p align="center">
    <a href="https://asciinema.org/a/327835" target="_blank"><img src="https://asciinema.org/a/327835.svg" width="900"></a>
</p>

## Installation and usage

The import path for the package is `github.com/tarampampam/go-jsonrpc`.

To install it, run:

```shell script
go get github.com/tarampampam/go-jsonrpc
```

> API documentation can be [found here](https://godoc.org/github.com/tarampampam/go-jsonrpc).

### Usage example

Working examples can be found [in `examples` directory](./examples).

RPC methods definition is very simple:

```go
package main

import "github.com/tarampampam/go-jsonrpc"

type myRpcMethod struct{} // Implements `jsonrpc.Method` interface

// GetParamsType says to Router a structure (or nil) which must be used for params parsing (fields, etc.).
func (*myRpcMethod) GetParamsType() interface{} { return nil }

// GetName returns method name in string representation.
func (*myRpcMethod) GetName() string { return "my.method" }

// Handle will be called by Router when method with current name will be requested.
func (*myRpcMethod) Handle(_ interface{}) (interface{}, jsonrpc.Error) {
	return "your method response", nil
}
```

You can define method params as a structure too:

```go
package main

import "errors"

type (
    myRpcMethod struct{} // Implements `jsonrpc.Method` interface

    myRpcMethodParams struct { // Implements `jsonrpc.Validator` interface
        Number     int     `json:"number"`
        IsOptional *string `json:"is_optional"` // optional value (can be nil)
    }
)

// Implement `jsonrpc.Validator` interface for easy incoming params validation
func (p *myRpcMethodParams) Validate() error {
	if p.Number <= 0 {
		return errors.New("number must be positive")
	}

	return nil // all is ok
}

// GetParamsType NOW returns structure for method params
func (*myRpcMethod) GetParamsType() interface{} { return &myRpcMethodParams{} }

// ...
```

And use use provided kernel and router:

```go
package main

import (
	"fmt"

	"github.com/tarampampam/go-jsonrpc"
	rpcKernel "github.com/tarampampam/go-jsonrpc/kernel"
	rpcRouter "github.com/tarampampam/go-jsonrpc/router"
)

// ... methods and params definitions ...

func main () {
    // create router instance
	router := rpcRouter.New()

    // register our RPC method
    router.RegisterMethod(new(myRpcMethod))

    // create kernel using our router
    kernel := rpcKernel.New(router)

    // handle RPC request
    responseAsJSON := kernel.HandleJSONRequest([]byte(`{"jsonrpc":"2.0", "method":"ping", "id":1}`))

    fmt.Println(responseAsJSON)
}
```

### Testing

For application testing we use built-in golang testing feature and `docker-ce` + `docker-compose` as develop environment. So, just write into your terminal after repository cloning:

```shell script
$ make test
```

## Changelog

[![Release date][badge_release_date]][link_releases]
[![Commits since latest release][badge_commits_since_release]][link_commits]

Changes log can be [found here][link_changes_log].

## Support

[![Issues][badge_issues]][link_issues]
[![Issues][badge_pulls]][link_pulls]

If you will find any package errors, please, [make an issue][link_create_issue] in current repository.

## License

This is open-sourced software licensed under the [MIT License][link_license].

[badge_build]:https://img.shields.io/github/workflow/status/tarampampam/go-jsonrpc/build?maxAge=30&logo=github
[badge_coverage]:https://img.shields.io/codecov/c/github/tarampampam/go-jsonrpc/master.svg?maxAge=30
[badge_goreport]:https://goreportcard.com/badge/github.com/tarampampam/go-jsonrpc
[badge_release_version]:https://img.shields.io/github/release/tarampampam/go-jsonrpc.svg?maxAge=30
[badge_language]:https://img.shields.io/github/go-mod/go-version/tarampampam/go-jsonrpc?longCache=true
[badge_license]:https://img.shields.io/github/license/tarampampam/go-jsonrpc.svg?longCache=true
[badge_release_date]:https://img.shields.io/github/release-date/tarampampam/go-jsonrpc.svg?maxAge=180
[badge_commits_since_release]:https://img.shields.io/github/commits-since/tarampampam/go-jsonrpc/latest.svg?maxAge=45
[badge_issues]:https://img.shields.io/github/issues/tarampampam/go-jsonrpc.svg?maxAge=45
[badge_pulls]:https://img.shields.io/github/issues-pr/tarampampam/go-jsonrpc.svg?maxAge=45
[link_goreport]:https://goreportcard.com/report/github.com/tarampampam/go-jsonrpc

[link_coverage]:https://codecov.io/gh/tarampampam/go-jsonrpc
[link_build]:https://github.com/tarampampam/go-jsonrpc/actions
[link_license]:https://github.com/tarampampam/go-jsonrpc/blob/master/LICENSE
[link_releases]:https://github.com/tarampampam/go-jsonrpc/releases
[link_commits]:https://github.com/tarampampam/go-jsonrpc/commits
[link_changes_log]:https://github.com/tarampampam/go-jsonrpc/blob/master/CHANGELOG.md
[link_issues]:https://github.com/tarampampam/go-jsonrpc/issues
[link_create_issue]:https://github.com/tarampampam/go-jsonrpc/issues/new/choose
[link_pulls]:https://github.com/tarampampam/go-jsonrpc/pulls
