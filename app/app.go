package app

import (
	"backend/core/config"
	"backend/core/constants"
	"backend/core/pkg/centrifugo"
	"backend/core/pkg/errorsx"
	"backend/core/pkg/lifecycle"
	"backend/core/pkg/logger"
	"backend/core/pkg/memory"
	"backend/core/pkg/pguof"
	"backend/core/pkg/repository"
	"backend/core/pkg/scope"
	"backend/core/pkg/storage"
	"backend/core/pkg/transport"
	"backend/internal/workers/scheduler"
	"backend/internal/workers/server"
	"sync"
	"time"
)

func Start(cfg *config.Config) {
	ctx, stop := lifecycle.Context()
	defer stop()

	shutdownCtx, shutdown := lifecycle.Timeout(ctx, 20*time.Second)
	defer shutdown()

	log := logger.New(cfg)
	database := storage.New(ctx, cfg, log)
	defer database.PG.Close()

	defer func() {
		if err := database.Redis.Close(); err != nil {
			log.Warn(errorsx.WrapJSON(err, "Redis disconnect error"))
		}
	}()

	var wg sync.WaitGroup

	socket := centrifugo.New(ctx, &wg, cfg, log, 4069)

	defer func() {
		socket.Close()

		wg.Wait()
		log.Info("Core services drained.")
	}()

	container := scope.New(
		ctx,
		cfg,
		log,
		database,
		socket,
		scope.Support{
			Factory: &scope.Factory{
				Memory:     memory.New(ctx, database.Redis),
				Transport:  transport.New(ctx, log, socket),
				Repository: repository.New(ctx, database.PG),
			},
			UnitOfWork: pguof.New(ctx, database.PG, log),
		},
	)

	workers := lifecycle.NewGroup(log)
	switch cfg.Service {
	case constants.ServiceServer:
		workers.Go(func() { server.Start(container) })
	case constants.ServiceScheduler:
		workers.Go(func() { scheduler.Start(container) })
	}

	log.Info("All services started")
	<-ctx.Done()

	log.Info("Shutdown signal received")
	<-shutdownCtx.Done()

	workers.Wait()
	log.Info("Worker services stopped.")
}
