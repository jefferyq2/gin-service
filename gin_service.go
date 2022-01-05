package master_gin

import (
	"github.com/pkg/errors"
	"log"
	"net"
	"sync"

	"github.com/acl-dev/master-go"
	"github.com/gin-gonic/gin"
)

var (
	g sync.WaitGroup	// Used to wait for service to stop.
)

type GinService struct {
	Alone bool
	Listeners []net.Listener
	Engines []*gin.Engine
}

// Run begin to start all the listening servers after Init() called.
func (service *GinService) Run()  {
	g.Add(len(service.Listeners))
	for i := 0; i < len(service.Engines); i++ {
		startServer(service.Listeners[i], service.Engines[i])
	}

	g.Wait()
}

func startServer(listener net.Listener, engine *gin.Engine)  {
	go func() {
		defer g.Done()

		_ = engine.RunListener(listener)
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

	var engines []*gin.Engine
	for i := 0; i < len(listeners); i++ {
		engine := gin.Default()
		engines = append(engines, engine)
	}
	service := &GinService{
		Alone: master.Alone,
		Listeners: listeners,
		Engines: engines,
	}
	return service, nil
}
