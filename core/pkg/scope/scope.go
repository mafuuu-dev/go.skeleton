package scope

import (
	"backend/core/config"
	"backend/core/pkg/centrifugo"
	"backend/core/pkg/memory"
	"backend/core/pkg/pguof"
	"backend/core/pkg/repository"
	"backend/core/pkg/storage"
	"backend/core/pkg/transport"
	"context"

	"go.uber.org/zap"
)

type Factory struct {
	Memory     *memory.Factory
	Transport  *transport.Factory
	Repository *repository.Factory
}

type Support struct {
	Factory    *Factory
	UnitOfWork *pguof.UnitOfWork
}

type Scope struct {
	Context context.Context
	Config  *config.Config
	DB      *storage.Storage
	Socket  *centrifugo.Centrifugo
	Log     *zap.SugaredLogger
	Support *Support
}

func New(
	ctx context.Context,
	cfg *config.Config,
	logger *zap.SugaredLogger,
	database *storage.Storage,
	socket *centrifugo.Centrifugo,
	support Support,
) *Scope {
	return &Scope{
		Context: ctx,
		Config:  cfg,
		DB:      database,
		Socket:  socket,
		Log:     logger,
		Support: &support,
	}
}
