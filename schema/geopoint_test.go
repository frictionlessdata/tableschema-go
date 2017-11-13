package schema

import (
	"testing"

	"github.com/matryer/is"
)

func TestCastGeoPoint(t *testing.T) {
	data := []struct {
		desc   string
		format string
		value  string
		want   GeoPoint
	}{
		{"DefaultNoParentheses", defaultFieldFormat, "90,40", GeoPoint{90, 40}},
		{"DefaultNoParenthesesFloats", defaultFieldFormat, "90.5,40.44", GeoPoint{90.5, 40.44}},
		{"DefaultNoParenthesesNegative", defaultFieldFormat, "-90.10,-40", GeoPoint{-90.10, -40}},
		{"DefaultNoParenthesesEmptyFormat", "", "90,40", GeoPoint{90, 40}},
		{"DefaultWithSpace", "", "90, 40", GeoPoint{90, 40}},
		{"DefaultWithSpaceNegative", "", "-90, -40", GeoPoint{-90, -40}},
		{"Array", GeoPointArrayFormat, "[90,40]", GeoPoint{90, 40}},
		{"ArrayFloat", GeoPointArrayFormat, "[90.5,40.44]", GeoPoint{90.5, 40.44}},
		{"ArrayNegative", GeoPointArrayFormat, "[-90.5,-40]", GeoPoint{-90.5, -40}},
		{"ArrayWithSpace", GeoPointArrayFormat, "[90, 40]", GeoPoint{90, 40}},
		{"ArrayWithSpaceNegative", GeoPointArrayFormat, "[-90, -40]", GeoPoint{-90, -40}},
		{"Object", GeoPointObjectFormat, `{"lon": 90, "lat": 45}`, GeoPoint{90, 45}},
	}
	for _, d := range data {
		t.Run(d.desc, func(t *testing.T) {
			is := is.New(t)
			got, err := castGeoPoint(d.format, d.value)
			is.NoErr(err)
			is.Equal(got, d.want)
		})
	}
	t.Run("Error", func(t *testing.T) {
		data := []struct {
			desc   string
			format string
			value  string
		}{
			{"BadJSON", GeoPointObjectFormat, ""},
			{"BadGeoPointJSON", GeoPointObjectFormat, `{"longi": 90, "lat": 45}`},
			{"BadFormat", "badformat", `{"longi": 90, "lat": 45}`},
			{"InvalidDefault", defaultFieldFormat, "/10,10/"},
			{"InvalidArray", defaultFieldFormat, "/[10,10]/"},
		}
		for _, d := range data {
			t.Run(d.desc, func(t *testing.T) {
				is := is.New(t)
				_, err := castGeoPoint(d.format, d.value)
				is.True(err != nil)
			})
		}
	})
}

func TestUncastGeoPoint(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		data := []struct {
			desc   string
			format string
			value  interface{}
			want   string
		}{
			{"GeoPointObject", GeoPointObjectFormat, GeoPoint{10, 10}, "{Lon:10 Lat:10}"},
			{"GeoPointArray", GeoPointArrayFormat, "[10,10]", "[10,10]"},
			{"GeoPointDefault", defaultFieldFormat, "10,10", "10,10"},
		}
		for _, d := range data {
			t.Run(d.desc, func(t *testing.T) {
				is := is.New(t)
				got, err := uncastGeoPoint(d.format, d.value)
				is.NoErr(err)
				is.Equal(d.want, got)
			})
		}
	})
	t.Run("Error", func(t *testing.T) {
		data := []struct {
			desc   string
			format string
			value  interface{}
		}{
			{"InvalidObjectType_Object", GeoPointObjectFormat, int(10)},
			{"InvalidObjectType_Array", GeoPointArrayFormat, int(10)},
			{"InvalidArray", GeoPointArrayFormat, "10,10"},
			{"InvalidObjectType_Empty", "", int(10)},
			{"InvalidObjectType_Default", defaultFieldFormat, int(10)},
			{"InvalidDefault", defaultFieldFormat, "/10,10/"},
			{"InvalidFormat", "badFormat", int(10)},
		}
		for _, d := range data {
			t.Run(d.desc, func(t *testing.T) {
				is := is.New(t)
				_, err := uncastGeoPoint(d.format, d.value)
				is.True(err != nil)
			})
		}
	})
}
