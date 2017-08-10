package schema

import "testing"

func TestCastGeoPoint_Success(t *testing.T) {
	data := []struct {
		desc   string
		format string
		value  string
		want   GeoPoint
	}{
		{"DefaultNoParentheses", defaultFieldFormat, "90,40", GeoPoint{90, 40}},
		{"DefaultNoParenthesesEmptyFormat", "", "90,40", GeoPoint{90, 40}},
		{"DefaultWithSpace", "", "90, 40", GeoPoint{90, 40}},
		{"Array", GeoPointArrayFormat, "[90,40]", GeoPoint{90, 40}},
		{"ArrayWithSpace", GeoPointArrayFormat, "[90, 40]", GeoPoint{90, 40}},
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
}

func TestCastGeoPoint_Error(t *testing.T) {
	data := []struct {
		desc   string
		format string
		value  string
	}{
		{"BadJSON", GeoPointObjectFormat, `{"longi": 90, "lat": 45}`},
	}
	for _, d := range data {
		t.Run(d.desc, func(t *testing.T) {
			_, err := castGeoPoint(d.format, d.value)
			if err == nil {
				t.Errorf("want:err got:nil")
			}
		})
	}
}
