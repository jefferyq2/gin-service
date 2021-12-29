package main

import (
	"flag"
	"fmt"
	"github.com/acl-dev/master-gin"
	"github.com/gin-gonic/gin"
)

var (
	addresses = flag.String(
		"listen",
		"127.0.0.1:8088",
		"listen addresses in alone mode",
	)
	debugMode = flag.Bool(
		"debug",
		false,
		"If running gin server in debug mode",
	)
)

func main()  {
	flag.Parse()

	if !*debugMode {
		gin.SetMode(gin.ReleaseMode)
	}

	service, err := master_gin.Init(*addresses, onStop)
	if err != nil {
		fmt.Println("Init master gin service failed:", err)
		return
	}

	setupRoute(service.Engines)
	fmt.Println("Listen and running ...")
	service.Run()
}

func setupRoute(engines []*gin.Engine)  {
	for _, e := range engines {
		e.GET("/", func(context *gin.Context) {
			context.String(200, "hello world!\r\n")
		})
	}
}

func onStop(bool) {
	fmt.Println("The process stopped!")
}