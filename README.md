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
    t, _ := csv.New(csv.FileSource("data.csv"), csv.LoadHeaders())  // load table
    t.Infer()  // infer the table schema
    t.Schema.SaveToFile("schema.json")  // save inferred schema to file
    data := []person
    t.CastAll(&data)  // casts the table data into the data slice.
}
```
# Documentation

## Table

A table is a core concept in a tabular data world. It represents a data with a metadata (Table Schema). Let's see how we could use it in practice.

Consider we have some local CSV file, `data.csv`:

```csv
city,location
london,"51.50,-0.11"
paris,"48.85,2.30"
rome,N/A
```

To read its contents we use [csv.New](https://godoc.org/github.com/frictionlessdata/tableschema-go/csv#New) to create a table and use the [File](https://godoc.org/github.com/frictionlessdata/tableschema-go/csv#FromFile) [Source](https://godoc.org/github.com/frictionlessdata/tableschema-go/csv#Source).

```go
    table, _ := csv.New(csv.FromFile("data.csv"), csv.LoadHeaders())
    table.Headers // ["city", "location"]
    table.All() // [[london 51.50,-0.11], [paris 48.85,2.30], [rome N/A]]
```

As we could see our locations are just a strings. But it should be geopoints. Also Rome's location is not available but it's also just a N/A string instead of go's zero value. First we have to infer Table Schema:

```go
    table.Infer()
    fmt.Println(table.Schema)
	// "fields": [
	//     {"name": "city", "type": "string", "format": "default"},
	//     {"name": "location", "type": "geopoint", "format": "default"},
	// ],
	// "missingValues": []
    // ...
```

Then we could create a struct and automatically cast the table data to schema types. It is like [json.Unmarshal](https://golang.org/pkg/encoding/json/#Unmarshal), but for table rows. First thing we need is to create the struct which will represent each row.

```go
type Location struct {
    City string
    Location schema.GeoPoint
}
```

Then we are ready to cast the table.

```go
var locations []Location
table.CastAll(&locations)
// Fails with cast error: "Invalid geopoint:\"N/A\""
```

The problem is that the library does not know that N/A is not an empty value. For those cases, there is a `missingValues` property in Table Schema specification. As a first try we set `missingValues` to N/A in table.Schema.

```go
table.Schema.MissingValues = []string{"N/A"}
var locations []Location
table.CastAll(&locations)
fmt.Println(rows)
// [{london {51.5 -0.11}} {paris {48.85 2.3}} {rome {0 0}}]
```

And because there are no errors on data reading we could be sure that our data is valid againt our schema. Let's save it:

```go
table.Schema.SaveToFile("schema.json")
```