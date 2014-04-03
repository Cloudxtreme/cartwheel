package balancers

import "fmt"
import "github.com/tobz/cartwheel/interfaces"

type balancerCreator func() interfaces.Balancer

var balancerCreators map[string]balancerCreator

func init() {
	balancerCreators = make(map[string]balancerCreator)
}

func RegisterBalancer(balancerName string, balancerCreator func() interfaces.Balancer) {
	if _, ok := balancerCreators[balancerName]; ok {
		panic(fmt.Sprintf("balancer '%s' already registered!", balancerName))
	}

	balancerCreators[balancerName] = balancerCreator
}

func CreateBalancer(balancerName string) interfaces.Balancer {
	if balancerCreator, ok := balancerCreators[balancerName]; ok {
		return balancerCreator()
	}

	panic(fmt.Sprintf("balancer '%s' does not exist!", balancerName))
}
