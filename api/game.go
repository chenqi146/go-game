package main

import (
	_ "embed"
	"flag"
	"fmt"
	"go-game/common/middleware"
	"net/http"
	
	"go-game/api/internal/config"
	"go-game/api/internal/handler"
	"go-game/api/internal/svc"
	
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/game.yaml", "the config file")

//go:embed etc/words-cn.txt
var WordPairs string

// todo 房间游戏心跳
func main() {
	flag.Parse()
	
	var c config.Config
	conf.MustLoad(*configFile, &c)
	
	server := rest.MustNewServer(c.RestConf, rest.WithCustomCors(cors, nil, "*"))
	
	defer server.Stop()
	
	server.Use(middleware.UserIdMiddleware)
	
	ctx := svc.NewServiceContext(c, WordPairs)
	handler.RegisterHandlers(server, ctx)
	
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

func cors(header http.Header) {
	header.Set("Access-Control-Allow-Headers", "*")
	header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
	header.Set("Access-Control-Expose-Headers", "Content-Length, Content-Type, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
	header.Set("Access-Control-Allow-Credentials", "true")
}
