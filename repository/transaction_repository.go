package repository

import (
	"database/sql"

	"github.com/albar2305/payment-app/model"
)

type TransactionRepository interface {
	Create(arg model.Transaction) (model.Transaction, error)
	GetByCustomerId(id string, params model.PaginationParams) ([]model.Transaction, error)
	List(params model.PaginationParams) ([]model.Transaction, error)
}

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

// Create implements TransactionRepository.
func (repo *transactionRepository) Create(arg model.Transaction) (model.Transaction, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return model.Transaction{}, err
	}
	sql := `
	INSERT INTO transactions (
		id,
		sender_customer_id,
		receiver_merchant_id,
		amount
	  ) VALUES (
		$1, $2, $3, $4
	  ) RETURNING id, sender_customer_id, receiver_merchant_id, amount, created_at`
	row := tx.QueryRow(sql, arg.ID, arg.SenderCustomerId, arg.ReceiverMerchantId, arg.Amount)
	var i model.Transaction
	err = row.Scan(
		&i.ID,
		&i.SenderCustomerId,
		&i.ReceiverMerchantId,
		&i.Amount,
		&i.CreatedAt,
	)
	sql = `UPDATE merchants
	SET balance = balance + $1
	WHERE id = $2`
	_, _ = tx.Exec(sql, arg.Amount, arg.ReceiverMerchantId)

	sql = `UPDATE customers
	SET balance = balance + $1
	WHERE user_id = $2`

	_, _ = tx.Exec(sql, -arg.Amount, arg.SenderCustomerId)

	if err := tx.Commit(); err != nil {
		return model.Transaction{}, err
	}
	return i, err

}

// Get implements TransactionRepository.
func (repo *transactionRepository) GetByCustomerId(id string, params model.PaginationParams) ([]model.Transaction, error) {
	sql := `SELECT id, sender_customer_id,receiver_merchant_id,amount,created_at from transactions WHERE sender_customer_id = $1
	ORDER BY created_at
	LIMIT $2
	OFFSET $3`
	rows, err := repo.db.Query(sql, id, params.Limit, params.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []model.Transaction{}
	for rows.Next() {
		var i model.Transaction
		if err := rows.Scan(
			&i.ID,
			&i.SenderCustomerId,
			&i.ReceiverMerchantId,
			&i.Amount,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

// List implements TransactionRepository.
func (repo *transactionRepository) List(params model.PaginationParams) ([]model.Transaction, error) {
	sql := `SELECT id, sender_customer_id,receiver_merchant_id,amount,created_at from transactions
	ORDER BY created_at
	LIMIT $1
	OFFSET $2`
	rows, err := repo.db.Query(sql, params.Limit, params.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []model.Transaction{}
	for rows.Next() {
		var i model.Transaction
		if err := rows.Scan(
			&i.ID,
			&i.SenderCustomerId,
			&i.ReceiverMerchantId,
			&i.Amount,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
