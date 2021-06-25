package handler

import (
	"context"
	"github.com/RosettaFlow/Carrier-Go/common/runutil"
	"sync"
)

// processes pending blocks queue on every processPendingBlocksPeriod
func (s *Service) processPendingBlocksQueue() {
	// Prevents multiple queue processing goroutines (invoked by RunEvery) from contending for data.
	locker := new(sync.Mutex)
	runutil.RunEvery(s.ctx, 1, func() {
		locker.Lock()
		if err := s.processPendingBlocks(s.ctx); err != nil {
			log.WithError(err).Debug("Could not process pending blocks")
		}
		locker.Unlock()
	})
}

// processes the block tree inside the queue
func (s *Service) processPendingBlocks(ctx context.Context) error {
	return nil
}


