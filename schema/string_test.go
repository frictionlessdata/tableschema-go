package schema

import "testing"

// To be in par with the python library.
func TestCastString_URIMustRequireScheme(t *testing.T) {
	if _, err := castString(stringURI, "google.com"); err == nil {
		t.Errorf("want:err got:nil")
	}
}

func TestCastString_InvalidUUIDVersion(t *testing.T) {
	// This is a uuid3: namespace DNS and python.org.
	if _, err := castString(stringUUID, "6fa459ea-ee8a-3ca4-894e-db77e160355e"); err == nil {
		t.Errorf("want:err got:nil")
	}
}

func TestCastString_Success(t *testing.T) {
	var data = []struct {
		Desc   string
		Value  string
		Format string
	}{
		{"URI", "http://google.com", stringURI},
		{"Email", "foo@bar.com", stringEmail},
		{"UUID", "C56A4180-65AA-42EC-A945-5FD21DEC0538", stringUUID},
	}
	for _, d := range data {
		v, err := castString(d.Format, d.Value)
		if err != nil {
			t.Errorf("want:nil got:%q", err)
		}
		if v != d.Value {
			t.Errorf("want:%s got:%s", d.Value, v)
		}
	}
}
