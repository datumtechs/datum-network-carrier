package testing

import (
	"github.com/RosettaFlow/Carrier-Go/common/feed"
	"github.com/RosettaFlow/Carrier-Go/event"
	"sync"
)

// MockStateNotifier mocks the state notifier.
type MockStateNotifier struct {
	feed     *event.Feed
	feedLock sync.Mutex

	recv     []*feed.Event
	recvLock sync.Mutex
	recvCh   chan *feed.Event

	RecordEvents bool
}

// ReceivedEvents returns the events received by the state feed in this mock.
func (msn *MockStateNotifier) ReceivedEvents() []*feed.Event {
	msn.recvLock.Lock()
	defer msn.recvLock.Unlock()
	return msn.recv
}

// StateFeed returns a state feed.
func (msn *MockStateNotifier) StateFeed() *event.Feed {
	msn.feedLock.Lock()
	defer msn.feedLock.Unlock()

	if msn.feed == nil && msn.recvCh == nil {
		msn.feed = new(event.Feed)
		if msn.RecordEvents {
			msn.recvCh = make(chan *feed.Event)
			sub := msn.feed.Subscribe(msn.recvCh)

			go func() {
				select {
				case evt := <-msn.recvCh:
					msn.recvLock.Lock()
					msn.recv = append(msn.recv, evt)
					msn.recvLock.Unlock()
				case <-sub.Err():
					sub.Unsubscribe()
				}
			}()
		}
	}
	return msn.feed
}
