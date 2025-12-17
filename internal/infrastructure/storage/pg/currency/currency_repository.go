package pg_currency

import (
	"backend/core/pkg/repository"
	"backend/internal/domain/currency/entity"
	"backend/internal/domain/currency/enum"
	"backend/internal/infrastructure/storage/pg/currency/query"
)

type CurrencyRepository struct {
	*repository.Repository
}

func NewRepository(factory *repository.Factory) *CurrencyRepository {
	return &CurrencyRepository{
		Repository: factory.Instance(),
	}
}

func (r *CurrencyRepository) GetCurrencyByCode(code currency_enum.Code) (*currency_entity.Currency, error) {
	return pg_currency_query.NewGetCurrencyByCode(r.Query()).
		SetCode(code).
		Execute()
}
