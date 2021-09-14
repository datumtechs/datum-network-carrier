// Package feed defines evengine feed types for inter-service communication
// during a beacon node's runtime.
package feed

// How to add a new evengine to the feed:
//   1. Add a file for the new type of feed.
//   2. Add a constant describing the list of events.
//   3. Add a structure with the name `<evengine>GetData` containing any data fields that should be supplied with the evengine.
//
// Note that the same evengine is supplied to all subscribers, so the evengine received by subscribers should be considered read-only.

// EventType is the type that defines the type of evengine.
type EventType int

// Event is the evengine that is sent with operation feed updates.
type Event struct {
	// Type is the type of evengine.
	Type EventType
	// Data is evengine-specific data.
	Data interface{}
}
