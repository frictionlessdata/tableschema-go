package schema

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/matryer/is"
)

func ExampleField_Cast() {
	in := `{
		"name": "id",
		"type": "string",
		"format": "default",
		"constraints": {
			"required": true,
			"minLen": "5",
			"maxLen": "10",
			"pattern": ".*11$",
			"enum":["1234511"]
		}
	}`
	var field Field
	json.Unmarshal([]byte(in), &field)
	v, err := field.Cast("1234511")
	if err != nil {
		panic(err)
	}
	fmt.Println(v)
	// Output: 1234511
}

func TestDefaultValues(t *testing.T) {
	data := []struct {
		Desc  string
		JSON  string
		Field Field
	}{
		{
			"Default Values",
			`{"name":"n1"}`,
			Field{Name: "n1", Type: defaultFieldType, Format: defaultFieldFormat, TrueValues: defaultTrueValues, FalseValues: defaultFalseValues,
				DecimalChar: defaultDecimalChar, GroupChar: defaultGroupChar, BareNumber: defaultBareNumber},
		},
		{
			"Overrinding default values",
			`{"name":"n2","type":"t2","format":"f2","falseValues":["f2"],"trueValues":["t2"]}`,
			Field{Name: "n2", Type: "t2", Format: "f2", TrueValues: []string{"t2"}, FalseValues: []string{"f2"},
				DecimalChar: defaultDecimalChar, GroupChar: defaultGroupChar, BareNumber: defaultBareNumber},
		},
	}
	for _, d := range data {
		t.Run(d.Desc, func(t *testing.T) {
			is := is.New(t)
			var f Field
			is.NoErr(json.Unmarshal([]byte(d.JSON), &f))
			is.Equal(f, d.Field)
		})
	}
}

func TestField_Cast(t *testing.T) {
	data := []struct {
		Desc     string
		Value    string
		Field    Field
		Expected interface{}
	}{
		{"Integer", "42", Field{Type: IntegerType}, int64(42)},
		{"String_URI", "http:/frictionlessdata.io", Field{Type: StringType, Format: "uri"}, "http:/frictionlessdata.io"},
		{"Boolean_TrueValues", "1", Field{Type: BooleanType, TrueValues: []string{"1"}}, true},
		{"Boolean_FalseValues", "0", Field{Type: BooleanType, FalseValues: []string{"0"}}, false},
		{"Number", "42.5", Field{Type: NumberType}, 42.5},
		{"Date_NoFormat", "2015-10-15", Field{Type: DateType}, time.Date(2015, time.October, 15, 0, 0, 0, 0, time.UTC)},
		{"Date_DefaultFormat", "2015-10-15", Field{Type: DateType, Format: defaultFieldFormat}, time.Date(2015, time.October, 15, 0, 0, 0, 0, time.UTC)},
		{"Date_CustomFormat", "15/10/2015", Field{Type: DateType, Format: "%d/%m/%Y"}, time.Date(2015, time.October, 15, 0, 0, 0, 0, time.UTC)},
		{"Time_NoFormat", "10:10:10", Field{Type: TimeType}, time.Date(0000, time.January, 01, 10, 10, 10, 00, time.UTC)},
		{"Time_DefaultFormat", "10:10:10", Field{Type: TimeType, Format: defaultFieldFormat}, time.Date(0000, time.January, 01, 10, 10, 10, 00, time.UTC)},
		{"Time_CustomFormat", "10-10-10", Field{Type: TimeType, Format: "%H-%M-%S"}, time.Date(0000, time.January, 01, 10, 10, 10, 00, time.UTC)},
		{"YearMonth", "2017-08", Field{Type: YearMonthType}, time.Date(2017, time.August, 01, 00, 00, 00, 00, time.UTC)},
		{"Year", "2017", Field{Type: YearType}, time.Date(2017, time.January, 01, 00, 00, 00, 00, time.UTC)},
		{"DateTime_NoFormat", "2008-09-15T10:53:00Z", Field{Type: DateTimeType}, time.Date(2008, time.September, 15, 10, 53, 00, 00, time.UTC)},
		{"DateTime_DefaultFormat", "2008-09-15T10:53:00Z", Field{Type: DateTimeType, Format: defaultFieldFormat}, time.Date(2008, time.September, 15, 10, 53, 00, 00, time.UTC)},
		{"Duration", "P2H", Field{Type: DurationType}, 2 * time.Hour},
		{"GeoPoint", "90,45", Field{Type: GeoPointType}, GeoPoint{90, 45}},
		{"Any", "10", Field{Type: AnyType}, "10"},
	}
	for _, d := range data {
		t.Run(d.Desc, func(t *testing.T) {
			is := is.New(t)
			c, err := d.Field.Cast(d.Value)
			is.NoErr(err)
			is.Equal(c, d.Expected)
		})
	}
	t.Run("Object_Success", func(t *testing.T) {
		is := is.New(t)
		f := Field{Type: ObjectType}
		obj, err := f.Cast(`{"name":"foo"}`)
		is.NoErr(err)

		objMap, ok := obj.(map[string]interface{})
		is.True(ok)
		is.Equal(len(objMap), 1)
		is.Equal(objMap["name"], "foo")
	})
	t.Run("Object_Failure", func(t *testing.T) {
		is := is.New(t)
		f := Field{Type: ObjectType}
		_, err := f.Cast(`{"name"}`)
		is.True(err != nil)
	})
	t.Run("Array_Success", func(t *testing.T) {
		is := is.New(t)
		f := Field{Type: ArrayType}
		obj, err := f.Cast(`["foo"]`)
		is.NoErr(err)

		arr, ok := obj.([]interface{})
		is.True(ok)
		is.Equal(len(arr), 1)
		is.Equal(arr[0], "foo")
	})
	t.Run("Array_Failure", func(t *testing.T) {
		is := is.New(t)
		f := Field{Type: ArrayType}
		_, err := f.Cast(`{"name":"foo"}`)
		is.True(err != nil)
	})
	t.Run("InvalidDate", func(t *testing.T) {
		data := []struct {
			desc  string
			field Field
			value string
		}{
			{"InvalidFormat_Any", Field{Type: DateType, Format: "any"}, "2015-10-15"},
			{"InvalidFormat_Strftime", Field{Type: DateType, Format: "Fooo"}, "2015-10-15"},
		}
		for _, d := range data {
			t.Run(d.desc, func(t *testing.T) {
				is := is.New(t)
				_, err := d.field.Cast(d.value)
				is.True(err != nil)
			})
		}
	})
	t.Run("InvalidFieldType", func(t *testing.T) {
		is := is.New(t)
		f := Field{Type: "invalidType"}
		_, err := f.Cast("42")
		is.True(err != nil)
	})
	t.Run("Constraints", func(t *testing.T) {
		t.Run("Required", func(t *testing.T) {
			is := is.New(t)
			f := Field{Type: StringType, Constraints: Constraints{Required: true}, MissingValues: map[string]struct{}{"NA": struct{}{}}}
			_, err := f.Cast("NA")
			is.True(err != nil)
		})
		t.Run("Enum", func(t *testing.T) {
			data := []struct {
				desc  string
				field Field
				value string
			}{
				{
					"SimpleCase",
					Field{Type: IntegerType, Constraints: Constraints{rawEnum: map[string]struct{}{"1": struct{}{}}}},
					"1",
				},
				{
					"NilEnumList",
					Field{Type: IntegerType},
					"10",
				},
				{
					"EmptyEnumList",
					Field{Type: IntegerType, Constraints: Constraints{rawEnum: map[string]struct{}{}}},
					"10",
				},
			}
			for _, d := range data {
				t.Run(d.desc, func(t *testing.T) {
					is := is.New(t)
					_, err := d.field.Cast(d.value)
					is.NoErr(err)
				})
			}
		})
		t.Run("EnumError", func(t *testing.T) {
			data := []struct {
				desc  string
				field Field
				value string
			}{
				{"NonEmptyEnumList", Field{Type: IntegerType, Constraints: Constraints{rawEnum: map[string]struct{}{"8": struct{}{}, "9": struct{}{}}}}, "10"},
			}
			for _, d := range data {
				t.Run(d.desc, func(t *testing.T) {
					is := is.New(t)
					_, err := d.field.Cast(d.value)
					is.True(err != nil)
				})
			}
		})
	})
}

