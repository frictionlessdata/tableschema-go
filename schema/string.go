package schema

import (
	"fmt"
	"github.com/satori/go.uuid"
	"net/mail"
	"net/url"
)

// Valid string formats and configuration.
const (
	stringURI         = "uri"
	stringEmail       = "email"
	stringUUID        = "uuid"
	stringBinary      = "binary"
	stringUUIDVersion = 4
)

func checkStringConstraints(v string, c Constraints) error {
	minLength := c.MinLength
	maxLength := c.MaxLength
	re := c.compiledPattern

	if minLength != 0 && len(v) < minLength {
		return fmt.Errorf("constraint check error: %v %v < minimum:%v", v, len(v), minLength)
	}
	if maxLength != 0 && len(v) > maxLength {
		return fmt.Errorf("constraint check error: %v %v > maximum:%v", v, len(v), maxLength)
	}

	if re != nil && !re.MatchString(v) {
		return fmt.Errorf("constraint check error: %v don't fit pattern : %v ", v, c.Pattern)
	}
	return nil
}

func castString(format, value string, c Constraints) (string, error) {
	err := checkStringConstraints(value, c)
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
