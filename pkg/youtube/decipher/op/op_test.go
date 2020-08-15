package op

import (
	"io/ioutil"
	"math/rand"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDecryptOps(t *testing.T) {
	type testCase struct {
		provider *decryptOpFuncProvider
		input    []byte
		param    interface{}
		expected []byte
	}
	cases := []testCase{
		{
			provider: reverseOpFuncProvider,
			input:    []byte{0x10, 0x20, 0x30, 0x40},
			param:    rand.Int(),
			expected: []byte{0x40, 0x30, 0x20, 0x10},
		},
		{
			provider: spliceOpFuncProvider,
			input:    []byte{0x10, 0x20, 0x30, 0x40},
			param:    2,
			expected: []byte{0x30, 0x40},
		},
		{
			provider: swapOpFuncProvider,
			input:    []byte{0x10, 0x20, 0x30, 0x40},
			param:    2,
			expected: []byte{0x30, 0x20, 0x10, 0x40},
		},
	}
	for _, c := range cases {
		op := c.provider.DecryptOpFunc(c.param)
		actual := op(c.input)
		if !cmp.Equal(c.expected, actual) {
			t.Error(cmp.Diff(c.expected, actual))
		}
	}
}

func TestDecryptOps_FindFunctionName(t *testing.T) {
	type testCase struct {
		provider *decryptOpFuncProvider
		expected string
	}
	cases := []testCase{
		{
			provider: reverseOpFuncProvider,
			expected: "ch",
		},
		{
			provider: spliceOpFuncProvider,
			expected: "ox",
		},
		{
			provider: swapOpFuncProvider,
			expected: "EQ",
		},
	}

	f, err := os.Open("testdata/player.txt")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		t.Error(err)
	}

	for _, c := range cases {
		name, err := c.provider.FindFunctionNameFunc(b)
		if err != nil {
			t.Error(err)
		}
		if !cmp.Equal(c.expected, name) {
			t.Error(cmp.Diff(c.expected, name))
		}
	}
}

func TestDecryptOps_UnknownJs_FindFunctionName(t *testing.T) {
	f, err := os.Open("testdata/unknown_player.txt")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		t.Error(err)
	}

	_, err = reverseOpFuncProvider.FindFunctionNameFunc(b)
	if err == nil {
		t.Error("error is not thrown")
	}
}
