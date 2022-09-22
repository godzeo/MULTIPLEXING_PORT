package root

import "fmt"

func WinDel() { //删除转发规则
	Cmd("netsh interface portproxy reset")
	fmt.Println("The forwarding rule is cleared successfully!")
}

func LinuxDel() { //删除转发规则

	LinuxCmd("iptables -t nat -D PREROUTING 1")
	LinuxCmd("iptables -t nat -L -n --line-numbers")
	fmt.Println("[Info] The forwarding rule is cleared successfully!")
}
