package schema

import (
	"fmt"
	"sort"
	"testing"

	"github.com/matryer/is"

	"github.com/frictionlessdata/tableschema-go/table"
)

// Demonstrations.

/* Demonstrations of incorrect inference.
	NOTE 	These are here only for review experimentation -
			they fail and won't pass CI tests.

// ExampleInferBoolWrongly should infer Value to be number, but infers bool.
/* TODO: uncomment once discussion of InferWithPrecedence is complete.
		 Commented out for now to avoid failing tests, as it's only a demo.
func ExampleInferBoolWrongly() {
	tab := table.FromSlices(
		"ExampleTable",
		[]string{"Item", "Value"},
		[][]string{
			[]string{"A", "1.2"},
			[]string{"B", "0"},
			[]string{"C", "0"},
			[]string{"D", "0"},
			[]string{"E", "5.2"},
		})
	s, _ := Infer(tab)
	fmt.Println("Fields:")
	for _, f := range s.Fields {
		fmt.Printf("{Name:%s Type:%s Format:%s}\n", f.Name, f.Type, f.Format)
	}
	// Output: Fields:
	// {Name:Item Type:string Format:default}
	// {Name:Value Type:number Format:default}
}

// ExampleInferStringWrongly should infer Value to be number, but infers string.
func ExampleInferStringWrongly() {
	tab := table.FromSlices(
		"ExampleTable",
		[]string{"Item", "Value"},
		[][]string{
			[]string{"A", "1.2"},
			[]string{"B", "0"},
			[]string{"C", ""},
			[]string{"D", ""},
			[]string{"E", ""},
			[]string{"F", ""},
			[]string{"G", ""},
			[]string{"H", "0"},
			[]string{"I", "5.2"},
		})
	s, _ := Infer(tab)
	fmt.Println("Fields:")
	for _, f := range s.Fields {
		fmt.Printf("{Name:%s Type:%s Format:%s}\n", f.Name, f.Type, f.Format)
	}
	// Output: Fields:
	// {Name:Item Type:string Format:default}
	// {Name:Value Type:number Format:default}
}
*/

// Demostration of precedence inference.

// ExampleInferPrecedenceNumberNotBool correctly infers as number.
func ExampleInferPrecedenceNumberNotBool() {
	tab := table.FromSlices(
		"ExampleTable",
		[]string{"Item", "Value"},
		[][]string{
			[]string{"A", "1.2"},
			[]string{"B", "0"},
			[]string{"C", "0"},
			[]string{"D", "0"},
			[]string{"E", "5.2"},
		})
	s, _ := Infer(tab, InferWithPrecedence(true))
	fmt.Println("Fields:")
	for _, f := range s.Fields {
		fmt.Printf("{Name:%s Type:%s Format:%s}\n", f.Name, f.Type, f.Format)
	}
	// Output: Fields:
	// {Name:Item Type:string Format:default}
	// {Name:Value Type:number Format:default}
}

// ExampleInferPrecedenceNumberNotString correctly infers as number.
func ExampleInferPrecedenceNumberNotString() {
	tab := table.FromSlices(
		"ExampleTable",
		[]string{"Item", "Value"},
		[][]string{
			[]string{"A", "1.2"},
			[]string{"B", "0"},
			[]string{"C", ""},
			[]string{"D", ""},
			[]string{"E", ""},
			[]string{"F", ""},
			[]string{"G", ""},
			[]string{"H", "0"},
			[]string{"I", "5.2"},
		})
	s, _ := Infer(tab, InferWithPrecedence(true))
	fmt.Println("Fields:")
	for _, f := range s.Fields {
		fmt.Printf("{Name:%s Type:%s Format:%s}\n", f.Name, f.Type, f.Format)
	}
	// Output: Fields:
	// {Name:Item Type:string Format:default}
	// {Name:Value Type:number Format:default}
}

// End of Demonstrations

