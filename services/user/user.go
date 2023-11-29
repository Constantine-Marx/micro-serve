package user

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID       int    `json:"ID"`
	Username string `json:"Username"`
	Email    string `json:"Email"`
}

type UserService interface {
	GetUserByID(ctx context.Context, args *User, reply *User) error
	CreateUser(ctx context.Context, args *User, reply *User) error
}

type userServiceImpl struct {
	db *sql.DB
}

func (s *userServiceImpl) GetUserByID(ctx context.Context, args *User, reply *User) error {
	row := s.db.QueryRow("SELECT id, username, email FROM users WHERE id = ?", args.ID)
	err := row.Scan(&reply.ID, &reply.Username, &reply.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user not found")
		}
		return err
	}
	return nil
}

func (s *userServiceImpl) CreateUser(ctx context.Context, args *User, reply *User) error {
	result, err := s.db.Exec("INSERT INTO users (id, username, email) VALUES (?, ?, ?)", args.ID, args.Username, args.Email)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user already exists")
	}
	return nil
}

func NewUserService(db *sql.DB) UserService {
	return &userServiceImpl{
		db: db,
	}
}
