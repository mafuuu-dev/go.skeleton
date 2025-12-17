package scheduler

import (
	"backend/core/pkg/scope"
	"backend/internal/constants"
	"backend/internal/workers/scheduler/services/revoke_tokens"
	"time"
)

func Start(scope *scope.Scope) {
	scope.Log.Info("Scheduler loop started")

	revokeTokensTicker := time.NewTicker(internal_constants.LoopSpeedRevokeTokens)
	revokeTokens := revoke_tokens.NewRevokeTokens(scope)
	defer revokeTokensTicker.Stop()

	for {
		select {
		case <-scope.Context.Done():
			scope.Log.Info("Stopping scheduler loop..")
			return
		case <-revokeTokensTicker.C:
			revokeTokens.Run()
		}
	}
}
