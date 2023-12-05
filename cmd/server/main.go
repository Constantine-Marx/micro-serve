// Description: 服务端启动入口
// ./cmd/server/main.go
package main

import (
	"database/sql"
	"flag"
	"log"
	_ "net/http/pprof"
	"os/exec"
	"time"

	"github.com/rpcxio/rpcx-etcd/serverplugin"
	"github.com/smallnest/rpcx/server"

	"rpcx/services/movie"
	"rpcx/services/movie_schedule"
	"rpcx/services/order"
	"rpcx/services/user"
	"rpcx/utils/storage"
	"rpcx/utils/utils"
)

var (
	etcdAddr = flag.String("etcdAddr", "localhost:2379", "etcd address")
	basePath = flag.String("base", "rpcx_test/", "prefix path")
)

func startUserService(addr string, db *sql.DB) {
	s := server.NewServer()
	addRegistryPlugin(s, addr)

	err := s.RegisterName("UserService", user.NewUserService(db), "")
	if err != nil {
		panic(err)
	}

	err = s.Serve("tcp", addr)
	if err != nil {
		panic(err)
	}
}

func startMovieService(addr string, db *sql.DB) {
	s := server.NewServer()
	addRegistryPlugin(s, addr)

	err := s.RegisterName("MovieService", movie.NewMovieService(db), "")
	if err != nil {
		panic(err)
	}

	err = s.Serve("tcp", addr)
	if err != nil {
		panic(err)
	}
}

func startOrderService(addr string, db *sql.DB) {
	s := server.NewServer()
	addRegistryPlugin(s, addr)

	err := s.RegisterName("OrderService", order.NewOrderService(db), "")
	if err != nil {
		panic(err)
	}

	err = s.Serve("tcp", addr)
	if err != nil {
		panic(err)
	}
}

func startUtilService(addr string, db *sql.DB) {
	s := server.NewServer()
	addRegistryPlugin(s, addr)

	err := s.RegisterName("UtilService", utils.NewExtractService(db), "")
	if err != nil {
		panic(err)
	}

	err = s.Serve("tcp", addr)
	if err != nil {
		panic(err)
	}
}

// main.go
func startMovieScheduleService(addr string, db *sql.DB) {
	s := server.NewServer()
	addRegistryPlugin(s, addr)

	// 创建 MovieScheduleService 实例并注册
	movieScheduleService := movie_schedule.NewMovieScheduleService(db)
	err := s.RegisterName("MovieScheduleService", movieScheduleService, "")
	if err != nil {
		panic(err)
	}

	err = s.Serve("tcp", addr)
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()

	// 启动 etcd
	cmd := exec.Command("etcd")
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	// 连接到数据库
	db, err := storage.ConnectDB("root", "228809", "localhost:3306", "movie_ticket_service")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	// 并行启动服务
	go startUserService(":8972", db)
	go startMovieService(":8973", db)
	go startOrderService(":8974", db)
	go startUtilService(":8975", db)
	go startMovieScheduleService(":8976", db)

	// 阻塞 main goroutine，以便服务继续运行
	select {}
}

func addRegistryPlugin(s *server.Server, addr string) {
	r := &serverplugin.EtcdV3RegisterPlugin{
		ServiceAddress: "tcp@" + addr, // 使用传入的 addr 而不是全局 *addr
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
