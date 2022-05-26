package state

import "github.com/datumtechs/datum-network-carrier/event"

// Notifier interface defines the methods of the service that provides state updates to consumers.
type Notifier interface {
	StateFeed() *event.Feed
}
