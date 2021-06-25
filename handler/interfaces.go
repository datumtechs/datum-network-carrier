package handler

// This defines the interface for interacting with block chain service
type blockchainService interface {

}

// Checker defines a struct which can verify whether a node is currently
// synchronizing a chain with the rest of peers in the network.
type Checker interface {
	Initialized() bool
	Syncing() bool
	Status() error
	Resync() error
}