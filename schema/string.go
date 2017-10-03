package schema

import (
	"fmt"
	"net/mail"
	"net/url"

	"github.com/satori/go.uuid"
)

// Valid string formats and configuration.
const (
	stringURI         = "uri"
	stringEmail       = "email"
	stringUUID        = "uuid"
	stringBinary      = "binary"
	stringUUIDVersion = 4
)

func checkStringConstraints(v string, minLength, maxLength int, t string) error {
	if minLength != 0 && len(v) < minLength {
		return fmt.Errorf("constraint check error: %s:%v %v < minimum:%v", t, v, minLength)
	}
	if maxLength != 0 && len(v) > maxLength {
		return fmt.Errorf("constraint check error: %s:%v %v > maximum:%v", t, v, maxLength)
	}
	return nil
}

func decodeString(format, value string, c Constraints) (string, error) {
	err := checkStringConstraints(value, c.MinLength, c.MaxLength, StringType)
	if err != nil {
		return value, err
	}

	switch format {
	case stringURI:
		_, err := url.ParseRequestURI(value)
		return value, err
	case stringEmail:
		_, err := mail.ParseAddress(value)
		return value, err
	case stringUUID:
		v, err := uuid.FromString(value)
		if v.Version() != stringUUIDVersion {
			return value, fmt.Errorf("invalid UUID version - got:%d want:%d", v.Version(), stringUUIDVersion)
		}
		return value, err
	}
	// NOTE: Returning the value for unknown format is in par with the python library.
	return value, nil
}
