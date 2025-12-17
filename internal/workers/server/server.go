package server

import (
	"backend/core/pkg/errorsx"
	"backend/core/pkg/lifecycle"
	"backend/core/pkg/middleware"
	"backend/core/pkg/scope"
	"backend/internal/usecases/server"
	"backend/internal/workers/server/api"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
)

func Start(scope *scope.Scope) {
	app := fiber.New(fiber.Config{
		ErrorHandler: server_usecase.NewErrorHandler(scope).Handler(),
	})

	app.Use(middleware.Logger())
	app.Use(middleware.Recovery())
	app.Use(middleware.RequestID())
	app.Use(middleware.Cors(scope.Config))
	app.Use(middleware.Compress(compress.LevelBestSpeed, []string{
		"/api/v1/maintenance/ping",
		"/api/v1/maintenance/health",
	}))

	api.Register(app, scope)

	go func() {
		if err := app.Listen(":" + scope.Config.HTTPPort); err != nil {
			scope.Log.Fatalf("Fiber server error: %v", errorsx.Error(err))
		}
	}()

	<-scope.Context.Done()

	stop(app, scope)
}

func stop(app *fiber.App, scope *scope.Scope) {
	scope.Log.Infof("Shutting down HTTP server..")

	ctx, cancel := lifecycle.Timeout(scope.Context, 5*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		scope.Log.Errorf("Error during Fiber shutdown: %v", errorsx.Error(err))
	}
}
