package transport

import (
	"backend/core/pkg/centrifugo"
	"context"

	"go.uber.org/zap"
)

type Transport struct {
	ws       *centrifugo.Centrifugo
	bulk     *centrifugo.BulkPublish
	presence *centrifugo.Presence
}

func NewTransport(ctx context.Context, log *zap.SugaredLogger, ws *centrifugo.Centrifugo) *Transport {
	return &Transport{
		ws:       ws,
		bulk:     centrifugo.NewBulkPublish(ctx, ws),
		presence: centrifugo.NewPresence(ctx, log, ws),
	}
}

func (t *Transport) WS() *centrifugo.Centrifugo {
	return t.ws
}

func (t *Transport) Bulk() *centrifugo.BulkPublish {
	return t.bulk
}

func (t *Transport) Presence() *centrifugo.Presence {
	return t.presence
}
