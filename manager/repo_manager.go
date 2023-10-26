package manager

import "github.com/albar2305/payment-app/repository"

type RepoManager interface {
	UserRepo() repository.UserRepository
	CustomerRepo() repository.CustomerRepository
	MerchantRepo() repository.MerchantRepository
	TransactionRepo() repository.TransactionRepository
}

type repoManager struct {
	infra InfraManager
}

// TransactionRepo implements RepoManager.
func (r *repoManager) TransactionRepo() repository.TransactionRepository {
	return repository.NewTransactionRepository(r.infra.Conn())
}

// MerchantRepo implements RepoManager.
func (r *repoManager) MerchantRepo() repository.MerchantRepository {
	return repository.NewMerchanRepository(r.infra.Conn())
}

// CustomerRepo implements RepoManager.
func (r *repoManager) CustomerRepo() repository.CustomerRepository {
	return repository.NewCustomerRepository(r.infra.Conn())
}

func (r *repoManager) UserRepo() repository.UserRepository {
	return repository.NewUserRepository(r.infra.Conn())
}

func NewRepoManager(infra InfraManager) RepoManager {
	return &repoManager{infra: infra}
}
