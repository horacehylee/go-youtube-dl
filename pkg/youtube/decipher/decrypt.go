package decipher

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/horacehylee/go-youtube-dl/pkg/youtube/decipher/op"
)

func (d *Decipher) decryptSignature(videoID string, sig string) (string, error) {
	ops, err := d.getDecryptOps(videoID)
	if err != nil {
		return "", err
	}
	b := []byte(sig)
	for _, op := range ops {
		b = op(b)
	}
	return string(b), nil
}

func (d *Decipher) getDecryptOps(videoID string) ([]op.DecryptOp, error) {
	var ops []op.DecryptOp

	r, err := d.client.PlayerJS(videoID)
	if err != nil {
		return ops, err
	}
	defer r.Close()

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return ops, err
	}

	registry := newDecryptOpRegistry()
	registry.register(op.ReverseOpFuncProvider{}, b)
	registry.register(op.SwapOpFuncProvider{}, b)
	registry.register(op.SpliceOpFuncProvider{}, b)
	if registry.err != nil {
		return ops, registry.err
	}

	opsStrings, err := getDecryptOpStrings(b)
	if err != nil {
		return ops, err
	}

	ops = make([]op.DecryptOp, len(opsStrings))
	for i, opsString := range opsStrings {
		f, err := parseJsFunction(opsString)
		if err != nil {
			return ops, err
		}
		opsFunc, ok := registry.get(f.name)
		if !ok {
			return ops, fmt.Errorf("ops func cannot be found: %v", f.name)
		}
		ops[i] = opsFunc(f.param)
	}
	return ops, nil
}

var decryptOpStringsPattern = regexp.MustCompile(`function\(a\){a=a\.split\(""\);(.*);return a\.join\(""\)}`)

// getDecryptOpStrings splits function calls with ';' to list of calls
func getDecryptOpStrings(b []byte) ([]string, error) {
	matches := decryptOpStringsPattern.FindSubmatch(b)
	if matches == nil || len(matches) < 2 {
		return []string{}, fmt.Errorf("failed to find decrypt ops with pattern: %v", decryptOpStringsPattern)
	}
	opsStrings := strings.Split(string(matches[1]), ";")
	if len(opsStrings) == 0 {
		return opsStrings, fmt.Errorf("empty decrypt ops")
	}
	return opsStrings, nil
}

type decryptOpRegistry struct {
	// registry with key for function call name, and value as the decrypt op func
	registry map[string]op.DecryptOpFunc
	err      error
}

func newDecryptOpRegistry() *decryptOpRegistry {
	return &decryptOpRegistry{
		registry: make(map[string]op.DecryptOpFunc),
	}
}

func (r *decryptOpRegistry) register(p op.DecryptOpFuncProvider, b []byte) {
	if r.err != nil {
		return
	}
	regex := p.Regex()
	matches := regex.FindSubmatch(b)
	if matches == nil || len(matches) < 2 {
		r.err = fmt.Errorf("failed to find decrypt function with pattern: %v", regex.String())
	}
	key := string(matches[1])
	r.registry[key] = p.Provide()
}

func (r *decryptOpRegistry) get(key string) (op.DecryptOpFunc, bool) {
	opFunc, found := r.registry[key]
	return opFunc, found
}
