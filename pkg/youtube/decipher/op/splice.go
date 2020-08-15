package op

import "regexp"

var spliceOpRegex = regexp.MustCompile(`([a-zA-Z_\\$][a-zA-Z_0-9]*):function\(a,b\){a\.splice\(0,b\)}`)

var spliceOpFuncProvider = &decryptOpFuncProvider{
	Name:                 "splice",
	FindFunctionNameFunc: findFunctionNameRegex("splice", spliceOpRegex),
	DecryptOpFunc:        spliceOpFunc,
}

func spliceOpFunc(p interface{}) DecryptOp {
	return func(b []byte) []byte {
		return b[p.(int):]
	}
}
