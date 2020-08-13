package op

import "regexp"

var swapOpRegex = regexp.MustCompile(`([a-zA-Z_\\$][a-zA-Z_0-9]*):function\(a,b\){var c=a\[0\];a\[0\]=a\[b%a\.length\];a\[b%a\.length\]=c}`)

// SwapOpFuncProvider provide swap operations
var SwapOpFuncProvider = &DecryptOpFuncProvider{
	Name:                 "swap",
	FindFunctionNameFunc: findFunctionNameRegex("swap", swapOpRegex),
	DecryptOpFunc:        swapOpFunc,
}

func swapOpFunc(p interface{}) DecryptOp {
	return func(b []byte) []byte {
		pos := p.(int) % len(b)
		b[0], b[pos] = b[pos], b[0]
		return b
	}
}
