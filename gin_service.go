package master_gin

import (
	"github.com/pkg/errors"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/acl-dev/go-service"
	"github.com/gin-gonic/gin"
)

var (
	g sync.WaitGroup	// Used to wait for service to stop.
	Version = "1.0.0"
)
type AcceptFunc func(net.Conn)
type CloseFunc func(net.Conn)

type GinServ struct {
	Listener net.Listener
	Engine *gin.Engine
}

type GinService struct {
	Alone bool
	Servers []*GinServ
	AcceptHandler AcceptFunc
	CloseHandler  CloseFunc
}

// Run begin to start all the listening servers after Init() called.
func (service *GinService) Run()  {
	g.Add(len(service.Servers))
	for _, s := range service.Servers {
		service.startServer(s.Listener, s.Engine)
	}

	g.Wait()
}

func (service *GinService) startServer(listener net.Listener, engine *gin.Engine)  {
	go func() {
		defer g.Done()

		//_ = engine.RunListener(listener)
		server := &http.Server {
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
		server.Serve(listener)
	}()
}

// Init Create the gin service when starting in alone or daemon mode;
// The addresses must not be empty in alone mode, and will be ignored in daemon mode,
// the addresses' format lookup like "127.0.0.1:8080;127.0.0.1:8081;127.0.0.1:8082";
// The stopHandler is the callback when the process is exiting.
func Init(addresses string, stopHandler func(bool)) (*GinService, error) {
	master.Prepare()

	if master.Alone && len(addresses) == 0 {
		log.Println("Listening addresses shouldn't be empty in running alone mode!")
		return nil, errors.New("Listening addresses shouldn't be empty in alone mode")
	}

	listeners, err := master.ServiceInit(addresses, stopHandler)
	if err != nil {
		return nil, err
	}

	var service GinService

	for _, l := range listeners {
		engine := gin.Default()
		serv := &GinServ{ Listener: l, Engine: engine }
		service.Servers = append(service.Servers, serv)
	}

	return &service, nil
}
