package balancers

import "math/rand"
import "github.com/tobz/cartwheel/interfaces"

func init() {
	RegisterBalancer("random", NewRandomBalancer)
}

type randomBalancer struct {
	randSource rand.Source
}

func NewRandomBalancer() *randomBalancer {
	return &randomBalancer{
		randSource: rand.NewSource(rand.Int63()),
	}
}

func (rb *randomBalancer) Choose(backends []interfaces.Backend) (interfaces.Backend, error) {
	i := rb.randSource.Int31n(len(backends) - 1)
	return backends[i], nil
}