func Exampleinfer() {
	tab := table.FromSlices(
		"ExampleTable",
		[]string{"Person", "Height"},
		[][]string{
			[]string{"Foo", "5"},
			[]string{"Bar", "4"},
			[]string{"Bez", "5.5"},
		})
	s, _ := Infer(tab)
	fmt.Println("Fields:")
	for _, f := range s.Fields {
		fmt.Printf("{Name:%s Type:%s Format:%s}\n", f.Name, f.Type, f.Format)
	}
	// Output: Fields:
	// {Name:Person Type:string Format:default}
	// {Name:Height Type:integer Format:default}
}

func ExampleInfer_withPrecedence() {
	tab := table.FromSlices(
		[]string{"Person", "Height"},
		[][]string{
			[]string{"Foo", "0"},
			[]string{"Bar", "0"},
		})

	s, _ := Infer(
		tab,
		WithPriorityOrder([]FieldType{NumberType, BooleanType, YearType, IntegerType, GeoPointType, YearMonthType, DateType, DateTimeType, TimeType, DurationType, ArrayType, ObjectType}))
	fmt.Println("Fields:")
	for _, f := range s.Fields {
		fmt.Printf("{Name:%s Type:%s Format:%s}\n", f.Name, f.Type, f.Format)
	}
	// Output: Fields:
	// {Name:Person Type:string Format:default}
	// {Name:Height Type:number Format:default}
}

func ExampleInferImplicitCasting() {
	tab := table.FromSlices(
		"ExampleTable",
		[]string{"Person", "Height"},
		[][]string{
			[]string{"Foo", "5"},
			[]string{"Bar", "4"},
			[]string{"Bez", "5.5"},
		})
	s, _ := InferImplicitCasting(tab)
	fmt.Println("Fields:")
	for _, f := range s.Fields {
		fmt.Printf("{Name:%s Type:%s Format:%s}\n", f.Name, f.Type, f.Format)
	}
	// Output: Fields:
	// {Name:Person Type:string Format:default}
	// {Name:Height Type:number Format:default}
}

func TestInferSampleLimit(t *testing.T) {
	data := []struct {
		desc        string
		sampleLimit int
		headers     []string
		table       [][]string
		want        int
	}{
		{"SampleZero", 0, []string{"Age"}, [][]string{[]string{"1"}, []string{"2"}, []string{"3"}}, 3},
		{"SampleOne", 1, []string{"Age"}, [][]string{[]string{"1"}, []string{"2"}, []string{"3"}}, 1},
		{"SampleTwo", 2, []string{"Age"}, [][]string{[]string{"1"}, []string{"2"}, []string{"3"}}, 2},
		{"SampleThree", 3, []string{"Age"}, [][]string{[]string{"1"}, []string{"2"}, []string{"3"}}, 3},
		{"SampleTen", 10, []string{"Age"}, [][]string{[]string{"1"}, []string{"2"}, []string{"3"}}, 3},
		{"SampleMinusOne", -1, []string{"Age"}, [][]string{[]string{"1"}, []string{"2"}, []string{"3"}}, 3},
		{"SampleMinusTen", -10, []string{"Age"}, [][]string{[]string{"1"}, []string{"2"}, []string{"3"}}, 3},
		{"SampleEmptyZero", 0, []string{"Age"}, [][]string{}, 0},
		{"SampleEmptyOne", 1, []string{"Age"}, [][]string{}, 0},
		{"SampleEmptyMinusTen", -10, []string{"Age"}, [][]string{}, 0},
	}
	for _, d := range data {
		t.Run(d.desc, func(t *testing.T) {
			is := is.New(t)
			s, err := sample(table.FromSlices(d.desc, d.headers, d.table), &inferConfig{sampleLimit: d.sampleLimit})
			is.NoErr(err)

			is.Equal(len(s), d.want)
		})
	}
	t.Run("LimitNotSpecified", func(t *testing.T) {
		data := []struct {
			desc    string
			headers []string
			table   [][]string
			want    int
		}{
			{"SampleDefault", []string{"Age"}, [][]string{[]string{"1"}, []string{"2"}, []string{"3"}}, 3},
			{"SampleDefault", []string{"Age"}, [][]string{}, 0},
		}
		for _, d := range data {
			t.Run(d.desc, func(t *testing.T) {
				is := is.New(t)
				s, err := sample(table.FromSlices(d.desc, d.headers, d.table), &inferConfig{})
				is.NoErr(err)

				is.Equal(len(s), d.want)
			})
		}
	})
}

