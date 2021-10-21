// Package runutil includes helpers for scheduling runnable, periodic functions.
package runutil

import (
	"context"
	"reflect"
	"runtime"
	"time"

	log "github.com/sirupsen/logrus"
)

// RunEvery runs the provided command periodically.
// It runs in a goroutine, and can be cancelled by finishing the supplied context.
func RunEvery(ctx context.Context, period time.Duration, f func()) {
	funcName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	ticker := time.NewTicker(period)
	go func() {
		for {
			select {
			case <-ticker.C:
				//log.WithField("function", funcName).Trace("running")
				f()
			case <-ctx.Done():
				log.WithField("function", funcName).Debug("ticker context is closed, exiting")
				ticker.Stop()
				return
			}
		}
	}()
}

// RunOnce delays running the supplied command.
// It runs in goroutine and can be cancelled by completing the provided context.
func RunOnce(ctx context.Context, period time.Duration, f func()) {
	funcName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	timer := time.NewTimer(period)
	go func() {
		for {
			select {
			case <-timer.C:
				//log.WithField("function", funcName).Trace("running")
				f()
			case <-ctx.Done():
				log.WithField("function", funcName).Debug("timer context is closed, exiting")
				timer.Stop()
				return
			}
		}
	}()
}
