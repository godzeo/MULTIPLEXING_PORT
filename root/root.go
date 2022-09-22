package root

import (
	"github.com/rs/zerolog"
	"multiplexing_port_socks5/protocols"
	"multiplexing_port_socks5/socks5"
)

func WinRun(listenport string, address string, connectport string, connectaddress string, socketport string, p []*protocols.Protocol, logger zerolog.Logger) {

	//添加转发规则
	WinAdd(listenport, address, connectport, connectaddress)

	//创建一个socks5
	go socks5.CreateForwardSocks(socketport)
	protocols.RunServer(connectaddress+":"+connectport, p, logger)
}
