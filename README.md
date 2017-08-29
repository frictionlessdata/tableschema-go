[![Build Status](https://travis-ci.org/frictionlessdata/tableschema-go.svg?branch=master)](https://travis-ci.org/frictionlessdata/tableschema-go) [![Coverage Status](https://coveralls.io/repos/github/frictionlessdata/tableschema-go/badge.svg?branch=master)](https://coveralls.io/github/frictionlessdata/tableschema-go?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/frictionlessdata/tableschema-go)](https://goreportcard.com/report/github.com/frictionlessdata/tableschema-go) [![Gitter chat](https://badges.gitter.im/gitterHQ/gitter.png)](https://gitter.im/frictionlessdata/chat) [![GoDoc](https://godoc.org/github.com/frictionlessdata/tableschema-go?status.svg)](https://godoc.org/github.com/frictionlessdata/tableschema-go)

# tableschema-go
A Go library for working with [Table Schema](http://specs.frictionlessdata.io/table-schema/).


# Main Features

* [table](https://godoc.org/github.com/frictionlessdata/tableschema-go/table) package defines [Table](https://godoc.org/github.com/frictionlessdata/tableschema-go/csv#Table) and `Iterator` interfaces, which are used to manipulate and/or explore tabular data;
* [csv](https://godoc.org/github.com/frictionlessdata/tableschema-go/csv) package contains implementation of Table and Iterator interfaces to deal with CSV format;
* [schema](https://github.com/frictionlessdata/tableschema-go/tree/master/schema) package contains classes and funcions for working with table schemas:
     * [Schema](https://godoc.org/github.com/frictionlessdata/tableschema-go/schema#Schema): main entry point, used to validate and deal with schemas
     * [Field](https://godoc.org/github.com/frictionlessdata/tableschema-go/schema#Field): for working with schema fields
     * [Infer](https://godoc.org/github.com/frictionlessdata/tableschema-go/schema#Schema) and [InferImplicitCasting](https://godoc.org/github.com/frictionlessdata/tableschema-go/schema#InferImplicitCasting): for inferring a schema of tabular data
     

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

# Examples

Code examples in this readme requires Go 1.8+. You can find more examples in the [examples](https://github.com/frictionlessdata/tableschema-go/tree/master/examples) directory.

```go
import (
    "fmt"
    "github.com/github.com/frictionlessdata/tableschema-go/csv"
)

// struct representing each row of the table.
type person struct {
    Name string
    Age uint16
}

func main() {
    t, _ := csv.New(
        csv.FileSource("data.csv"), csv.LoadHeaders())  // load table
    t.Infer()  // infer the table schema
    t.Schema.SaveToFile("schema.json")  // save inferred schema to file
    data := []person
    t.CastAll(&data)  // casts the table data into the data slice.
}
```
