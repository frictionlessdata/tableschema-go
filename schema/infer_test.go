package schema

import (
	"fmt"
	"sort"
	"testing"

	"github.com/matryer/is"

	"github.com/frictionlessdata/tableschema-go/table"
)

func Exampleinfer() {
	tab := table.FromSlices(
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

func ExampleInferImplicitCasting() {
	tab := table.FromSlices(
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
			s, err := infer(d.headers, d.table)
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
				_, err := infer(d.headers, d.table)
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
		infer(benchmarkHeaders, generateBenchmarkTable(growthMultiplier))
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
