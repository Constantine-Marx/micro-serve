// Project Name: go_learning_rpcx
package main

import (
	"flag"
	etcdclient "github.com/rpcxio/rpcx-etcd/client"
	"github.com/smallnest/rpcx/client"
)

var (
	etcdAddr = flag.String("etcdAddr", "localhost:2379", "etcd address")
	basePath = flag.String("base", "rpcx_test", "prefix path")
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
}
