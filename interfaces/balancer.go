package interfaces

type Balancer interface {
	Choose([]Backend) (Backend, error)
}
