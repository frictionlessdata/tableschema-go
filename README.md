[![Build Status](https://travis-ci.org/frictionlessdata/tableschema-go.svg?branch=master)](https://travis-ci.org/frictionlessdata/tableschema-go) [![Coverage Status](https://coveralls.io/repos/github/frictionlessdata/tableschema-go/badge.svg?branch=master)](https://coveralls.io/github/frictionlessdata/tableschema-go?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/frictionlessdata/tableschema-go)](https://goreportcard.com/report/github.com/frictionlessdata/tableschema-go) [![Gitter chat](https://badges.gitter.im/gitterHQ/gitter.png)](https://gitter.im/frictionlessdata/chat) [![GoDoc](https://godoc.org/github.com/frictionlessdata/tableschema-go?status.svg)](https://godoc.org/github.com/frictionlessdata/tableschema-go)

# tableschema-go
A Go library for working with Table Schema.

* [Design doc](https://goo.gl/ExQbi6)

# Getting started

## Installation

This package uses [semantic versioning 2.0.0](http://semver.org/). 

### Using dep

```sh
$ dep init
$ dep ensure -add github.com/frictionlessdata/tableschema-go@0.1
```

### Using govendor

```sh
$ govendor init
$ govendor fetch github.com/frictionlessdata/tableschema-go@0.1
```

# Example

```go
import (
    "fmt"
    "github.com/github.com/frictionlessdata/tableschema-go/table"
)

func main() {
    t, _ := table.CSVFile("data.csv", table.LoadCSVHeaders())
    fmt.Println(t.Headers)

    t.Infer()
    t.Schema.SaveToFile("schema.json")

    data := []struct {
        Name string
        Age uint16
    }{}
    t.CastAll(&data)
    fmt.Println(data)
}
```
