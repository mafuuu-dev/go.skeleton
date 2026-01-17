package lifecycle

import (
	"backend/core/pkg/errorsx"
	"sync"
	"time"

	"go.uber.org/zap"
)

type Group struct {
	wg  sync.WaitGroup
	log *zap.SugaredLogger
}

func NewGroup(log *zap.SugaredLogger) *Group {
	return &Group{
		log: log,
	}
}

func (g *Group) Go(fn func()) {
	g.wg.Add(1)

	go func() {
		defer g.wg.Done()

		defer func() {
			if r := recover(); r != nil {
				g.log.Warn(errorsx.Recover(r).(*errorsx.Error).ToJSON())
				time.Sleep(5 * time.Second)

				g.Go(fn)
			}
		}()

		fn()
	}()
}

func (g *Group) Wait() {
	g.wg.Wait()
}
