package main

import (
	"flag"
	"fmt"
	"github.com/acl-dev/master-gin"
	"github.com/acl-dev/master-go"
	"github.com/gin-gonic/gin"
	"net/http"
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
	fmt.Println("master-go version:", master.Version)
	fmt.Println("mager-gin version:", master_gin.Version)

	flag.Parse()

	if !*debugMode {
		gin.SetMode(gin.ReleaseMode)
	}

	fmt.Println("addr:", *addresses)
	service, err := master_gin.Init(*addresses, onStop)
	if err != nil {
		fmt.Println("Init master gin service failed:", err)
		return
	}

	setRoute(*service)
	fmt.Println("Listen and running ...")
	service.Run()
}

func ginWrap(f func(w http.ResponseWriter, r *http.Request)) gin.HandlerFunc {
	return func(context *gin.Context) {
		f(context.Writer, context.Request)
	}
}

func onTest(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintf(w, "test: Hello World!\r\n")
}

func setRoute(service master_gin.GinService) {
	for _, s := range service.Servers {
		s.Engine.GET("/", func(context *gin.Context) {
			context.String(200, "hello world!\r\n")
		})
		s.Engine.GET("/test", ginWrap(onTest))
	}
}

func onStop(bool) {
	fmt.Println("The process stopped!")
}