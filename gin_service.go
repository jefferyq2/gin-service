package master_gin

import (
	"github.com/pkg/errors"
	"log"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/acl-dev/go-service"
	"github.com/gin-gonic/gin"
)

var (
	Version = "1.0.1"
)

type AcceptFunc func(net.Conn)
type CloseFunc func(net.Conn)

type GinServ struct {
	Listener net.Listener
	Engine   *gin.Engine
}

type GinService struct {
	Alone         bool
	Servers       []*GinServ
	AcceptHandler AcceptFunc
	CloseHandler  CloseFunc
}

// Run begin to start all the listening servers after Init() called.
func (service *GinService) Run() {
	var g sync.WaitGroup // Used to wait for all service to stop.

	g.Add(len(service.Servers))

	for _, s := range service.Servers {
		go func(serv *GinServ) {
			defer g.Done()
			service.run(serv.Listener, serv.Engine)
		}(s)
	}

	// Waiting the disconnect status from acl_master.
	master.Wait()

	// Waiting all the gin services to stop.
	g.Wait()
}

func (service *GinService) run(listener net.Listener, engine *gin.Engine) {
	//_ = engine.RunListener(listener)
	server := &http.Server{
		Handler: engine,
		ConnState: func(conn net.Conn, state http.ConnState) {
			switch state {
			case http.StateNew:
				master.ConnCountInc()
				if service.AcceptHandler != nil {
					service.AcceptHandler(conn)
				}
				break
			case http.StateActive:
				break
			case http.StateIdle:
				break
			case http.StateHijacked:
				master.ConnCountDec()
				if service.CloseHandler != nil {
					service.CloseHandler(conn)
				}
				break
			case http.StateClosed:
				master.ConnCountDec()
				if service.CloseHandler != nil {
					service.CloseHandler(conn)
				}
				break
			default:
				break
			}
		},
	}

	err := server.Serve(listener)
	if err != nil {
		log.Printf("pid=%d: Http Service over, err=%s\r\n",
			os.Getpid(), err.Error())
	}
}

// Init Create the gin service when starting in alone or daemon mode;
// The addresses must not be empty in alone mode, and will be ignored in daemon mode,
// the addresses' format lookup like "127.0.0.1:8080;127.0.0.1:8081;127.0.0.1:8082";
func Init(addresses string) (*GinService, error) {
	master.Prepare()

	if master.Alone && len(addresses) == 0 {
		log.Println("Listening addresses shouldn't be empty in running alone mode!")
		return nil, errors.New("Listening addresses shouldn't be empty in alone mode")
	}

	listeners, err := master.ServiceInit(addresses)
	if err != nil {
		return nil, err
	}

	var service GinService

	for _, l := range listeners {
		engine := gin.Default()
		serv := &GinServ{Listener: l, Engine: engine}
		service.Servers = append(service.Servers, serv)
	}

	return &service, nil
}
