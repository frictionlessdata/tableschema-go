package main

import (
	"fmt"
	"log"
	"os"

	"github.com/frictionlessdata/tableschema-go/csv"
)

type user struct {
	ID   int
	Age  int
	Name string
}

func main() {
	table, err := csv.New(csv.FromFile("data_infer_utf8.csv"), csv.SetHeaders("id", "age", "name"))
	if err != nil {
		log.Fatal(err)
	}
	if err := table.Infer(); err != nil {
		log.Fatal(err)
	}
	// Writing schema to stdout.
	table.Schema.Write(os.Stdout)
	// Casting and writing data to stdout.
	var data []user
	if err := table.CastAll(&data); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n\nData:%+v\n", data)
}
