package schema

import (
	"regexp"
	"testing"

	"github.com/matryer/is"
)

// To be in par with the python library.
func TestCastString_URIMustRequireScheme(t *testing.T) {
	is := is.New(t)
	_, err := castString(stringURI, "google.com", Constraints{})
	is.True(err != nil)
}

func TestCastString_InvalidUUIDVersion(t *testing.T) {
	is := is.New(t)
	// This is a uuid3: namespace DNS and python.org.
	_, err := castString(stringUUID, "6fa459ea-ee8a-3ca4-894e-db77e160355e", Constraints{})
	is.True(err != nil)
}

func TestCastString_ErrorCheckingConstraints(t *testing.T) {
	data := []struct {
		desc        string
		value       string
		format      string
		constraints Constraints
	}{
		{"InvalidMinLength_UUID", "6fa459ea-ee8a-3ca4-894e-db77e160355e", stringUUID, Constraints{MinLength: 100}},
		{"InvalidMinLength_Email", "foo@bar.com", stringEmail, Constraints{MinLength: 100}},
		{"InvalidMinLength_URI", "http://google.com", stringURI, Constraints{MinLength: 100}},
		{"InvalidMaxLength_UUID", "6fa459ea-ee8a-3ca4-894e-db77e160355e", stringUUID, Constraints{MaxLength: 1}},
		{"InvalidMaxLength_Email", "foo@bar.com", stringEmail, Constraints{MaxLength: 1}},
		{"InvalidMaxLength_URI", "http://google.com", stringURI, Constraints{MaxLength: 1}},
		{"InvalidPattern_UUID", "6fa459ea-ee8a-3ca4-894e-db77e160355e", stringUUID, Constraints{compiledPattern: regexp.MustCompile("^[0-9a-f]{1}-.*"), Pattern: "^[0-9a-f]{1}-.*"}},
		{"InvalidPattern_Email", "foo@bar.com", stringEmail, Constraints{compiledPattern: regexp.MustCompile("[0-9].*"), Pattern: "[0-9].*"}},
		{"InvalidPattern_URI", "http://google.com", stringURI, Constraints{compiledPattern: regexp.MustCompile("^//.*"), Pattern: "^//.*"}},
	}
	for _, d := range data {
		t.Run(d.desc, func(t *testing.T) {
			is := is.New(t)
			_, err := castString(d.format, d.value, d.constraints)
			is.True(err != nil)
		})
	}
}

func TestCastString_Success(t *testing.T) {
	var data = []struct {
		desc        string
		value       string
		format      string
		constraints Constraints
	}{
		{"URI", "http://google.com", stringURI, Constraints{MinLength: 1, compiledPattern: regexp.MustCompile("^http://.*"), Pattern: "^http://.*"}},
		{"Email", "foo@bar.com", stringEmail, Constraints{MinLength: 1, compiledPattern: regexp.MustCompile(".*@.*"), Pattern: ".*@.*"}},
		{"UUID", "C56A4180-65AA-42EC-A945-5FD21DEC0538", stringUUID, Constraints{MinLength: 36, MaxLength: 36, compiledPattern: regexp.MustCompile("[0-9A-F]{8}-[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{12}"), Pattern: "[0-9A-F]{8}-[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{12}"}},
	}
	for _, d := range data {
		t.Run(d.desc, func(t *testing.T) {
			is := is.New(t)
			v, err := castString(d.format, d.value, d.constraints)
			is.NoErr(err)
			is.Equal(v, d.value)
		})
	}
}
