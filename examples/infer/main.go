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
	tab, err := csv.NewTable(csv.FromFile("data_infer_utf8.csv"), csv.SetHeaders("id", "age", "name"))
	if err != nil {
		panic(err)
	}
	sch, err := schema.Infer(tab)
	if err != nil {
		panic(err)
	}
	var users []user
	sch.DecodeTable(tab, &users)
}
