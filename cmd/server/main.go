package main

import (
	"net"
	"net/http"
	"net/rpc"

	"github.com/umayr/zero"
)

func main() {
	rpc.RegisterName("store", zero.NewPlug())
	rpc.HandleHTTP()

	l, err := net.Listen("tcp", ":7161")
	if err != nil {
		panic(err)
	}

	go http.Serve(l, nil)

	select {}
}
