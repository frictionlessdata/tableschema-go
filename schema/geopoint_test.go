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
