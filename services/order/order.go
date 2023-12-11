// Package order services/order/order.go
package order

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"sync"
	"time"
)

type Order struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	MovieID    int       `json:"movie_id"`
	SeatNumber int       `json:"seat_number"` // 更新座位号字段类型
	OrderTime  time.Time `json:"order_time"`
}

type OrderService interface {
	GetOrderByID(ctx context.Context, args *PageRequest, reply *Order) error
	CreateOrder(ctx context.Context, args *Order, reply *struct{ OrderNumber string }) error
}

type orderServiceImpl struct {
	orders     map[int]*Order
	orderMutex sync.Mutex
	db         *sql.DB
}

type PageRequest struct {
	Args      *Order `json:"args,omitempty"`
	QueryType string `json:"serach_type,omitempty"`
}

func (s *orderServiceImpl) GetOrderByID(ctx context.Context, args *PageRequest, reply *Order) error {
	queryId := args.Args.ID
	switch args.QueryType {
	case "order_id":
		queryId = args.Args.ID
	case "user_id":
		queryId = args.Args.UserID
	case "movie_id":
		queryId = args.Args.MovieID
	}
	row := s.db.QueryRow("SELECT id, user_id, movie_id, seat_number, order_time FROM orders WHERE ? = ?", args.QueryType, queryId)
	err := row.Scan(&reply.ID, &reply.UserID, &reply.MovieID, &reply.SeatNumber, &reply.OrderTime)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("order not found")
		}
		return err
	}

	return nil
}

func (s *orderServiceImpl) CreateOrder(ctx context.Context, args *Order, reply *struct{ OrderNumber string }) error {
	args.OrderTime = time.Now() // 设置下单时间为当前时间
	result, err := s.db.Exec("INSERT INTO orders (user_id, movie_id, seat_number, order_time) VALUES (?, ?, ?, ?)", args.UserID, args.MovieID, args.SeatNumber, args.OrderTime)
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

	// 生成订单号
	hasher := md5.New()
	hasher.Write([]byte(args.OrderTime.String() + strconv.Itoa(args.UserID)))
	orderNumber := hex.EncodeToString(hasher.Sum(nil))
	reply.OrderNumber = orderNumber // 设置返回值中的订单号
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	args.ID = int(lastInsertID)

	return nil
}

func NewOrderService(db *sql.DB) OrderService {
	return &orderServiceImpl{
		db: db,
	}
}
