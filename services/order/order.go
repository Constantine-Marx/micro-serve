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
	"strings"
	"sync"
	"time"
)

type Order struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	ScheduleID int       `json:"schedule_id"`
	MovieID    int       `json:"movie_id"`
	SeatRow    int       `json:"seat_row"`
	SeatColumn int       `json:"seat_column"`
	OrderTime  time.Time `json:"order_time"`
}

type OrderService interface {
	GetOrderByID(ctx context.Context, args *PageRequest, reply *Order) error
	CreateOrder(ctx context.Context, args *Order, reply *struct {
		OrderNumber string
		status      int
	}) error
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
	row := s.db.QueryRow("SELECT id, user_id, movie_id, schedule_id, seat_row, seat_column, order_time FROM orders WHERE id = ?", queryId)
	err := row.Scan(&reply.ID, &reply.UserID, &reply.MovieID, &reply.SeatRow, &reply.SeatColumn, &reply.OrderTime)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("order not found")
		}
		return err
	}

	return nil
}

func (s *orderServiceImpl) CreateOrder(ctx context.Context, args *Order, reply *struct {
	OrderNumber string
	status      int
}) error {

	args.OrderTime = time.Now() // 设置下单时间为当前时间
	result, err := s.db.Exec("INSERT INTO orders (user_id, movie_id, schedule_id, seat_row, seat_column, order_time) VALUES (?, ?, ?, ?, ?, ?)", args.UserID, args.MovieID, args.ScheduleID, args.SeatRow, args.SeatColumn, args.OrderTime)
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

	// 更新座位状态
	//scheduleID := 1 // 从请求参数或数据库获取相应的排期ID
	//seatRows := []int{args.SeatRow}
	//seatColumns := []int{args.SeatColumn}
	//err = s.updateSeats(scheduleID, seatRows, seatColumns)
	if err != nil {
		return err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	args.ID = int(lastInsertID)

	return nil
}

func (s *orderServiceImpl) updateSeats(scheduleID int, seatRows []int, seatColumns []int) error {
	// 查询 movie_schedules 表以获取当前座位状态
	var seats string
	err := s.db.QueryRow("SELECT seats FROM movie_schedules WHERE id = ?", scheduleID).Scan(&seats)
	if err != nil {
		return err
	}

	// 将座位字符串转换为二维数组
	seatsArray := make([][]int, 0)
	for _, row := range strings.Split(seats, ";") {
		seatRow := make([]int, 0)
		for _, seat := range strings.Split(row, ",") {
			seatInt, _ := strconv.Atoi(seat)
			seatRow = append(seatRow, seatInt)
		}
		seatsArray = append(seatsArray, seatRow)
	}

	// 更新座位状态
	for i, seatRow := range seatRows {
		seatColumn := seatColumns[i]
		seatsArray[seatRow][seatColumn] = 1
	}

	// 将二维数组转换回字符串
	updatedSeats := ""
	for i, row := range seatsArray {
		if i > 0 {
			updatedSeats += ";"
		}
		for j, seat := range row {
			if j > 0 {
				updatedSeats += ","
			}
			updatedSeats += strconv.Itoa(seat)
		}
	}

	// 更新 movie_schedules 表中的座位信息
	_, err = s.db.Exec("UPDATE movie_schedules SET seats = ? WHERE id = ?", updatedSeats, scheduleID)
	if err != nil {
		return err
	}

	return nil
}

func NewOrderService(db *sql.DB) OrderService {
	return &orderServiceImpl{
		db: db,
	}
}
