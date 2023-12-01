// Package user services/user/user.go
package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID       int    `json:"ID"`
	Username string `json:"Username"`
	Email    string `json:"Email"`
	Password string `json:"Password,omitempty"`
}

type UserService interface {
	GetUserByID(ctx context.Context, args *User, reply *User) error
	CreateUser(ctx context.Context, args *User, reply *User) error
	Login(ctx context.Context, args *User, reply *User) error
}

type userServiceImpl struct {
	db *sql.DB
}

func (s *userServiceImpl) GetUserByID(ctx context.Context, args *User, reply *User) error {
	log.Printf("GetUserByID called with args: %+v", args)
	row := s.db.QueryRow("SELECT id, username, email FROM users WHERE id = ?", args.ID)
	err := row.Scan(&reply.ID, &reply.Username, &reply.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("User not found: %v", err)
			return fmt.Errorf("user not found")
		}
		log.Printf("GetUserByID error: %v", err)
		return err
	}
	log.Printf("GetUserByID success, reply: %+v", reply)
	return nil
}
func (s *userServiceImpl) CreateUser(ctx context.Context, args *User, reply *User) error {
	result, err := s.db.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", args.Username, args.Email, args.Password)
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

	// Get the auto-generated ID
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	args.ID = int(lastInsertID)

	return nil
}

func (s *userServiceImpl) Login(ctx context.Context, args *User, reply *User) error {
	log.Printf("Login called with args: %+v", args)
	row := s.db.QueryRow("SELECT id, username, email FROM users WHERE username = ? AND password = ?", args.Username, args.Password)
	err := row.Scan(&reply.ID, &reply.Username, &reply.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("Invalid credentials: %v", err)
			return fmt.Errorf("invalid credentials")
		}
		log.Printf("Login error: %v", err)
		return err
	}
	log.Printf("Login success, reply: %+v", reply)
	return nil
}

func NewUserService(db *sql.DB) UserService {
	return &userServiceImpl{
		db: db,
	}
}