func TestUnmarshalJSON_InvalidField(t *testing.T) {
	is := is.New(t)
	var f Field
	is.True(json.Unmarshal([]byte("{Foo:1}"), &f) != nil)
}

func TestTestString(t *testing.T) {
	is := is.New(t)
	f := Field{Type: "integer"}
	is.True(f.TestString("42"))
	is.True(!f.TestString("boo"))
}

func TestField_Uncast(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		data := []struct {
			desc  string
			field Field
			value interface{}
			want  string
		}{
			{"Int", Field{Type: IntegerType}, 1, "1"},
			{"Number", Field{Type: NumberType}, 1.0, "1"},
			{"IntNumberImplicitCast", Field{Type: NumberType}, 100, "100"},
			{"NumberToIntImplicitCast", Field{Type: IntegerType}, 100.5, "100"},
			{"Boolean", Field{Type: BooleanType}, true, "true"},
			{"Duration", Field{Type: DurationType}, 1 * time.Second, "P0Y0M0DT1S"},
			{"GeoPoint", Field{Type: GeoPointType}, "10,10", "10,10"},
			{"String", Field{Type: StringType}, "foo", "foo"},
			{"Array", Field{Type: ArrayType}, []string{"foo"}, "[foo]"},
			{"Date", Field{Type: DateType}, time.Unix(1, 0), "1970-01-01T00:00:01Z"},
			{"Year", Field{Type: YearType}, time.Unix(1, 0), "1970-01-01T00:00:01Z"},
			{"YearMonth", Field{Type: YearMonthType}, time.Unix(1, 0), "1970-01-01T00:00:01Z"},
			{"DateTime", Field{Type: DateTimeType}, time.Unix(1, 0), "1970-01-01T00:00:01Z"},
			{"Date", Field{Type: DateType}, time.Unix(1, 0), "1970-01-01T00:00:01Z"},
			{"Object", Field{Type: ObjectType}, eoStruct{Name: "Foo"}, `{"name":"Foo"}`},
			{"Any", Field{Type: AnyType}, "10", "10"},
		}
		for _, d := range data {
			t.Run(d.desc, func(t *testing.T) {
				is := is.New(t)
				got, err := d.field.Uncast(d.value)
				is.NoErr(err)
				is.Equal(d.want, got)
			})
		}
	})
	t.Run("Error", func(t *testing.T) {
		data := []struct {
			desc  string
			field Field
			value interface{}
		}{
			{"StringToIntCast", Field{Type: IntegerType}, "1.5"},
			{"StringToNumberCast", Field{Type: NumberType}, "1.5"},
			{"InvalidType", Field{Type: "Boo"}, "1"},
		}
		for _, d := range data {
			t.Run(d.desc, func(t *testing.T) {
				is := is.New(t)
				_, err := d.field.Uncast(d.value)
				is.True(err != nil)
			})
		}
	})
}
