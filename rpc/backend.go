package rpc

type Backend interface {
	Start()
	Stop() error
	Status() error
}
