package core


type Scheduler interface {
	OnStart() error
	OnError () error
	Name() string
}