package pg_token_query

import (
	"backend/core/pkg/errorsx"
	"backend/core/pkg/query"
)

type RevokeExpiredTokens struct {
	*query.Query
}

func NewRevokeExpiredTokens(factory *query.Factory) *RevokeExpiredTokens {
	return &RevokeExpiredTokens{
		Query: factory.Instance(),
	}
}

func (q *RevokeExpiredTokens) Execute() error {
	_, err := q.Exec(q)
	return errorsx.Wrap(err, "failed to revoke expired tokens")
}

func (q *RevokeExpiredTokens) Sql() string {
	return `DELETE FROM tokens WHERE CURRENT_TIMESTAMP >= created_at + INTERVAL '24 hours'`
}
