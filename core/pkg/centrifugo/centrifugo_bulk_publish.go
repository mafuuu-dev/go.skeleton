package centrifugo

import (
	"backend/core/pkg/lifecycle"
	"backend/core/types"
	"context"
	"time"

	"golang.org/x/sync/errgroup"
)

type BulkMessage struct {
	Channel types.Channel
	Data    any
}

type BulkPublish struct {
	ctx        context.Context
	centrifugo *Centrifugo
}

func NewBulkPublish(ctx context.Context, centrifugo *Centrifugo) *BulkPublish {
	return &BulkPublish{
		ctx:        ctx,
		centrifugo: centrifugo,
	}
}

func (s *BulkPublish) Publish(messages []BulkMessage) error {
	ctx, cancel := lifecycle.Timeout(s.ctx, 3*time.Second)
	defer cancel()

	g, _ := errgroup.WithContext(ctx)
	for _, message := range messages {
		g.Go(func() error {
			s.centrifugo.Publish(message.Channel, message.Data)
			return nil
		})
	}

	return g.Wait()
}
