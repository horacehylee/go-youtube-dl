package op

import "regexp"

var swapOpRegex = regexp.MustCompile(`([a-zA-Z_\\$][a-zA-Z_0-9]*):function\(a,b\){var c=a\[0\];a\[0\]=a\[b%a\.length\];a\[b%a\.length\]=c}`)

// SwapOpFuncProvider provide swap operations
type SwapOpFuncProvider struct {
}

// Regex needed for identify swap operation
func (s SwapOpFuncProvider) Regex() *regexp.Regexp {
	return swapOpRegex
}

// Provide swap DecryptOpFunc
func (s SwapOpFuncProvider) Provide() DecryptOpFunc {
	return swapOpFunc
}

func swapOpFunc(p interface{}) DecryptOp {
	return func(b []byte) []byte {
		pos := p.(int) % len(b)
		b[0], b[pos] = b[pos], b[0]
		return b
	}
}
