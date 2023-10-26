package repository

import (
	"database/sql"

	"github.com/albar2305/payment-app/model"
)

type MerchantRepository interface {
	Create(arg model.Merchant) (model.Merchant, error)
	Delete(id string) error
	Get(id string) (model.Merchant, error)
	List(params model.PaginationParams) ([]model.Merchant, error)
}

type merchantRepository struct {
	db *sql.DB
}

func NewMerchanRepository(db *sql.DB) MerchantRepository {
	return &merchantRepository{db: db}
}

// Create implements MerchantRepository.
func (repo *merchantRepository) Create(arg model.Merchant) (model.Merchant, error) {
	sql := `
	INSERT INTO merchants (
		id, name, description, business_type, balance
	  ) VALUES (
		$1, $2, $3, $4, $5
	  ) RETURNING id, name, description, business_type, balance, created_at`

	row := repo.db.QueryRow(sql, arg.ID, arg.Name, arg.Description, arg.BusinesType, arg.Balance)
	var i model.Merchant
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.BusinesType,
		&i.Balance,
		&i.CreatedAt,
	)

	return i, err
}

// Delete implements MerchantRepository.
func (repo *merchantRepository) Delete(id string) error {
	sql := `
	DELETE FROM merchants
	WHERE id = $1
	`

	_, err := repo.db.Exec(sql, id)
	return err
}

// Get implements MerchantRepository.
func (repo *merchantRepository) Get(id string) (model.Merchant, error) {
	sql := `SELECT id, name, description, business_type, balance, created_at FROM merchants
	WHERE name = $1 LIMIT 1`
	row := repo.db.QueryRow(sql, id)
	var i model.Merchant
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.BusinesType,
		&i.Balance,
		&i.CreatedAt,
	)
	return i, err
}

// List implements MerchantRepository.
func (repo *merchantRepository) List(params model.PaginationParams) ([]model.Merchant, error) {
	sql := `SELECT id, name, description, business_type, balance, created_at FROM merchants
	ORDER BY created_at
	LIMIT $1
	OFFSET $2`
	rows, err := repo.db.Query(sql, params.Limit, params.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []model.Merchant{}
	for rows.Next() {
		var i model.Merchant
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.BusinesType,
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
