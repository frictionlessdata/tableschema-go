package schema

import "testing"

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
			got, err := castGeoPoint(d.format, d.value)
			if err != nil {
				t.Errorf("want:nil got:%q", err)
			}
			if got != d.want {
				t.Errorf("want:%+v got:%+v", d.want, got)
			}
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
				_, err := castGeoPoint(d.format, d.value)
				if err == nil {
					t.Errorf("want:err got:nil")
				}
			})
		}
	})
}

func TestEncodeGeoPoint(t *testing.T) {
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
				got, err := encodeGeoPoint(d.format, d.value)
				if err != nil {
					t.Errorf("err want:nil got:%q", err)
				}
				if d.want != got {
					t.Errorf("val want:%s got:%s", d.want, got)
				}
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
				_, err := encodeGeoPoint(d.format, d.value)
				if err == nil {
					t.Errorf("want:err got:nil")
				}
			})
		}
	})
}
