package op

import "regexp"

var spliceOpRegex = regexp.MustCompile(`([a-zA-Z_\\$][a-zA-Z_0-9]*):function\(a,b\){a\.splice\(0,b\)}`)

// SpliceOpFuncProvider provide splice operations
type SpliceOpFuncProvider struct {
}

// Regex needed for identify splice operation
func (s SpliceOpFuncProvider) Regex() *regexp.Regexp {
	return spliceOpRegex
}

// Provide splice DecryptOpFunc
func (s SpliceOpFuncProvider) Provide() DecryptOpFunc {
	return spliceOpFunc
}

func spliceOpFunc(p interface{}) DecryptOp {
	return func(b []byte) []byte {
		return b[p.(int):]
	}
}
