package main

import (
	"flag"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/acl-dev/master-gin"
)

var (
	filePath = flag.String("c", "gin-server.cf", "configure file")
	listenAddrs = flag.String("listen", "127.0.0.1:8088", "listen addresses")
)

func main()  {
	flag.Parse()

	var engines []*gin.Engine
	var err error

	master_gin.Init()

	if master_gin.Alone {
		fmt.Println("listen:", listenAddrs)
		engines, err = master_gin.AloneStart(*listenAddrs)

	} else {
		engines, err = master_gin.DaemonStart()
	}

	if err != nil {
		fmt.Println("start server failed on", listenAddrs)
		return
	}

	setupRoute(engines)
	master_gin.Wait()
}

func setupRoute(engines []*gin.Engine)  {
	for _, e := range engines {
		e.GET("/", func(context *gin.Context) {
			context.String(200, "hello world!\r\n")
		})
	}
}