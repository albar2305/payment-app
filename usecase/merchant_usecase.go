package usecase

import (
	"github.com/albar2305/payment-app/model"
	"github.com/albar2305/payment-app/repository"
	"github.com/albar2305/payment-app/utils/common"
)

type MerchantUseCase interface {
	RegisterNewMerchant(payload model.CreateMerchantRequest) (model.Merchant, error)
	GetMerchant(id string) (model.Merchant, error)
	ListMerchant(params model.PaginationParams) ([]model.Merchant, error)
	DeleteMerchant(id string) error
}

type merchantUseCase struct {
	repo repository.MerchantRepository
}

func NewMerchantUseCase(repo repository.MerchantRepository) MerchantUseCase {
	return &merchantUseCase{
		repo: repo,
	}
}

// DeleteMerchants implements MerchantUseCase.
func (usecase *merchantUseCase) DeleteMerchant(id string) error {
	return usecase.repo.Delete(id)
}

// GetMerchants implements MerchantUseCase.
func (usecase *merchantUseCase) GetMerchant(id string) (model.Merchant, error) {
	return usecase.repo.Get(id)
}

// ListMerchants implements MerchantUseCase.
func (usecase *merchantUseCase) ListMerchant(params model.PaginationParams) ([]model.Merchant, error) {
	return usecase.repo.List(params)
}

// RegisterNewMerchants implements MerchantUseCase.
func (usecase *merchantUseCase) RegisterNewMerchant(payload model.CreateMerchantRequest) (model.Merchant, error) {
	merchant := model.Merchant{
		ID:          common.GenerateID(),
		Name:        payload.Name,
		Description: payload.Description,
		BusinesType: payload.BusinesType,
		Balance:     0,
	}

	return usecase.repo.Create(merchant)
}
