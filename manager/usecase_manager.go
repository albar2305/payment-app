package manager

import (
	"github.com/albar2305/payment-app/usecase"
)

type UseCaseManager interface {
	UserUseCase() usecase.UserUseCase
	CustomerUseCase() usecase.CustomerUseCase
	MerchantUseCase() usecase.MerchantUseCase
	TransactionUseCase() usecase.TransactionUseCase
}

type useCaseManager struct {
	repoManager RepoManager
}

// TransactionUseCase implements UseCaseManager.
func (u *useCaseManager) TransactionUseCase() usecase.TransactionUseCase {
	return usecase.NewTransactionUseCase(u.repoManager.TransactionRepo(), u.UserUseCase(), u.CustomerUseCase(), u.MerchantUseCase())
}

// MerchantUseCase implements UseCaseManager.
func (u *useCaseManager) MerchantUseCase() usecase.MerchantUseCase {
	return usecase.NewMerchantUseCase(u.repoManager.MerchantRepo())
}

// CustomerUseCase implements UseCaseManager.
func (u *useCaseManager) CustomerUseCase() usecase.CustomerUseCase {
	return usecase.NewCustomerUseCase(u.repoManager.CustomerRepo(), u.UserUseCase())
}

// UserUseCase implements UseCaseManager.
func (u *useCaseManager) UserUseCase() usecase.UserUseCase {
	return usecase.NewUserUseCase(u.repoManager.UserRepo())
}

func NewUseCaseManager(repoManager RepoManager) UseCaseManager {
	return &useCaseManager{repoManager: repoManager}
}