func TestInfer(t *testing.T) {
	data := []struct {
		desc    string
		headers []string
		table   [][]string
		want    Schema
	}{
		{"1Cell_Date", []string{"Birthday"}, [][]string{[]string{"1983-10-15"}}, Schema{Fields: []Field{{Name: "Birthday", Type: DateType, Format: defaultFieldFormat}}}},
		{"1Cell_Integer", []string{"Age"}, [][]string{[]string{"10"}}, Schema{Fields: []Field{{Name: "Age", Type: IntegerType, Format: defaultFieldFormat}}}},
		{"1Cell_Number", []string{"Weight"}, [][]string{[]string{"20.2"}}, Schema{Fields: []Field{{Name: "Weight", Type: NumberType, Format: defaultFieldFormat}}}},
		{"1Cell_Boolean", []string{"Foo"}, [][]string{[]string{"0"}}, Schema{Fields: []Field{{Name: "Foo", Type: BooleanType, Format: defaultFieldFormat}}}},
		{"1Cell_Object", []string{"Foo"}, [][]string{[]string{`{"name":"foo"}`}}, Schema{Fields: []Field{{Name: "Foo", Type: ObjectType, Format: defaultFieldFormat}}}},
		{"1Cell_Array", []string{"Foo"}, [][]string{[]string{`["name"]`}}, Schema{Fields: []Field{{Name: "Foo", Type: ArrayType, Format: defaultFieldFormat}}}},
		{"1Cell_String", []string{"Foo"}, [][]string{[]string{"name"}}, Schema{Fields: []Field{{Name: "Foo", Type: StringType, Format: defaultFieldFormat}}}},
		{"1Cell_Time", []string{"Foo"}, [][]string{[]string{"10:15:50"}}, Schema{Fields: []Field{{Name: "Foo", Type: TimeType, Format: defaultFieldFormat}}}},
		{"1Cell_YearMonth", []string{"YearMonth"}, [][]string{[]string{"2017-08"}}, Schema{Fields: []Field{{Name: "YearMonth", Type: YearMonthType, Format: defaultFieldFormat}}}},
		{"1Cell_Year", []string{"Year"}, [][]string{[]string{"2017"}}, Schema{Fields: []Field{{Name: "Year", Type: YearType, Format: defaultFieldFormat}}}},
		{"1Cell_DateTime", []string{"DateTime"}, [][]string{[]string{"2008-09-15T15:53:00+05:00"}}, Schema{Fields: []Field{{Name: "DateTime", Type: DateTimeType, Format: defaultFieldFormat}}}},
		{"1Cell_Duration", []string{"Duration"}, [][]string{[]string{"P3Y6M4DT12H30M5S"}}, Schema{Fields: []Field{{Name: "Duration", Type: DurationType, Format: defaultFieldFormat}}}},
		{"1Cell_GeoPoint", []string{"GeoPoint"}, [][]string{[]string{"90,45"}}, Schema{Fields: []Field{{Name: "GeoPoint", Type: GeoPointType, Format: defaultFieldFormat}}}},
		{"ManyCells",
			[]string{"Name", "Age", "Weight", "Bogus", "Boolean", "Boolean1"},
			[][]string{
				[]string{"Foo", "10", "20.2", "1", "1", "1"},
				[]string{"Foo", "10", "30", "1", "1", "1"},
				[]string{"Foo", "10", "30", "Daniel", "1", "2"},
			},
			Schema{Fields: []Field{
				{Name: "Name", Type: StringType, Format: defaultFieldFormat},
				{Name: "Age", Type: IntegerType, Format: defaultFieldFormat},
				{Name: "Weight", Type: IntegerType, Format: defaultFieldFormat},
				{Name: "Bogus", Type: BooleanType, Format: defaultFieldFormat},
				{Name: "Boolean", Type: BooleanType, Format: defaultFieldFormat},
				{Name: "Boolean1", Type: BooleanType, Format: defaultFieldFormat},
			}},
		},
	}
	for _, d := range data {
		t.Run(d.desc, func(t *testing.T) {
			is := is.New(t)
			s, err := infer(d.headers, d.table, orderedTypes)
			is.NoErr(err)

			sort.Sort(s.Fields)
			sort.Sort(d.want.Fields)
			is.Equal(s, &d.want)
		})
	}
	t.Run("Error", func(t *testing.T) {
		data := []struct {
			desc    string
			headers []string
			table   [][]string
		}{
			{"NotATable", []string{}, [][]string{[]string{"1"}}},
		}
		for _, d := range data {
			t.Run(d.desc, func(t *testing.T) {
				is := is.New(t)
				_, err := infer(d.headers, d.table, orderedTypes)
				is.True(err != nil)
			})
		}
	})
}

