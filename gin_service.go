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

		engine.RunListener(listener)
	}()
}

func Init(addrs string) (*GinService, error) {
	master.Prepare()

	if master.Alone && len(addrs) == 0 {
		log.Println("Listening addresses shouldn't be empty in running alone mode!")
		return nil, errors.New("Listening addresses shouldn't be empty in alone mode")
	}

	if !master.Alone {
		addrs = ""
	}

	listeners, err := master.ServiceInit(addrs, nil)
	if err != nil {
		return nil, err
	}

	var engines []*gin.Engine
	for i := 0; i < len(listeners); i++ {
		engine := gin.Default()
		engines = append(engines, engine)
	}
	service := &GinService{ Alone: master.Alone, Listeners: listeners, Engines: engines }
	return service, nil
}
