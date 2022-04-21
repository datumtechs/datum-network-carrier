package state

import "github.com/Metisnetwork/Metis-Carrier/event"

// Notifier interface defines the methods of the service that provides state updates to consumers.
type Notifier interface {
	StateFeed() *event.Feed
}
