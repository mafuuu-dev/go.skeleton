package transport

import (
	"backend/core/pkg/centrifugo"
	"context"

	"go.uber.org/zap"
)

type Factory struct {
	ctx context.Context
	log *zap.SugaredLogger
	ws  *centrifugo.Centrifugo
}

func New(ctx context.Context, log *zap.SugaredLogger, ws *centrifugo.Centrifugo) *Factory {
	return &Factory{
		ctx: ctx,
		log: log,
		ws:  ws,
	}
}

func (factory *Factory) Instance() *Transport {
	return NewTransport(factory.ctx, factory.log, factory.ws)
}
