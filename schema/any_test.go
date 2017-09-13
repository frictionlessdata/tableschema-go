package schema

import "testing"

func TestCastAny(t *testing.T) {
	got, err := castAny("foo")
	if err != nil {
		t.Fatalf("err want:nil got:%q", err)
	}
	want := "foo"
	if want != got {
		t.Fatalf("val want:%s got:%s", want, got)
	}
}

func TestEncodeAny(t *testing.T) {
	got, err := encodeAny(10)
	if err != nil {
		t.Fatalf("err want:nil got:%q", err)
	}
	want := "10"
	if want != got {
		t.Fatalf("val want:%s got:%s", want, got)
	}
}
