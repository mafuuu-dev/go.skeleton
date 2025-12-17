package revoke_tokens

import (
	"backend/core/pkg/errorsx"
	"backend/core/pkg/scope"
	"backend/internal/domain/token/repository"
	"backend/internal/infrastructure/storage/pg/token"
)

type repositories struct {
	Token token_repository.Repository
}

type RevokeTokens struct {
	scope        *scope.Scope
	repositories repositories
}

func NewRevokeTokens(scope *scope.Scope) *RevokeTokens {
	return &RevokeTokens{
		scope: scope,
		repositories: repositories{
			Token: pg_token.NewRepository(scope.Support.Factory.Repository),
		},
	}
}

func (s *RevokeTokens) Run() {
	if err := s.repositories.Token.RevokeExpiredTokens(); err != nil {
		s.scope.Log.Warnf("Error revoke expired tokens: %v", errorsx.JSONTrace(err))
		return
	}

	s.scope.Log.Info("Expired tokens revoked")
}
