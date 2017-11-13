package schema

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
)

// Formats specific to GeoPoint field type.
const (
	GeoPointArrayFormat  = "array"
	GeoPointObjectFormat = "object"
)

// GeoPoint represents a "geopoint" cell.
// More at: https://specs.frictionlessdata.io/table-schema/#geopoint
type GeoPoint struct {
	Lon float64 `json:"lon,omitempty"`
	Lat float64 `json:"lat,omitempty"`
}

// UnmarshalJSON sets *f to a copy of data. It will respect the default values
func (p *GeoPoint) UnmarshalJSON(data []byte) error {
	type geoPointAlias struct {
		Lon *float64 `json:"lon,omitempty"`
		Lat *float64 `json:"lat,omitempty"`
	}
	var a geoPointAlias
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}
	if a.Lon == nil || a.Lat == nil {
		return fmt.Errorf("Invalid geopoint:\"%s\"", string(data))
	}
	p.Lon = *a.Lon
	p.Lat = *a.Lat
	return nil
}

var (
	geoPointDefaultRegexp = regexp.MustCompile(`^([-+]?[0-9]*\.?[0-9]*), ?([-+]?[0-9]*\.?[0-9]*)$`)
	geoPointArrayRegexp   = regexp.MustCompile(`^\[([-+]?[0-9]*\.?[0-9]+), ?([-+]?[0-9]*\.?[0-9]+)\]$`)
)

func castGeoPoint(format, value string) (GeoPoint, error) {
	switch format {
	case "", defaultFieldFormat:
		return applyGeoPointRegexp(geoPointDefaultRegexp, value)
	case GeoPointArrayFormat:
		return applyGeoPointRegexp(geoPointArrayRegexp, value)
	case GeoPointObjectFormat:
		var p GeoPoint
		if err := json.Unmarshal([]byte(value), &p); err != nil {
			return GeoPoint{}, err
		}
		return p, nil
	}
	return GeoPoint{}, fmt.Errorf("invalid geopoint format:%s", format)
}

func applyGeoPointRegexp(r *regexp.Regexp, value string) (GeoPoint, error) {
	matches := r.FindStringSubmatch(value)
	if len(matches) == 0 || len(matches[1]) == 0 || len(matches[2]) == 0 {
		return GeoPoint{}, fmt.Errorf("Invalid geopoint:\"%s\"", value)
	}
	lon, _ := strconv.ParseFloat(matches[1], 64)
	lat, _ := strconv.ParseFloat(matches[2], 64)
	return GeoPoint{lon, lat}, nil
}

func uncastGeoPoint(format string, gp interface{}) (string, error) {
	switch format {
	case "", defaultFieldFormat:
		value, ok := gp.(string)
		if ok {
			_, err := applyGeoPointRegexp(geoPointDefaultRegexp, value)
			if err != nil {
				return "", err
			}
			return value, nil
		}
		return "", fmt.Errorf("invalid object type to uncast to geopoint dfault format. want:string got:%v", reflect.TypeOf(gp).String())
	case GeoPointArrayFormat:
		value, ok := gp.(string)
		if ok {
			_, err := applyGeoPointRegexp(geoPointArrayRegexp, value)
			if err != nil {
				return "", err
			}
			return value, nil
		}
		return "", fmt.Errorf("invalid object type to uncast to geopoint %s format. want:string got:%v", GeoPointArrayFormat, reflect.TypeOf(gp).String())
	case GeoPointObjectFormat:
		value, ok := gp.(GeoPoint)
		if ok {
			return fmt.Sprintf("%+v", value), nil
		}
		return "", fmt.Errorf("invalid object type to uncast to geopoint %s format. want:schema.Geopoint got:%v", GeoPointObjectFormat, reflect.TypeOf(gp).String())
	}
	return "", fmt.Errorf("invalid geopoint - type:%v value:\"%v\" format:%s", gp, reflect.ValueOf(gp).Type(), format)
}
