package consensus

type Consensus interface {
	OnStart()
	Error() error
}
