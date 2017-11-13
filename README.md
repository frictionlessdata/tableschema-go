[![Build Status](https://travis-ci.org/frictionlessdata/tableschema-go.svg?branch=master)](https://travis-ci.org/frictionlessdata/tableschema-go) [![Coverage Status](https://coveralls.io/repos/github/frictionlessdata/tableschema-go/badge.svg?branch=master)](https://coveralls.io/github/frictionlessdata/tableschema-go?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/frictionlessdata/tableschema-go)](https://goreportcard.com/report/github.com/frictionlessdata/tableschema-go) [![Gitter chat](https://badges.gitter.im/gitterHQ/gitter.png)](https://gitter.im/frictionlessdata/chat) [![GoDoc](https://godoc.org/github.com/frictionlessdata/tableschema-go?status.svg)](https://godoc.org/github.com/frictionlessdata/tableschema-go)

# tableschema-go

[Table schema](http://specs.frictionlessdata.io/table-schema/) tooling in Go.

# Getting started

## Installation

This package uses [semantic versioning 2.0.0](http://semver.org/). 

### Using dep

```sh
$ dep init
$ dep ensure -add github.com/frictionlessdata/tableschema-go/csv@>=0.1
```


# Main Features

## Tabular Data Load

Have tabular data stored in local files? Remote files? Packages like the [csv](https://godoc.org/github.com/frictionlessdata/tableschema-go/csv) are going to help on loading the data you need and making it ready for processing. 

```go
package main

import "github.com/frictionlessdata/tableschema-go/csv"

func main() {
   tab, err := csv.NewTable(csv.Remote("myremotetable"), csv.LoadHeaders())
   // Error handling.
}
```

Supported physical representations:

* [CSV](https://godoc.org/github.com/frictionlessdata/tableschema-go/csv)

You would like to use tableschema-go but the physical representation you use is not listed here? No problem! Please create an issue before start contributing. We will be happy to help you along the way.

## Schema Inference and Configuration

Got that new dataset and wants to start getting your hands dirty ASAP? No problems, let the [schema package](https://github.com/frictionlessdata/tableschema-go/tree/master/schema) try to infer
the data types based on the table data.

```go
package main

import (
   "github.com/frictionlessdata/tableschema-go/csv"
   "github.com/frictionlessdata/tableschema-go/schema"
)

func main() {
   tab, _ := csv.NewTable(csv.Remote("myremotetable"), csv.LoadHeaders())
   sch, _ := schema.Infer(tab)
   fmt.Printf("%+v", sch)
}
```

> Want to go faster? Please give [InferImplicitCasting](https://godoc.org/github.com/frictionlessdata/tableschema-go/schema#InferImplicitCasting) a try and let us know how it goes.

There might be cases in which the inferred schema is not correct. One of those cases is when your data use strings like "N/A" to represent missing cells. That would usually make our inferential algorithm think the field is a string.

When that happens, you can manually perform those last minutes tweaks [Schema](https://godoc.org/github.com/frictionlessdata/tableschema-go/schema#Schema).

```go
   sch.MissingValues = []string{"N/A"}
   sch.GetField("ID").Type = schema.IntegerType
```

After all that, you could persist your schema to disk:

```go
sch.SaveToFile("users_schema.json")
```

And use the local schema later:

```go
sch, _ := sch.LoadFromFile("users_schema.json")
```

Finally, if your schema is saved remotely, you can also use it:

```go
sch, _ := schema.LoadRemote("http://myfoobar/users/schema.json")
```

## Processing Tabular Data

Once you have the data, you would like to process using language data types. [schema.CastTable](https://godoc.org/github.com/frictionlessdata/tableschema-go/schema#example-Schema-CastTable) and [schema.CastRow](https://godoc.org/github.com/frictionlessdata/tableschema-go/schema#example-Schema-CastRow) are your friends on this journey.

```go
package main

import (
   "github.com/frictionlessdata/tableschema-go/csv"
   "github.com/frictionlessdata/tableschema-go/schema"
)

type user struct {
   ID   int
   Age  int
   Name string
}

func main() {
   tab, _ := csv.NewTable(csv.FromFile("users.csv"), csv.LoadHeaders())
   sch, _ := schema.Infer(tab)
   var users []user
   sch.CastTable(tab, &users)
   // Users slice contains the table contents properly raw into
   // language types. Each row will be a new user appended to the slice.
}
```

If you have a lot of data and can no load everything in memory, you can easily iterate trough it:

```go
...
   iter, _ := sch.Iter()
   for iter.Next() {
      var u user
      sch.CastRow(iter.Row(), &u)
      // Variable u is now filled with row contents properly raw
      // to language types.
   }
...
```

> Even better if you could do it regardless the physical representation! The [table](https://godoc.org/github.com/frictionlessdata/tableschema-go/table) package declares some interfaces that will help you to achieve this goal:

* [Table](https://godoc.org/github.com/frictionlessdata/tableschema-go/table#Table)
* [Iterator](https://godoc.org/github.com/frictionlessdata/tableschema-go/table#Iterator)

### Field

Class represents field in the schema.

For example, data values can be castd to native Go types. Decoding a value will check if the value is of the expected type, is in the correct format, and complies with any constraints imposed by a schema.

```javascript
{
    'name': 'birthday',
    'type': 'date',
    'format': 'default',
    'constraints': {
        'required': True,
        'minimum': '2015-05-30'
    }
}
```

The following example will raise exception the passed-in is less than allowed by `minimum` constraints of the field. `Errors` will be returned as well when the user tries to cast values which are not well formatted dates.

```go
date, err := field.Cast("2014-05-29")
// uh oh, something went wrong
```

Values that can't be castd will return an `error`.
Casting a value that doesn't meet the constraints will return an `error`.

Available types, formats and resultant value of the cast:

| Type | Formats | Casting result |
| ---- | ------- | -------------- |
| any | default | interface{} |
| object | default | interface{} |
| array | default | []interface{} |
| boolean | default | bool |
| duration | default | time.Time |
| geopoint | default, array, object | [float64, float64] |
| integer | default | int64 |
| number | default | float64 |
| string | default, uri, email, binary | string |
| date | default, any, <PATTERN> | time.Time |
| datetime | default, any, <PATTERN> | time.Time |
| time | default, any, <PATTERN> | time.Time |
| year | default | time.Time |
| yearmonth | default | time.Time |

## Saving Tabular Data

Once you're done processing the data, it is time to persist results. As an example, let us assume we have a remote table schema called `summary`, which contains two fields:

* `Date`: of type [date](https://specs.frictionlessdata.io/table-schema/#date)
* `AverageAge`: of type [number](https://specs.frictionlessdata.io/table-schema/#number) 


```go
import (
   "github.com/frictionlessdata/tableschema-go/csv"
   "github.com/frictionlessdata/tableschema-go/schema"
)


type summaryEntry struct {
    Date time.Time
    AverageAge float64
}

func WriteSummary(summary []summaryEntry, path string) {
   sch, _ := schema.LoadRemote("http://myfoobar/users/summary/schema.json")

   f, _ := os.Create(path)
   defer f.Close()

   w := csv.NewWriter(f)
   defer w.Flush()

   w.Write([]string{"Date", "AverageAge"})
   for _, summ := range summary{
       row, _ := sch.UncastRow(summ)
       w.Write(row)
   }
}
```

# API Reference and More Examples

More detailed documentation about API methods and plenty of examples is available at [https://godoc.org/github.com/frictionlessdata/tableschema-go](https://godoc.org/github.com/frictionlessdata/tableschema-go)

# Contributing

Found a problem and would like to fix it? Have that great idea and would love to see it in the repository?

> Please open an issue before start working

That could save a lot of time from everyone and we are super happy to answer questions and help you alonge the way. Furthermore, feel free to join [frictionlessdata Gitter chat room](https://gitter.im/frictionlessdata/chat) and ask questions.

This project follows the [Open Knowledge International coding standards](https://github.com/okfn/coding-standards)

* Before start coding:
     * Fork and pull the latest version of the master branch
     * Make sure you have go 1.8+ installed and you're using it
     * Make sure you [dep](https://github.com/golang/dep) installed

* Before sending the PR:

```sh
$ cd $GOPATH/src/github.com/frictionlessdata/tableschema-go
$ dep ensure
$ go test ./..
```

And make sure your all tests pass.
