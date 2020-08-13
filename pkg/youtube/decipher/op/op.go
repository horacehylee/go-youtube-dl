package op

import "regexp"

// DecryptOp manipulates byte slice for decryption
type DecryptOp func(b []byte) []byte

// DecryptOpFunc returns decrypt operations with param provided
type DecryptOpFunc func(p interface{}) DecryptOp

// DecryptOpFuncProvider provides DecryptOpFunc for easy to extend to include more ops
type DecryptOpFuncProvider interface {
	Regex() *regexp.Regexp
	Provide() DecryptOpFunc
}
