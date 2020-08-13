package op

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewDecryptOpRegistry(t *testing.T) {
	provider := ReverseOpFuncProvider
	r := NewDecryptOpRegistry(provider)

	expectedProviderCount := 1
	actualProviderCount := len(r.providers)
	if !cmp.Equal(expectedProviderCount, actualProviderCount) {
		t.Fatal(cmp.Diff(expectedProviderCount, actualProviderCount))
	}

	expectedProviderName := provider.Name
	actualProviderName := r.providers[0].Name
	if !cmp.Equal(expectedProviderName, actualProviderName) {
		t.Fatal(cmp.Diff(expectedProviderName, actualProviderName))
	}
}

func TestRegistryLoadAndGet(t *testing.T) {
	r := NewDecryptOpRegistry(
		ReverseOpFuncProvider,
	)
	f, err := os.Open("testdata/player.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	err = r.Load(b)
	if err != nil {
		t.Fatal(err)
	}

	_, found := r.Get("ch")
	if !found {
		t.Fatal(`"ch" should be found in registry`)
	}

	_, found = r.Get("EQ")
	if found {
		t.Fatal(`"EQ" should not be found`)
	}
}

func TestRegistryDuplicateProvider(t *testing.T) {
	r := NewDecryptOpRegistry(
		ReverseOpFuncProvider,
		ReverseOpFuncProvider,
	)
	f, err := os.Open("testdata/player.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	err = r.Load(b)
	if err == nil {
		t.Fatal("error not thrown")
	}
}

func TestRegistryProviderFindFunctionNameFailed(t *testing.T) {
	provider := &DecryptOpFuncProvider{
		Name: "testing",
		FindFunctionNameFunc: func(b []byte) (string, error) {
			return "", errors.New("testing error")
		},
	}
	r := NewDecryptOpRegistry(
		provider,
	)
	f, err := os.Open("testdata/player.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	err = r.Load(b)
	if err == nil {
		t.Fatal("error not thrown")
	}
}
