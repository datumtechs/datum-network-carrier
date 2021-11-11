package testing

import (
	"errors"
	"github.com/RosettaFlow/Carrier-Go/common/feed"
	"github.com/RosettaFlow/Carrier-Go/event"
	taskmngpb "github.com/RosettaFlow/Carrier-Go/lib/netmsg/taskmng"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/types"
	"github.com/libp2p/go-libp2p-core/peer"
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

// Sync defines a mock for the sync service.
type Sync struct {
	IsSyncing     bool
	IsInitialized bool
}

// Syncing --
func (s *Sync) Syncing() bool {
	return s.IsSyncing
}

// Initialized --
func (s *Sync) Initialized() bool {
	return s.IsInitialized
}

// Status --
func (s *Sync) Status() error {
	return nil
}

// Resync --
func (s *Sync) Resync() error {
	return nil
}

// Mock for TaskManager
type MockTaskManager struct {
}

func (m *MockTaskManager) Start() error                                               { return nil }
func (m *MockTaskManager) Stop() error                                                { return nil }
func (m *MockTaskManager) ValidateTaskResultMsg(pid peer.ID, taskResultMsg *taskmngpb.TaskResultMsg) error {
	return errors.New("invalid check")
}
func (m *MockTaskManager) OnTaskResultMsg(pid peer.ID, taskResultMsg *taskmngpb.TaskResultMsg) error {
	return nil
}
func (m *MockTaskManager) ValidateTaskResourceUsageMsg(pid peer.ID, taskResourceUsageMsg *taskmngpb.TaskResourceUsageMsg) error {
	return errors.New("invalid check")
}
func (m *MockTaskManager) OnTaskResourceUsageMsg(pid peer.ID, taskResourceUsageMsg *taskmngpb.TaskResourceUsageMsg) error {
	return nil
}
func (m *MockTaskManager) ValidateTaskTerminateMsg(pid peer.ID, terminateMsg *taskmngpb.TaskTerminateMsg) error {
	return errors.New("invalid check")
}
func (m *MockTaskManager) OnTaskTerminateMsg (pid peer.ID, terminateMsg *taskmngpb.TaskTerminateMsg) error { return nil }
func (m *MockTaskManager) SendTaskEvent(event *libtypes.TaskEvent) error     { return nil }
func (m *MockTaskManager) HandleResourceUsage(usage *types.TaskResuorceUsage) error { return nil }
