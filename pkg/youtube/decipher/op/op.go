package op

import (
	"fmt"
	"regexp"
)

// DecryptOp manipulates byte slice for decryption
type DecryptOp func(b []byte) []byte

// DecryptOpFunc returns decrypt operations with param provided
type DecryptOpFunc func(p interface{}) DecryptOp

type findFunctionNameFunc func(b []byte) (string, error)

// DecryptOpFuncProvider describes a DecryptOpFunc
type DecryptOpFuncProvider struct {
	Name                 string
	FindFunctionNameFunc findFunctionNameFunc
	DecryptOpFunc        DecryptOpFunc
}

func findFunctionNameRegex(name string, regex *regexp.Regexp) findFunctionNameFunc {
	return func(b []byte) (string, error) {
		matches := regex.FindSubmatch(b)
		if matches == nil || len(matches) < 2 {
			return "", fmt.Errorf("failed to find %v decrypt function with pattern: %v", name, regex.String())
		}
		name := string(matches[1])
		return name, nil
	}
}
