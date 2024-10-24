package main

import (
	"log"

	"github.com/48club/rpc-watchdog/message"
	"github.com/48club/rpc-watchdog/service"
	"github.com/48club/rpc-watchdog/types"
)

func main() {
	rpcChan := make(chan types.Chan)
	service.Watch(rpcChan)

	for {
		c := <-rpcChan
		if c.Err != nil {
			log.Printf("RPC [%s] is not working: %s", c.Rpc, c.Err)
			message.Notify(c)
		} else {
			log.Printf("RPC [%s] is working", c.Rpc)
		}
	}
}
