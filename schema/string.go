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

func castString(format, value string) (string, error) {
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