func TestInferImplicitCasting(t *testing.T) {
	data := []struct {
		desc    string
		headers []string
		table   [][]string
		want    Schema
	}{
		{"1Cell_Date", []string{"Birthday"}, [][]string{[]string{"1983-10-15"}}, Schema{Fields: []Field{{Name: "Birthday", Type: DateType, Format: defaultFieldFormat}}}},
		{"1Cell_Integer", []string{"Age"}, [][]string{[]string{"10"}}, Schema{Fields: []Field{{Name: "Age", Type: IntegerType, Format: defaultFieldFormat}}}},
		{"1Cell_Number", []string{"Weight"}, [][]string{[]string{"20.2"}}, Schema{Fields: []Field{{Name: "Weight", Type: NumberType, Format: defaultFieldFormat}}}},
		{"1Cell_Boolean", []string{"Foo"}, [][]string{[]string{"0"}}, Schema{Fields: []Field{{Name: "Foo", Type: BooleanType, Format: defaultFieldFormat}}}},
		{"1Cell_Object", []string{"Foo"}, [][]string{[]string{`{"name":"foo"}`}}, Schema{Fields: []Field{{Name: "Foo", Type: ObjectType, Format: defaultFieldFormat}}}},
		{"1Cell_Array", []string{"Foo"}, [][]string{[]string{`["name"]`}}, Schema{Fields: []Field{{Name: "Foo", Type: ArrayType, Format: defaultFieldFormat}}}},
		{"1Cell_String", []string{"Foo"}, [][]string{[]string{"name"}}, Schema{Fields: []Field{{Name: "Foo", Type: StringType, Format: defaultFieldFormat}}}},
		{"1Cell_Time", []string{"Foo"}, [][]string{[]string{"10:15:50"}}, Schema{Fields: []Field{{Name: "Foo", Type: TimeType, Format: defaultFieldFormat}}}},
		{"1Cell_YearMonth", []string{"YearMonth"}, [][]string{[]string{"2017-08"}}, Schema{Fields: []Field{{Name: "YearMonth", Type: YearMonthType, Format: defaultFieldFormat}}}},
		{"1Cell_Year", []string{"Year"}, [][]string{[]string{"2017"}}, Schema{Fields: []Field{{Name: "Year", Type: YearType, Format: defaultFieldFormat}}}},
		{"1Cell_DateTime", []string{"DateTime"}, [][]string{[]string{"2008-09-15T15:53:00+05:00"}}, Schema{Fields: []Field{{Name: "DateTime", Type: DateTimeType, Format: defaultFieldFormat}}}},
		{"1Cell_Duration", []string{"Duration"}, [][]string{[]string{"P3Y6M4DT12H30M5S"}}, Schema{Fields: []Field{{Name: "Duration", Type: DurationType, Format: defaultFieldFormat}}}},
		{"1Cell_GeoPoint", []string{"GeoPoint"}, [][]string{[]string{"90,45"}}, Schema{Fields: []Field{{Name: "GeoPoint", Type: GeoPointType, Format: defaultFieldFormat}}}},
		{"ManyCells",
			[]string{"Name", "Age", "Weight", "Bogus", "Boolean", "Int"},
			[][]string{
				[]string{"Foo", "10", "20.2", "1", "1", "1"},
				[]string{"Foo", "10", "30", "1", "1", "1"},
				[]string{"Foo", "10", "30", "Daniel", "1", "2"},
			},
			Schema{Fields: []Field{
				{Name: "Name", Type: StringType, Format: defaultFieldFormat},
				{Name: "Age", Type: IntegerType, Format: defaultFieldFormat},
				{Name: "Weight", Type: NumberType, Format: defaultFieldFormat},
				{Name: "Bogus", Type: StringType, Format: defaultFieldFormat},
				{Name: "Boolean", Type: BooleanType, Format: defaultFieldFormat},
				{Name: "Int", Type: IntegerType, Format: defaultFieldFormat},
			}},
		},
	}
	for _, d := range data {
		t.Run(d.desc, func(t *testing.T) {
			is := is.New(t)
			s, err := inferImplicitCasting(d.headers, d.table)
			is.NoErr(err)

			sort.Sort(s.Fields)
			sort.Sort(d.want.Fields)
			is.Equal(s, &d.want)
		})
	}
	t.Run("Error", func(t *testing.T) {
		data := []struct {
			desc    string
			headers []string
			table   [][]string
		}{
			{"NotATable", []string{}, [][]string{[]string{"1"}}},
		}
		for _, d := range data {
			t.Run(d.desc, func(t *testing.T) {
				is := is.New(t)
				_, err := inferImplicitCasting(d.headers, d.table)
				is.True(err != nil)
			})
		}
	})
}

