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

	registry := op.NewDecryptOpRegistry(
		op.ReverseOpFuncProvider,
		op.SpliceOpFuncProvider,
		op.SwapOpFuncProvider,
	)
	err = registry.Load(b)
	if err != nil {
		return ops, err
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
		opsFunc, ok := registry.Get(f.Name)
		if !ok {
			return ops, fmt.Errorf("ops func cannot be found: %v", f.Name)
		}
		ops[i] = opsFunc(f.Param)
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
