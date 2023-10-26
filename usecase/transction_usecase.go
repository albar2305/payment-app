package usecase

import (
	"fmt"

	"github.com/albar2305/payment-app/model"
	"github.com/albar2305/payment-app/repository"
	"github.com/albar2305/payment-app/utils/common"
)

type TransactionUseCase interface {
	RegisterNewTransaction(payload model.CreateTransactionRequest) (model.Transaction, error)
	GetTransactionByCustomerId(id string, params model.PaginationParams) ([]model.Transaction, error)
	ListTransaction(params model.PaginationParams) ([]model.Transaction, error)
}

type transactionUseCase struct {
	repo       repository.TransactionRepository
	userUC     UserUseCase
	customerUC CustomerUseCase
	merchantUC MerchantUseCase
}

func NewTransactionUseCase(repo repository.TransactionRepository, userUC UserUseCase, customerUC CustomerUseCase, merchantUC MerchantUseCase) TransactionUseCase {
	return &transactionUseCase{
		repo:       repo,
		userUC:     userUC,
		customerUC: customerUC,
		merchantUC: merchantUC,
	}
}

// GetTransaction implements TransactionUseCase.
func (usecase *transactionUseCase) GetTransactionByCustomerId(id string, params model.PaginationParams) ([]model.Transaction, error) {
	transactions, err := usecase.repo.GetByCustomerId(id, params)
	if err != nil {
		return []model.Transaction{}, err
	}
	return transactions, err
}

// ListTransaction implements TransactionUseCase.
func (usecase *transactionUseCase) ListTransaction(params model.PaginationParams) ([]model.Transaction, error) {
	transactions, err := usecase.repo.List(params)
	if err != nil {
		return []model.Transaction{}, err
	}

	return transactions, err
}

// RegisterNewTransaction implements TransactionUseCase.
func (usecase *transactionUseCase) RegisterNewTransaction(payload model.CreateTransactionRequest) (model.Transaction, error) {

	user, err := usecase.userUC.GetUserById(payload.UserId)
	if err != nil {
		return model.Transaction{}, fmt.Errorf("error getting user from user %v: %v", payload.UserId, err)
	}

	customer, err := usecase.customerUC.GetCustomerByUserId(user.ID)
	if err != nil {
		return model.Transaction{}, fmt.Errorf("error getting customer from customer with user id %v: %v", user.ID, err)
	}
	req := model.Transaction{
		ID:                 common.GenerateID(),
		SenderCustomerId:   customer.ID,
		ReceiverMerchantId: payload.ReceiverMerchantId,
		Amount:             payload.Amount,
	}

	transaction, err := usecase.repo.Create(req)
	if err != nil {
		return model.Transaction{}, err
	}

	return transaction, err
}
