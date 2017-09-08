package schema

import (
	"testing"
	"time"
)

func TestEncodeTime(t *testing.T) {
	data := []struct {
		desc  string
		value time.Time
		want  string
	}{
		{"SimpleDate", time.Unix(1, 0), "1970-01-01T00:00:01Z"},
	}
	for _, d := range data {
		t.Run(d.desc, func(t *testing.T) {
			got, err := encodeTime(d.value)
			if err != nil {
				t.Fatalf("err want:nil got:%q", err)
			}
			if d.want != got {
				t.Fatalf("val want:%s got:%s", d.want, got)
			}
		})
	}
}
