package core


type Scheduler interface {
	Start() error
	Stop() error
	Error () error
	Name() string
}