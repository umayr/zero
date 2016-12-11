package main

import (
	"net/rpc"
	"net"
	"net/http"
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