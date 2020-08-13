package op

// DecryptOpRegistry with key for function call name, and value as the decrypt op func
type DecryptOpRegistry struct {
	registry  map[string]DecryptOpFunc
	providers []*DecryptOpFuncProvider
}

// NewDecryptOpRegistry return new instance of registry
func NewDecryptOpRegistry(providers ...*DecryptOpFuncProvider) *DecryptOpRegistry {
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
		r.registry[name] = p.DecryptOpFunc
	}
	return nil
}

// Get from registry with key of function name
func (r *DecryptOpRegistry) Get(name string) (DecryptOpFunc, bool) {
	opFunc, found := r.registry[name]
	return opFunc, found
}
