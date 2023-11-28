package order

import (
	"fmt"
	"sync"
	"time"
)

type Order struct {
	ID        int
	UserID    int
	MovieID   int
	TicketNum int
	Date      time.Time
}

type OrderService interface {
	GetOrderByID(id int) (*Order, error)
	CreateOrder(order *Order) error
}

type orderServiceImpl struct {
	orders     map[int]*Order
	orderMutex sync.Mutex
}

func NewOrderService() OrderService {
	return &orderServiceImpl{
		orders: make(map[int]*Order),
	}
}

func (s *orderServiceImpl) GetOrderByID(id int) (*Order, error) {
	order, ok := s.orders[id]
	if !ok {
		return nil, fmt.Errorf("order not found")
	}
	return order, nil
}

func (s *orderServiceImpl) CreateOrder(order *Order) error {
	s.orderMutex.Lock()
	defer s.orderMutex.Unlock()

	if _, ok := s.orders[order.ID]; ok {
		return fmt.Errorf("order already exists")
	}
	s.orders[order.ID] = order
	return nil
}
