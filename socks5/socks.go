package socks5

import (
	"fmt"
	socks5 "github.com/armon/go-socks5"
	"github.com/hashicorp/yamux"
	"multiplexing_port_socks5/protocols"
)

var session *yamux.Session

func CreateForwardSocks2(socketport string) error {
	protocols.Slave("0.0.0.0", socketport)
	return nil
}

func CreateForwardSocks(address string) error {
	server, err := socks5.New(&socks5.Config{})
	if err != nil {
		return err
	}
	fmt.Println("[Info] socks5 success！port:", address)
	if err := server.ListenAndServe("tcp", "0.0.0.0:"+address); err != nil {
		return err
	}
	return nil
}
