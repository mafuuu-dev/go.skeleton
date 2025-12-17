package token_repository

type Repository interface {
	RevokeExpiredTokens() error
}
