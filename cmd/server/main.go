package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os/exec"
	"rpcx/services/storage"
	"time"

	"github.com/rpcxio/rpcx-etcd/serverplugin"
	"github.com/smallnest/rpcx/server"

	"rpcx/services/movie"
	"rpcx/services/order"
	"rpcx/services/user"
)

var (
	addr     = flag.String("addr", "localhost:8972", "server address")
	etcdAddr = flag.String("etcdAddr", "localhost:2379", "etcd address")
	basePath = flag.String("base", "/rpcx_test", "prefix path")
)

func main() {
	flag.Parse()

	//启动etcd
	cmd := exec.Command("etcd")
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Second * 2)

	go func() {
		_ = http.ListenAndServe(":9981", nil)
	}()

	// 连接到数据库
	db, err := storage.ConnectDB("root", "user", "localhost:3306", "movie_ticket_service")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	s := server.NewServer()
	addRegistryPlugin(s)

	_ = s.RegisterName("UserService", user.NewUserService(db), "")
	_ = s.RegisterName("MovieService", movie.NewMovieService(), "")
	_ = s.RegisterName("OrderService", order.NewOrderService(), "")

	err = s.Serve("tcp", *addr)
	if err != nil {
		panic(err)
	}
}

func addRegistryPlugin(s *server.Server) {
	r := &serverplugin.EtcdV3RegisterPlugin{
		ServiceAddress: "tcp@" + *addr,
		EtcdServers:    []string{*etcdAddr},
		BasePath:       *basePath,
		UpdateInterval: time.Minute,
	}
	err := r.Start()
	if err != nil {
		log.Fatal(err)
	}
	s.Plugins.Add(r)
}
