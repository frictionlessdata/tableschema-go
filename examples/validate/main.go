package main

import (
	"log"

	"github.com/frictionlessdata/tableschema-go/csv"
	"github.com/frictionlessdata/tableschema-go/schema"
)

// Example of how to read, validate and change a schema.
func main() {
	// Reading schem.
	capitalSchema, err := schema.LoadFromFile("schema.json")
	if err != nil {
		log.Fatal(err)
	}
	// Validate schema.
	if err := capitalSchema.Validate(); err != nil {
		log.Fatal(err)
	}

	// Printing schema fields names.
	log.Println("Fields:")
	for i, f := range capitalSchema.Fields {
		log.Printf("%d - %s\n", i, f.Name)
	}

	// Working with schema fields.
	if capitalSchema.HasField("Capital") {
		log.Println("Field capital exists in schema")
	} else {
		log.Fatalf("Schema must have the field capital")
	}
	field, _ := capitalSchema.GetField("URL")
	if field.TestString("http://new.url.com") {
		value, err := field.Cast("http://new.url.com")
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
		if err := capitalSchema.CastRow(iter.Row(), &capitalRow); err != nil {
			log.Fatalf("Couldn't unmarshal row:%v err:%q", iter.Row(), err)
		}
		log.Printf("Cast Row: %+v\n", capitalRow)
	}
}
