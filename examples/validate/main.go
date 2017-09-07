package main

import (
	"log"

	"github.com/frictionlessdata/tableschema-go/csv"
	"github.com/frictionlessdata/tableschema-go/schema"
)

// Example of how to read, validate and change a schema.
func main() {
	// Reading schem.
	capitalSchema, err := schema.ReadFromFile("schema.json")
	if err != nil {
		log.Fatal(err)
	}
	// Validate schema.
	if err := capitalSchema.Validate(); err != nil {
		log.Fatal(err)
	}

	// Printing headers.
	log.Printf("Headers: %v\n", capitalSchema.Headers())

	// Working with schema fields.
	if capitalSchema.HasField("capital") {
		log.Println("Field capital exists in schema")
	} else {
		log.Fatalf("Schema must have the field capital")
	}
	field, _ := capitalSchema.GetField("url")
	if field.TestString("http://new.url.com") {
		value, err := field.UnmarshalString("http://new.url.com")
		log.Printf("URL unmarshal to value: %v\n", value)
		if err != nil {
			log.Fatalf("Error casting value: %q", err)
		}
	} else {
		log.Fatalf("Value http://new.url.com must fit in field capital.")
	}

	// Dealing with tabular data associated with the schema.
	table, err := csv.NewTable(csv.FromFile("capital.csv"), csv.LoadHeaders())
	capitalRow := struct {
		ID      int
		Capital float64
		URL     string
	}{}

	iter, _ := table.Iter()
	for iter.Next() {
		if err := capitalSchema.Decode(iter.Row(), &capitalRow); err != nil {
			log.Fatalf("Couldn't unmarshal row:%v err:%q", iter.Row(), err)
		}
		log.Printf("Unmarshal Row: %+v\n", capitalRow)
	}
}
