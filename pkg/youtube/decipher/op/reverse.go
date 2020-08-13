package op

import "regexp"

var reverseOpRegex = regexp.MustCompile(`([a-zA-Z_\\$][a-zA-Z_0-9]*):function\(a\){a\.reverse\(\)}`)

// ReverseOpFuncProvider provide reverse operations
type ReverseOpFuncProvider struct {
}

// Regex needed for identify reverse operation
func (r ReverseOpFuncProvider) Regex() *regexp.Regexp {
	return reverseOpRegex
}

// Provide reverse DecryptOpFunc
func (r ReverseOpFuncProvider) Provide() DecryptOpFunc {
	return reverseOpFunc
}

func reverseOpFunc(_ interface{}) DecryptOp {
	return func(b []byte) []byte {
		l, r := 0, len(b)-1
		for l < r {
			b[l], b[r] = b[r], b[l]
			l++
			r--
		}
		return b
	}
}
