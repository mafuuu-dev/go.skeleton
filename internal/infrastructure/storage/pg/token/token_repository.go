package pg_token

import (
	"backend/core/pkg/repository"
	"backend/internal/infrastructure/storage/pg/token/query"
)

type TokenRepository struct {
	*repository.Repository
}

func NewRepository(factory *repository.Factory) *TokenRepository {
	return &TokenRepository{
		Repository: factory.Instance(),
	}
}

func (r *TokenRepository) RevokeExpiredTokens() error {
	return pg_token_query.NewRevokeExpiredTokens(r.Query()).Execute()
}
