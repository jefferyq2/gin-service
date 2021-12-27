package master_gin

import (
	"log"
	"net"
	"sync"

	"github.com/pkg/errors"
	"github.com/acl-dev/master-go"
	"github.com/gin-gonic/gin"
)

var (
	g sync.WaitGroup	// Used to wait for service to stop.
	Alone bool
)

func Init(addrs string) ([]*gin.Engine, error) {
	master.Prepare()
	Alone = master.Alone
}

func AloneStart(addrs string) ([]*gin.Engine, error) {
	if len(addrs) == 0 {
		log.Println("Addrs empty")
		return nil, errors.New("No listen addresses")
	}

	return ginServiceStart(addrs)
}

func DaemonStart() ([]*gin.Engine, error) {
	return ginServiceStart("")
}

func Wait()  {
	g.Wait()
}

func ginServiceStart(addrs string) ([]*gin.Engine, error) {
	var listeners []net.Listener
	var err error

	listeners, err = master.ServiceInit(addrs, nil)
	if err != nil {
		return nil, err
	}

	g.Add(len(listeners))

	var engines []*gin.Engine

	for _, ln := range listeners {
		engine := gin.Default()
		engines = append(engines, engine)
		startServer(ln, engine)
	}

	return engines, nil
}

func startServer(listener net.Listener, engine *gin.Engine)  {
	go func() {
		defer g.Done()

		engine.RunListener(listener)
	}()
}
