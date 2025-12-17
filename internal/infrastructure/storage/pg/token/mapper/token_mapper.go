package pg_token_mapper

import (
	"backend/core/pkg/pgscan"
	"backend/internal/domain/token/entity"
	"backend/internal/infrastructure/storage/model"
)

func ToEntity(row pgscan.Scannable) (token_entity.Token, error) {
	var token token_entity.Token
	err := row.Scan(
		&token.JTI,
		&token.UserID,
		&token.Token,
		&token.CreatedAt,
	)

	return token, err
}

func ToModel(token token_entity.Token) model.TokenModel {
	return model.TokenModel{
		JTI:       token.JTI,
		UserID:    token.UserID,
		Token:     token.Token,
		CreatedAt: token.CreatedAt,
	}
}
