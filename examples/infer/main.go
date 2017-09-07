package main

import (
	"fmt"

	"github.com/frictionlessdata/tableschema-go/csv"
	"github.com/frictionlessdata/tableschema-go/schema"
)

type user struct {
	ID   int
	Age  int
	Name string
}

func main() {
	table, err := csv.NewTable(csv.FromFile("data_infer_utf8.csv"), csv.SetHeaders("id", "age", "name"))
	if err != nil {
		panic(err)
	}

	userSchema, err := schema.Infer(table)
	if err != nil {
		panic(err)
	}
	fmt.Println(userSchema)

	iter, _ := table.Iter()
	defer iter.Close()
	for iter.Next() {
		var data user
		if err := userSchema.Decode(iter.Row(), &data); err != nil {
			panic(err)
		}
		fmt.Println(data)
	}

	// decoder, err := schema.NewDecoder(table, userSchema)
	// for decoder.More() {
	// 	var data user
	// 	if err := decoder.Decode(&data); err != nil {
	// 		panic(err)
	// 	}
	// 	fmt.Println(data)
	// }
}
