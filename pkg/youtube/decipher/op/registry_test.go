package op

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewDecryptOpRegistry(t *testing.T) {
	provider := reverseOpFuncProvider
	r := newDecryptOpRegistry(provider)

	expectedProviderCount := 1
	actualProviderCount := len(r.providers)
	if !cmp.Equal(expectedProviderCount, actualProviderCount) {
		t.Error(cmp.Diff(expectedProviderCount, actualProviderCount))
	}

	expectedProviderName := provider.Name
	actualProviderName := r.providers[0].Name
	if !cmp.Equal(expectedProviderName, actualProviderName) {
		t.Error(cmp.Diff(expectedProviderName, actualProviderName))
	}
}

func TestDefaultDecryptOpRegistry(t *testing.T) {
	r := DefaultDecryptOpRegistry()

	expectedProviderCount := 3
	actualProviderCount := len(r.providers)
	if !cmp.Equal(expectedProviderCount, actualProviderCount) {
		t.Error(cmp.Diff(expectedProviderCount, actualProviderCount))
	}
}

func TestRegistryLoadAndGet(t *testing.T) {
	r := newDecryptOpRegistry(
		reverseOpFuncProvider,
	)
	f, err := os.Open("testdata/player.txt")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		t.Error(err)
	}
	err = r.Load(b)
	if err != nil {
		t.Error(err)
	}

	_, found := r.Get("ch")
	if !found {
		t.Error(`"ch" should be found in registry`)
	}

	_, found = r.Get("EQ")
	if found {
		t.Error(`"EQ" should not be found`)
	}
}

func TestRegistryDuplicateProvider(t *testing.T) {
	r := newDecryptOpRegistry(
		reverseOpFuncProvider,
		reverseOpFuncProvider,
	)
	f, err := os.Open("testdata/player.txt")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		t.Error(err)
	}
	err = r.Load(b)
	if err == nil {
		t.Error("error not thrown")
	}
}

func TestRegistryProviderFindFunctionNameFailed(t *testing.T) {
	provider := &decryptOpFuncProvider{
		Name: "testing",
		FindFunctionNameFunc: func(b []byte) (string, error) {
			return "", errors.New("testing error")
		},
	}
	r := newDecryptOpRegistry(
		provider,
	)
	f, err := os.Open("testdata/player.txt")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		t.Error(err)
	}
	err = r.Load(b)
	if err == nil {
		t.Error("error not thrown")
	}
}
