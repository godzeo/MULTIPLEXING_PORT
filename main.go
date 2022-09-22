package main

import (
	"flag"
	"fmt"
	"github.com/rs/zerolog"
	"multiplexing_port_socks5/protocols"
	"multiplexing_port_socks5/root"
	"os"
	"runtime"
)

func main() {

	fmt.Println("-h 查看使用参数")
	//fmt.Println("ex:  mute.exe -type add -lhost 0.0.0.0 -lport 80 -rhost 10.211.55.4 -rport 8889 -socks5 1234")
	//fmt.Println("本地开启一个socks5代理,端口 1234")
	//fmt.Println("将 80 端口的流量转发到 8889 端口")
	//fmt.Println("然后在 8889 端口进行流量分流：http请求就走 127.0.0.1:80 端口,socks5请求就走 127.0.0.1:1234 端口")

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	p := make([]*protocols.Protocol, 0, 7)
	socketport := flag.String("socks5", "", "socks5开启端口")
	cmdtype := flag.String("type", "", "执行类型(add/del)")
	listenport := flag.String("lport", "", "本地监听的端口")
	address := flag.String("lhost", "", "本地监听的ip地址")
	connectport := flag.String("rport", "", "需要转发到的端口")
	connectaddress := flag.String("rhost", "", "需要转发到的地址")
	ssh := flag.String("ssh", "127.0.0.1:22", "The SSH server address")
	tls := flag.String("tls", "127.0.0.1:443", "The TLS/HTTPS server address")
	openvpn := flag.String("ovpn", "", "The OpenVPN server address")
	rdp := flag.String("rdp", "127.0.0.1:3389", "The OpenVPN server address")

	flag.Parse()
	if *cmdtype == "add" {
		if *listenport == "" {
			fmt.Println("[Info] 请填写转发的端口！")
		} else if *socketport == "" {
			fmt.Println("[Info] 请填写socks5开启的端口！")
		} else if *address == "" {
			fmt.Println("[Info] 请填写转发的ip地址！")
		} else if *connectport == "" {
			fmt.Println("[Info] 请填写转发的目标端口！")
		} else if *connectaddress == "" {
			fmt.Println("[Info] 请填写转发的目标ip地址！")
		} else {
			//转发规则
			p = append(p, protocols.NewHTTPProtocol(*address+":"+*listenport))
			p = append(p, protocols.NewSOCKS5Protocol(*connectaddress+":"+*socketport))
			if *tls != "" {
				p = append(p, protocols.NewTLSProtocol(*tls))
			}
			if *ssh != "" {
				p = append(p, protocols.NewSSHProtocol(*ssh))
			}
			if *openvpn != "" {
				p = append(p, protocols.NewOpenVPNProtocol(*openvpn))
			}
			if *rdp != "" {
				p = append(p, protocols.NewRDPProtocol(*rdp))
			}

			sysType := runtime.GOOS
			if sysType == "windows" {
				root.WinRun(*listenport, *address, *connectport, *connectaddress, *socketport, p, logger)

			} else if sysType == "linux" {
				root.LinuxRun(*listenport, *address, *connectport, *connectaddress, *socketport, p, logger)
			} else {
				fmt.Println("[Info] System recognition err")
				os.Exit(1)
			}

		}
	} else if *cmdtype == "del" {
		sysType := runtime.GOOS
		if sysType == "windows" {
			root.WinDel()

		} else if sysType == "linux" {
			root.LinuxDel()
		} else {
			fmt.Println("[Info] System recognition err")
			os.Exit(1)
		}
	}
}
