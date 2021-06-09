package node

import (
	"github.com/RosettaFlow/Carrier-Go/p2p"
)

type ServiceContext struct {

}


// ServiceConstructor is the function signature of the constructors needed to be
// registered for service instantiation.
type ServiceConstructor func(ctx *ServiceContext) (Service, error)



type Service interface {

	// Protocols retrieves the P2P protocols the service wishes to start.
	Protocols() []p2p.Protocol

	// Start is called after all services have been constructed and the networking
	// layer was also initialized to spawn any goroutines required by the service.
	Start() error

	// Stop terminates all goroutines belonging to the service, blocking until they
	// are all terminated.
	Stop() error
}
