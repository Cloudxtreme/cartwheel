package interfaces

type Backend interface {
	RequestsInFlight() uint32
}
