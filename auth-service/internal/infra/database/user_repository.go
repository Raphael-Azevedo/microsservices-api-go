package database

import (
	"biz-hub-auth-service/internal/entity"
	"database/sql"
	"errors"
)

type UserRepository struct {
	Db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{Db: db}
}

func (r *UserRepository) Create(user *entity.User) error {
	stmt, err := r.Db.Prepare("INSERT INTO users (id, name, email, password) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(user.ID, user.Name, user.Email, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) FindUserByEmail(email string) (*entity.User, error) {
	row := r.Db.QueryRow("SELECT id, name, email, password FROM users WHERE email = ?", email)

	user := &entity.User{}
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}

	return user, nil
}
