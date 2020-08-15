package op

import "fmt"

// DecryptOpRegistry with key for function call name, and value as the decrypt op func
type DecryptOpRegistry struct {
	registry  map[string]DecryptOpFunc
	providers []*decryptOpFuncProvider
}

// DefaultDecryptOpRegistry for default set of decrypt ops
func DefaultDecryptOpRegistry() *DecryptOpRegistry {
	return newDecryptOpRegistry(
		reverseOpFuncProvider,
		spliceOpFuncProvider,
		swapOpFuncProvider,
	)
}

func newDecryptOpRegistry(providers ...*decryptOpFuncProvider) *DecryptOpRegistry {
	return &DecryptOpRegistry{
		registry:  make(map[string]DecryptOpFunc, len(providers)),
		providers: providers,
	}
}

// Load registry with Javascript byte slice
func (r *DecryptOpRegistry) Load(b []byte) error {
	for _, p := range r.providers {
		name, err := p.FindFunctionNameFunc(b)
		if err != nil {
			return err
		}
		if _, found := r.registry[name]; found {
			return fmt.Errorf(`same name "%v" is used already, cannot be used for %v`, name, p.Name)
		}
		r.registry[name] = p.DecryptOpFunc
	}
	return nil
}

// Get from registry with key of function name
func (r *DecryptOpRegistry) Get(name string) (DecryptOpFunc, bool) {
	opFunc, found := r.registry[name]
	return opFunc, found
}
