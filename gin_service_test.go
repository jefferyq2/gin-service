package master_gin

import (
	"fmt"
	"testing"
)

func TestInit(t *testing.T) {
	var addresses = "127.0.0.1:8088"
	_, err := Init(addresses)
	if err != nil {
		fmt.Printf("init error %s\r\n", err.Error())
		return
	}
}
