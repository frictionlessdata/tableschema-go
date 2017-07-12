package table

import (
	"reflect"
	"strings"
	"testing"
)

func TestRead(t *testing.T) {
	in := `first,last
"Foo","Bar"
"Bez","Boo"`
	table := Table{strings.NewReader(in)}
	rows, err := table.Read()
	if err != nil {
		t.Errorf("err want:nil got:%q", err)
	}
	expected := [][]string{{"first", "last"}, {"Foo", "Bar"}, {"Bez", "Boo"}}
	for i := range expected {
		if !reflect.DeepEqual(rows[i], expected[i]) {
			t.Errorf("val want:%v got:%v", expected[i], rows[i])
		}
	}
}
