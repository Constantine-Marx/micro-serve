package user

import (
	"database/sql"
	"fmt"
	"log"
)

type User struct {
	ID       int
	Username string
	Email    string
}

type Service interface {
	GetUserByID(id int) (*User, error)
	CreateUser(user *User) error
}

type userServiceImpl struct {
	users map[int]*User
	db    *sql.DB
}

func NewUserService(db *sql.DB) Service {
	return &userServiceImpl{
		db:    db,
		users: make(map[int]*User),
	}
}

func (s *userServiceImpl) GetUserByID(id int) (*User, error) {
	user, ok := s.users[id]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}

	// 示例：从数据库中读取用户信息
	row := s.db.QueryRow("SELECT email FROM users WHERE id = ?", id)
	var email string
	err := row.Scan(&email)
	if err != nil {
		log.Printf("Failed to get user email from database: %v", err)
	} else {
		user.Email = email
	}

	return user, nil
}

func (s *userServiceImpl) CreateUser(user *User) error {
	if _, ok := s.users[user.ID]; ok {
		return fmt.Errorf("user already exists")
	}
	s.users[user.ID] = user

	// 示例：将用户信息写入数据库
	_, err := s.db.Exec("INSERT INTO users (id, username, email) VALUES (?, ?, ?)", user.ID, user.Username, user.Email)
	if err != nil {
		log.Printf("Failed to insert user into database: %v", err)
	}

	return nil
}
