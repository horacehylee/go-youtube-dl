package decipher

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseJsFunction(t *testing.T) {
	s := "Mt.ox(a,3)"
	f, err := parseJsFunction(s)
	if err != nil {
		t.Error(err)
	}

	expected := jsFunction{
		Name:  "ox",
		Param: 3,
	}
	if !cmp.Equal(expected, f) {
		t.Error(cmp.Diff(expected, f))
	}
}

func TestParseJsFunctionFailed(t *testing.T) {
	cases := []string{
		"Mtox(a,3)",    // no dot
		"abc",          // not function call
		"Mt.ox(a,asd)", // param not number
	}
	for _, c := range cases {
		_, err := parseJsFunction(c)
		if err == nil {
			t.Error("error not thrown for", c)
		}
	}
}
