package op

import "regexp"

var reverseOpRegex = regexp.MustCompile(`([a-zA-Z_\\$][a-zA-Z_0-9]*):function\(a\){a\.reverse\(\)}`)

// ReverseOpFuncProvider provide reverse operations
var ReverseOpFuncProvider = &DecryptOpFuncProvider{
	Name:                 "reverse",
	FindFunctionNameFunc: findFunctionNameRegex("reverse", reverseOpRegex),
	DecryptOpFunc:        reverseOpFunc,
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
