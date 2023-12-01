// Path: cmd/gateway/main.go
package main

import (
	"flag"
	"github.com/rs/cors"
	"net/http"

	etcdclient "github.com/rpcxio/rpcx-etcd/client"
	gateway "github.com/rpcxio/rpcx-gateway"
	"github.com/smallnest/rpcx/client"
	"log"
)

var (
	addr     = flag.String("addr", "localhost:9981", "gateway address")
	etcdAddr = flag.String("etcdAddr", "localhost:2379", "etcd address")
	basePath = flag.String("base", "rpcx_test/", "prefix path")
)

type MyHTTPServer struct {
	server *http.Server
}

func (m *MyHTTPServer) RegisterHandler(base string, handler gateway.ServiceHandler) {
	http.HandleFunc(base, func(w http.ResponseWriter, r *http.Request) {
		meta, payload, err := handler(r, r.URL.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for k, v := range meta {
			w.Header().Set(k, v)
		}
		_, _ = w.Write(payload)
	})
}

func (m *MyHTTPServer) Serve() error {
	return m.server.ListenAndServe()
}

func main() {
	flag.Parse()

	// 创建 etcd 服务发现
	sd, err := etcdclient.NewEtcdV3DiscoveryTemplate(*basePath, []string{*etcdAddr}, false, nil)
	if err != nil {
		log.Fatalf("Failed to create etcd service discovery: %v", err)
	}

	// 创建一个 HTTP 服务器实例
	myHTTPServer := &MyHTTPServer{
		server: &http.Server{
			Addr: *addr,
		},
	}

	// 添加 CORS 处理
	corsWrapper := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // 允许所有域名
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		Debug:            false,
	})

	myHTTPServer.server.Handler = corsWrapper.Handler(http.DefaultServeMux)

	// 创建 RPCX 网关
	g := gateway.NewGateway("/", myHTTPServer, sd, client.Failover, client.RoundRobin, client.DefaultOption)

	log.Printf("RPCX gateway is running on %s", *addr)

	// 启动 HTTP 服务器
	if err := g.Serve(); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
