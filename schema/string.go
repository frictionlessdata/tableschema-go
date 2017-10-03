package schema

import (
	"fmt"
	"github.com/satori/go.uuid"
	"net/mail"
	"net/url"
	"regexp"
)

// Valid string formats and configuration.
const (
	stringURI         = "uri"
	stringEmail       = "email"
	stringUUID        = "uuid"
	stringBinary      = "binary"
	stringUUIDVersion = 4
)

func checkStringConstraints(v string, minLength, maxLength int, pattern, t string) error {
	if minLength != 0 && len(v) < minLength {
		return fmt.Errorf("constraint check error: %s:%v %v < minimum:%v", t, v, len(v), minLength)
	}
	if maxLength != 0 && len(v) > maxLength {
		return fmt.Errorf("constraint check error: %s:%v %v > maximum:%v", t, v, len(v), maxLength)
	}

	if pattern != "" {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return fmt.Errorf("constraint check error: invalid pattern %v for %v : %v ", pattern, v)
		}

		match := re.MatchString(v)
		if false == match {
			return fmt.Errorf("constraint check error: %s:%v don't fit pattern : %v ", t, v, pattern)
		}
	}
	return nil
}

func decodeString(format, value string, c Constraints) (string, error) {
	err := checkStringConstraints(value, c.MinLength, c.MaxLength, c.Pattern, StringType)
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
