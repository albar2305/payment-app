package usecase

import (
	"github.com/albar2305/payment-app/model"
	"github.com/albar2305/payment-app/repository"
	"github.com/albar2305/payment-app/utils/common"
)

type UserUseCase interface {
	RegisterNewUser(payload model.CreateUserRequest) (model.User, error)
	GetUser(id string) (model.User, error)
	GetUserById(id string) (model.User, error)
	ListUser(params model.PaginationParams) ([]model.UserResponse, error)
	UpdateUser(user model.User) (model.User, error)
}

type userUseCase struct {
	repo repository.UserRepository
}

// GetUser implements UserUseCase.
func (usecase *userUseCase) GetUser(username string) (model.User, error) {
	return usecase.repo.Get(username)
}

// GetUser implements UserUseCase.
func (usecase *userUseCase) GetUserById(id string) (model.User, error) {
	return usecase.repo.GetById(id)
}

// ListUser implements UserUseCase.
func (usecase *userUseCase) ListUser(params model.PaginationParams) ([]model.UserResponse, error) {
	return usecase.repo.List(params)
}

// UpdateUser implements UserUseCase.
func (usecase *userUseCase) UpdateUser(user model.User) (model.User, error) {
	return usecase.repo.Update(user)
}

// RegisterNewUser implements UserUseCase.
func (usecase *userUseCase) RegisterNewUser(payload model.CreateUserRequest) (model.User, error) {
	hashedPassword, err := common.HashPassword(payload.Password)
	if err != nil {
		return model.User{}, err
	}
	if payload.Role == "" {
		payload.Role = "user"
	}

	userRequest := model.User{
		ID:       common.GenerateID(),
		Email:    payload.Email,
		Username: payload.Username,
		Password: hashedPassword,
		Role:     payload.Role,
	}

	user, err := usecase.repo.Create(userRequest)
	return user, err
}

func NewUserUseCase(repo repository.UserRepository) UserUseCase {
	return &userUseCase{repo: repo}
}
