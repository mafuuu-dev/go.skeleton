package currency_usecase

import (
	"backend/core/pkg/errorsx"
	"backend/core/pkg/scope"
	"backend/internal/domain/currency/entity"
	"backend/internal/domain/currency/enum"
	"backend/internal/domain/currency/repository"
	"backend/internal/infrastructure/storage/pg/currency"
)

type getCurrencyByCodeRepositories struct {
	Currency currency_repository.Repository
}

type GetCurrencyByCode struct {
	scope        *scope.Scope
	repositories getCurrencyByCodeRepositories
}

func NewGetCurrencyByCode(scope *scope.Scope) *GetCurrencyByCode {
	return &GetCurrencyByCode{
		scope: scope,
		repositories: getCurrencyByCodeRepositories{
			Currency: pg_currency.NewRepository(scope.Support.Factory.Repository),
		},
	}
}

func (u *GetCurrencyByCode) Get(code currency_enum.Code) (*currency_entity.Currency, error) {
	currency, err := u.repositories.Currency.GetCurrencyByCode(code)
	if err != nil {
		return nil, errorsx.Wrap(err, "Failed to get currency by code")
	}

	return currency, nil
}