var (
	benchmarkHeaders = []string{"Name", "Birthday", "Weight", "Address", "Siblings"}
	benchmarkTable   = [][]string{
		[]string{"Foo", "2015-10-12", "20.2", `{"Street":"Foo", "Number":10, "City":"New York", "State":"NY"}`, `["Foo"]`},
		[]string{"Bar", "2015-10-12", "30", `{"Street":"Foo", "Number":10, "City":"New York", "State":"NY"}`, `["Foo"]`},
		[]string{"Bez", "2015-10-12", "30", `{"Street":"Foo", "Number":10, "City":"New York", "State":"NY"}`, `["Foo"]`},
	}
)

func benchmarkinfer(growthMultiplier int, b *testing.B) {
	for n := 0; n < b.N; n++ {
		infer(benchmarkHeaders, generateBenchmarkTable(growthMultiplier), orderedTypes)
	}
}

func benchmarkInferImplicitCasting(growthMultiplier int, b *testing.B) {
	for n := 0; n < b.N; n++ {
		inferImplicitCasting(benchmarkHeaders, generateBenchmarkTable(growthMultiplier))
	}
}

func generateBenchmarkTable(growthMultiplier int) [][]string {
	var t [][]string
	for i := 0; i < growthMultiplier; i++ {
		t = append(t, benchmarkTable...)
	}
	return t
}

func BenchmarkInferSmall(b *testing.B)                 { benchmarkinfer(1, b) }
func BenchmarkInferMedium(b *testing.B)                { benchmarkinfer(100, b) }
func BenchmarkInferBig(b *testing.B)                   { benchmarkinfer(1000, b) }
func BenchmarkInferImplicitCastingSmall(b *testing.B)  { benchmarkInferImplicitCasting(1, b) }
func BenchmarkInferImplicitCastingMedium(b *testing.B) { benchmarkInferImplicitCasting(100, b) }
func BenchmarkInferImplicitCastingBig(b *testing.B)    { benchmarkInferImplicitCasting(1000, b) }
