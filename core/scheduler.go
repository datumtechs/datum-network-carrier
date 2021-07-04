package core


type Scheduler interface {
	OnStart() error
	OnStop() error
	OnError () error
	Name() string
}