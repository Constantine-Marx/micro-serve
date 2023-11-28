package main

import (
	"context"
	"flag"
	"log"
	"time"

	etcdclient "github.com/rpcxio/rpcx-etcd/client"
	"github.com/smallnest/rpcx/client"

	"rpcx/services/movie"
	serviceorder "rpcx/services/order"
	"rpcx/services/user"
)

var (
	etcdAddr = flag.String("etcdAddr", "localhost:2379", "etcd address")
	basePath = flag.String("base", "/rpcx_test", "prefix path")
)

func main() {
	flag.Parse()

	d, _ := etcdclient.NewEtcdV3Discovery(*basePath, "UserService", []string{*etcdAddr}, false, nil)
	userClient := client.NewXClient("UserService", client.Failover, client.RoundRobin, d, client.DefaultOption)
	defer func(userClient client.XClient) {
		_ = userClient.Close()
	}(userClient)

	d, _ = etcdclient.NewEtcdV3Discovery(*basePath, "MovieService", []string{*etcdAddr}, false, nil)
	movieClient := client.NewXClient("MovieService", client.Failover, client.RoundRobin, d, client.DefaultOption)
	defer func(movieClient client.XClient) {
		_ = movieClient.Close()
	}(movieClient)

	d, _ = etcdclient.NewEtcdV3Discovery(*basePath, "OrderService", []string{*etcdAddr}, false, nil)
	orderClient := client.NewXClient("OrderService", client.Failover, client.RoundRobin, d, client.DefaultOption)
	defer func(orderClient client.XClient) {
		_ = orderClient.Close()
	}(orderClient)

	// 示例：创建用户
	users := &user.User{
		ID:       1,
		Username: "John",
		Email:    "john@example.com",
	}
	err := userClient.Call(context.Background(), "CreateUser", users, nil)
	if err != nil {
		log.Fatalf("failed to create user: %v", err)
	}

	// 示例：创建电影
	movies := &movie.Movie{
		ID:     1,
		Title:  "Inception",
		Rating: 8.8,
	}
	err = movieClient.Call(context.Background(), "CreateMovie", movies, nil)
	if err != nil {
		log.Fatalf("failed to create movie: %v", err)
	}

	// 示例：创建订单
	orders := &serviceorder.Order{
		ID:        1,
		UserID:    1,
		MovieID:   1,
		TicketNum: 2,
		Date:      time.Now(),
	}
	err = orderClient.Call(context.Background(), "CreateOrder", orders, nil)
	if err != nil {
		log.Fatalf("failed to create order: %v", err)
	}

	// 示例：获取订单
	var getOrder serviceorder.Order
	err = orderClient.Call(context.Background(), "GetOrderByID", 1, &getOrder)
	if err != nil {
		log.Fatalf("failed to get order: %v", err)
	}
	log.Printf("Order: %+v", getOrder)
}
