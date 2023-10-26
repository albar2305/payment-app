package repository

import (
	"database/sql"
	"strings"

	"github.com/albar2305/payment-app/model"
)

type UserRepository interface {
	Create(arg model.User) (model.User, error)
	Get(name string) (model.User, error)
	GetById(id string) (model.User, error)
	List(params model.PaginationParams) ([]model.UserResponse, error)
	Update(arg model.User) (model.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

// GetByUId implements UserRepository.
func (u *userRepository) GetById(id string) (model.User, error) {
	sql := `
	SELECT id, email, username, password, role, created_at from users where id = $1
	`
	row := u.db.QueryRow(sql, id)
	var i model.User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Username,
		&i.Password,
		&i.Role,
		&i.CreatedAt,
	)
	return i, err
}

// Update implements UserRepository.
func (u *userRepository) Update(arg model.User) (model.User, error) {
	sql := `UPDATE users
	SET
	  email = $1,
	  username = $2,
	  password = $3,
	  role = $4
	WHERE
	  id = $5
	RETURNING id,email, username, password, role, created_at`
	row := u.db.QueryRow(sql, arg.Email, arg.Username, arg.Password, arg.Role, arg.ID)
	var i model.User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Username,
		&i.Password,
		&i.Role,
		&i.CreatedAt,
	)
	return i, err
}

// Create implements UserRepository.
func (u *userRepository) Create(arg model.User) (model.User, error) {
	sql := `
	INSERT INTO users (
		id,
		email,
		username,
		password,
		role
	  ) VALUES (
		$1, $2, $3, $4, $5
	  ) RETURNING id,email, username, password, role, created_at`

	row := u.db.QueryRow(sql, arg.ID, arg.Email, arg.Username, arg.Password, strings.ToLower(arg.Role))
	var i model.User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Username,
		&i.Password,
		&i.Role,
		&i.CreatedAt,
	)
	return i, err
}

// Get implements UserRepository.
func (u *userRepository) Get(username string) (model.User, error) {
	sql := `
	SELECT id, email, username, password, role, created_at from users where username = $1
	`
	row := u.db.QueryRow(sql, username)
	var i model.User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Username,
		&i.Password,
		&i.Role,
		&i.CreatedAt,
	)
	return i, err
}

// List implements UserRepository.
func (u *userRepository) List(params model.PaginationParams) ([]model.UserResponse, error) {
	sql := `SELECT id, email, username, created_at from users
	ORDER BY created_at
	LIMIT $1
	OFFSET $2`
	rows, err := u.db.Query(sql, params.Limit, params.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []model.UserResponse{}
	for rows.Next() {
		var i model.UserResponse
		if err := rows.Scan(
			&i.ID,
			&i.Email,
			&i.Username,
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
