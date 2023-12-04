// Package order services/order/order.go
package order

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Order struct {
	ID        int       `json:"ID"`
	UserID    int       `json:"UserID"`
	MovieID   int       `json:"MovieID"`
	TicketNum int       `json:"TicketNum"`
	Date      time.Time `json:"Date"`
}

type OrderService interface {
	GetOrderByID(ctx context.Context, args *Order, reply *Order) error
	CreateOrder(ctx context.Context, args *Order, reply *struct{}) error
}

type orderServiceImpl struct {
	orders     map[int]*Order
	orderMutex sync.Mutex
	db         *sql.DB
}

func (s *orderServiceImpl) GetOrderByID(ctx context.Context, args *Order, reply *Order) error {
	var date []byte
	row := s.db.QueryRow("SELECT id, UserID, MovieID, TicketNum, date FROM orders WHERE id = ?", args.ID)
	err := row.Scan(&reply.ID, &reply.UserID, &reply.MovieID, &reply.TicketNum, &date)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("order not found")
		}
		return err
	}

	parsedDate, err := time.Parse("2006-01-02 15:04:05", string(date))
	if err != nil {
		return fmt.Errorf("failed to parse date: %v", err)
	}
	reply.Date = parsedDate

	return nil
}

func (s *orderServiceImpl) CreateOrder(ctx context.Context, order *Order, reply *struct{}) error {
	result, err := s.db.Exec("INSERT INTO orders (id, UserID, MovieID, TicketNum, date) VALUES (?, ?, ?, ?, ?)", order.ID, order.UserID, order.MovieID, order.TicketNum, order.Date)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("order already exists")
	}
	return nil
}

func NewOrderService(db *sql.DB) OrderService {
	return &orderServiceImpl{
		db: db,
	}
}
