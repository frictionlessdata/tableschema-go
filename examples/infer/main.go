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
	reader, err := csv.NewReader(csv.FromFile("data_infer_utf8.csv"), csv.SetHeaders("id", "age", "name"), csv.InferSchema())
	if err != nil {
		log.Fatal(err)
	}
	// Writing schema to stdout.
	reader.Schema.Write(os.Stdout)
	// Casting and writing data to stdout.
	var data []user
	if err := reader.UnmarshalAll(&data); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n\nData:%+v\n", data)
}
