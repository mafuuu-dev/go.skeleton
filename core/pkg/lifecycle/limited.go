package lifecycle

import (
	"golang.org/x/sync/errgroup"
)

type LimitedGroup struct {
	group errgroup.Group
	sem   chan struct{}
}

func NewLimitedGroup(max int) *LimitedGroup {
	if max <= 0 {
		max = 1
	}

	return &LimitedGroup{
		sem: make(chan struct{}, max),
	}
}

func (lg *LimitedGroup) Go(f func() error) {
	lg.group.Go(func() error {
		lg.sem <- struct{}{}
		defer func() { <-lg.sem }()

		return f()
	})
}

func (lg *LimitedGroup) Wait() error {
	return lg.group.Wait()
}
