package repository

import (
	"database/sql"
	"defect-management-system/backend/internal/model"
	"errors"
)

type UserRepository struct {
	db *DB
}

func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByLogin(login string) (*model.User, error) {
	query := "SELECT id, login, password_hash, role, created_at FROM users WHERE login = ?"
	if r.db.Driver != "sqlite" {
		query = "SELECT id, login, password_hash, role, created_at FROM users WHERE login = $1"
	}

	var user model.User
	err := r.db.QueryRow(query, login).Scan(&user.ID, &user.Login, &user.PasswordHash, &user.Role, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByID(id int) (*model.User, error) {
	query := "SELECT id, login, password_hash, role, created_at FROM users WHERE id = ?"
	if r.db.Driver != "sqlite" {
		query = "SELECT id, login, password_hash, role, created_at FROM users WHERE id = $1"
	}

	var user model.User
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Login, &user.PasswordHash, &user.Role, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
