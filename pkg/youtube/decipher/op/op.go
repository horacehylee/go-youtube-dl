package op

import (
	"fmt"
	"regexp"
)

// DecryptOp manipulates byte slice for decryption
type DecryptOp func(b []byte) []byte

// DecryptOpFunc is a function with parameter for a decrypt operation.
type DecryptOpFunc func(p interface{}) DecryptOp

type findFunctionNameFunc func(b []byte) (string, error)

type decryptOpFuncProvider struct {
	Name                 string
	FindFunctionNameFunc findFunctionNameFunc
	DecryptOpFunc        DecryptOpFunc
}

func findFunctionNameRegex(name string, regex *regexp.Regexp) findFunctionNameFunc {
	return func(b []byte) (string, error) {
		matches := regex.FindSubmatch(b)
		if matches == nil || len(matches) < 2 {
			return "", fmt.Errorf("failed to find %v decrypt function", name)
		}
		name := string(matches[1])
		return name, nil
	}
}
