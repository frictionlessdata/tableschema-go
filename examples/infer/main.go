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
	tab, err := csv.NewTable(csv.FromFile("data_infer_utf8.csv"), csv.SetHeaders("ID", "Age", "Name"))
	if err != nil {
		panic(err)
	}
	fmt.Println("## Raw Table ##")
	fmt.Println(tab)
	sch, err := schema.Infer(tab)
	if err != nil {
		panic(err)
	}

	fmt.Println("## Schema ##")
	fmt.Println(sch)
	var users []user
	sch.CastTable(tab, &users)

	fmt.Printf("\n## Cast Table ##\n%+v\n", users)
}
