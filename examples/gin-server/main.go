package main

import (
	"flag"
	"fmt"
	"github.com/acl-dev/gin-service"
	"github.com/acl-dev/go-service"
	"github.com/gin-gonic/gin"
	"log"
	"net"
	"net/http"
	"os"
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

func main() {
	fmt.Println("go-service version:", master.Version)
	fmt.Println("gin-service version:", master_gin.Version)

	flag.Parse()

	if !*debugMode {
		gin.SetMode(gin.ReleaseMode)
	}

	fmt.Println("addr:", *addresses)
	service, err := master_gin.Init(*addresses)
	if err != nil {
		log.Println("Init master gin service failed:", err)
		return
	}

	service.AcceptHandler = func(conn net.Conn) {
		log.Printf("Connect from %s\r\n", conn.RemoteAddr())
	}
	service.CloseHandler = func(conn net.Conn) {
		log.Printf("Disconnect from %s\r\n", conn.RemoteAddr())
	}

	fmt.Printf("ServiceType=%s, test_src=%s, test_bool=%t\r\n",
		master.ServiceType, master.AppConf.GetString("test_src"),
		master.AppConf.GetBool("test_bool"))

	setRoute(*service)
	log.Printf("pid=%d, start gin service...\r\n", os.Getpid())
	service.Run()
	log.Printf("pid=%d, gin service stopped!\r\n", os.Getpid())
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
