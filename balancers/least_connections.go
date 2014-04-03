package balancers

import "fmt"
import "math"
import "math/rand"
import "github.com/tobz/cartwheel/interfaces"

func init() {
	RegisterBalancer("least_connections", NewLeastConnectionsBalancer)
}

type leastConnectionsBalancer struct {
}

func NewLeastConnectionsBalancer() *lastConnectionsBalancer {
	return &leastConnectionsBalancer{}
}

func (lcb *leastConnectionsBalancer) Choose(backends []interfaces.Backend) (interfaces.Backend, error) {
	leastConnsIndex := -1
	leastConnsCount := math.MaxUint32

	for i, backend := range backends {
		if backend.RequestsInFlight() < leastConnsCount {
			leastConnsIndex = i
		}
	}

	if leastConnsIndex == -1 {
		return nil, fmt.Errorf("couldn't find backend with lowest connection count: index is -1")
	}

	return backends[i], nil
}
