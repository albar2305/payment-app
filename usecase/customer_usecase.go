package usecase

import (
	"fmt"

	"github.com/albar2305/payment-app/model"
	"github.com/albar2305/payment-app/repository"
	"github.com/albar2305/payment-app/utils/common"
)

type CustomerUseCase interface {
	RegisterNewCustomer(payload model.CreateCustomerRequest) (model.CustomerResponse, error)
	GetCustomerById(id string) (model.CustomerResponse, error)
	GetCustomerByUserId(userId string) (model.CustomerResponse, error)
	ListCustomer(params model.PaginationParams) ([]model.CustomerResponse, error)
	AddCustomerBalance(id string, amount int64) (model.CustomerResponse, error)
	DeleteCustomer(id string) error
}

type customerUseCase struct {
	repo        repository.CustomerRepository
	userUseCase UserUseCase
}

func NewCustomerUseCase(repo repository.CustomerRepository, userUseCase UserUseCase) CustomerUseCase {
	return &customerUseCase{
		repo: repo,
		// db:          db,
		userUseCase: userUseCase,
	}
}

// DeleteCustomer implements CustomerUseCase.
func (usecase *customerUseCase) DeleteCustomer(id string) error {
	customer, err := usecase.GetCustomerById(id)
	if err != nil {
		return fmt.Errorf("customer with ID %s not found", id)
	}

	err = usecase.repo.Delete(customer.ID)
	if err != nil {
		return fmt.Errorf("failed to delete customer: %v", err.Error())
	}
	return nil

}

// AddCustomerBalance implements CustomerUseCase.
func (usecase *customerUseCase) AddCustomerBalance(id string, amount int64) (model.CustomerResponse, error) {

	customer, err := usecase.repo.AddCustomerBalance(id, amount)
	if err != nil {
		return model.CustomerResponse{}, err
	}

	user, err := usecase.userUseCase.GetUserById(customer.UserID)
	if err != nil {
		return model.CustomerResponse{}, err
	}

	customerResponse := model.CustomerResponse{
		ID: customer.ID,
		User: model.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Username:  user.Username,
			CreatedAt: user.CreatedAt,
		},
		Name:      customer.Name,
		Balance:   customer.Balance,
		CreatedAt: customer.CreatedAt,
	}
	return customerResponse, err
}

// GetCustomer implements CustomerUseCase.
func (usecase *customerUseCase) GetCustomerByUserId(userId string) (model.CustomerResponse, error) {
	customer, err := usecase.repo.GetByUserId(userId)
	if err != nil {
		return model.CustomerResponse{}, fmt.Errorf("error getting customer from repository %v: %v", userId, err)
	}

	user, err := usecase.userUseCase.GetUserById(customer.UserID)
	if err != nil {
		return model.CustomerResponse{}, fmt.Errorf("error getting user %v: %v", customer.UserID, err)
	}

	customerResponse := model.CustomerResponse{
		ID: customer.ID,
		User: model.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Username:  user.Username,
			CreatedAt: user.CreatedAt,
		},
		Name:      customer.Name,
		Balance:   customer.Balance,
		CreatedAt: customer.CreatedAt,
	}
	return customerResponse, err
}

func (usecase *customerUseCase) GetCustomerById(id string) (model.CustomerResponse, error) {
	customer, err := usecase.repo.GetById(id)
	if err != nil {
		return model.CustomerResponse{}, err
	}

	user, err := usecase.userUseCase.GetUserById(customer.ID)
	if err != nil {
		return model.CustomerResponse{}, err
	}

	customerResponse := model.CustomerResponse{
		ID: customer.ID,
		User: model.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Username:  user.Username,
			CreatedAt: user.CreatedAt,
		},
		Name:      customer.Name,
		Balance:   customer.Balance,
		CreatedAt: customer.CreatedAt,
	}
	return customerResponse, err
}

// ListCustomer implements CustomerUseCase.
func (usecase *customerUseCase) ListCustomer(params model.PaginationParams) ([]model.CustomerResponse, error) {
	customers, err := usecase.repo.List(params)
	if err != nil {
		return []model.CustomerResponse{}, err
	}

	customerResponses := []model.CustomerResponse{}
	for _, customer := range customers {
		user, err := usecase.userUseCase.GetUserById(customer.UserID)
		if err != nil {
			return []model.CustomerResponse{}, err
		}
		customerResponse := model.CustomerResponse{
			ID: customer.ID,
			User: model.UserResponse{
				ID:        user.ID,
				Email:     user.Email,
				Username:  user.Username,
				CreatedAt: user.CreatedAt,
			},
			Name:      customer.Name,
			Balance:   customer.Balance,
			CreatedAt: customer.CreatedAt,
		}
		customerResponses = append(customerResponses, customerResponse)
	}

	return customerResponses, err
}

// RegisterNewCustomer implements CustomerUseCase.
func (usecase *customerUseCase) RegisterNewCustomer(payload model.CreateCustomerRequest) (model.CustomerResponse, error) {
	customer := model.Customer{
		ID:      common.GenerateID(),
		UserID:  payload.UserID,
		Name:    payload.Name,
		Balance: 0,
	}

	user, err := usecase.userUseCase.GetUserById(customer.UserID)
	if err != nil {
		return model.CustomerResponse{}, err
	}

	result, err := usecase.repo.Create(customer)
	if err != nil {
		return model.CustomerResponse{}, err
	}

	customerResponse := model.CustomerResponse{
		ID: result.ID,
		User: model.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Username:  user.Username,
			CreatedAt: user.CreatedAt,
		},
		Name:      result.Name,
		Balance:   result.Balance,
		CreatedAt: result.CreatedAt,
	}
	return customerResponse, err
}
