package repository

import (
	"database/sql"

	"github.com/albar2305/payment-app/model"
)

type CustomerRepository interface {
	Create(arg model.Customer) (model.Customer, error)
	Delete(id string) error
	GetByUserId(userId string) (model.Customer, error)
	GetById(id string) (model.Customer, error)
	List(params model.PaginationParams) ([]model.Customer, error)
	AddCustomerBalance(id string, amount int64) (model.Customer, error)
}

type customerRepository struct {
	db *sql.DB
}

// AddCustomerBalance implements CustomerRepository.
func (c *customerRepository) AddCustomerBalance(id string, amount int64) (model.Customer, error) {
	sql := `UPDATE customers
	SET balance = balance + $1
	WHERE user_id = $2
	RETURNING id, user_id,name, balance, created_at`
	row := c.db.QueryRow(sql, amount, id)
	var i model.Customer
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Name,
		&i.Balance,
		&i.CreatedAt,
	)

	return i, err
}

// Create implements CustomerRepository.
func (c *customerRepository) Create(arg model.Customer) (model.Customer, error) {
	sql := `
	INSERT INTO customers (
		id,
		user_id,
		name,
		balance
	  ) VALUES (
		$1, $2, $3, $4
	  ) RETURNING id,user_id, name, balance, created_at`

	row := c.db.QueryRow(sql, arg.ID, arg.UserID, arg.Name, arg.Balance)
	var i model.Customer
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Name,
		&i.Balance,
		&i.CreatedAt,
	)

	return i, err

}

// Delete implements CustomerRepository.
func (c *customerRepository) Delete(id string) error {
	sql := `
	DELETE FROM customers
	WHERE id = $1
	`

	_, err := c.db.Exec(sql, id)
	return err
}

// Get implements CustomerRepository.
func (c *customerRepository) GetByUserId(userId string) (model.Customer, error) {
	sql := `SELECT id, user_id, name, balance, created_at FROM customers
	WHERE user_id = $1 LIMIT 1`
	row := c.db.QueryRow(sql, userId)
	var i model.Customer
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Name,
		&i.Balance,
		&i.CreatedAt,
	)
	return i, err
}

func (c *customerRepository) GetById(id string) (model.Customer, error) {
	sql := `SELECT id, user_id, name, balance, created_at FROM customers
	WHERE id = $1 LIMIT 1`
	row := c.db.QueryRow(sql, id)
	var i model.Customer
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Name,
		&i.Balance,
		&i.CreatedAt,
	)
	return i, err
}

// List implements CustomerRepository.
func (c *customerRepository) List(params model.PaginationParams) ([]model.Customer, error) {
	sql := `SELECT id,user_id, name, balance, created_at FROM customers
	ORDER BY created_at
	LIMIT $1
	OFFSET $2`
	rows, err := c.db.Query(sql, params.Limit, params.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []model.Customer{}
	for rows.Next() {
		var i model.Customer
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Name,
			&i.Balance,
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

func NewCustomerRepository(db *sql.DB) CustomerRepository {
	return &customerRepository{db: db}
}
