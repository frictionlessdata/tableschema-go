package examples

import (
	"fmt"

	"github.com/frictionlessdata/tableschema-go/csv"
	"github.com/frictionlessdata/tableschema-go/schema"
)

type csvRow struct {
	Name string
}

func ExampleTable_Iter() {
	tab, _ := csv.New(csv.StringSource("\"name\"\nfoo\nbar"), csv.LoadHeaders())
	tab.Schema = &schema.Schema{Fields: []schema.Field{{Name: "name", Type: schema.StringType}}}
	iter, _ := tab.Iter()
	for iter.Next() {
		var data csvRow
		iter.CastRow(&data)
		fmt.Println(data.Name)
	}
	// Output:foo
	// bar
}

func ExampleTable_Infer() {
	tab, _ := csv.New(csv.StringSource("\"name\"\nfoo\nbar"), csv.LoadHeaders())
	if err := tab.Infer(); err != nil {
		fmt.Println(err)
	}
	iter, _ := tab.Iter()
	for iter.Next() {
		var data csvRow
		iter.CastRow(&data)
		fmt.Println(data.Name)
	}
	// Output:foo
	// bar
}
