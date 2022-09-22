package root

import (
	"fmt"
	"os"
	"strings"
)

func WinAdd(listenport string, address string, connectport string, connectaddress string) { //添加转发规则
	cmdstr := Cmd("netsh interface portproxy show all")
	if strings.Contains(cmdstr, connectaddress) && strings.Contains(cmdstr, connectport) && strings.Contains(cmdstr, connectaddress) && strings.Contains(cmdstr, listenport) {
		fmt.Println("[Info] listen_ip:", connectaddress, " 端口:", listenport)
		fmt.Println("[Info] forward_ip:", connectaddress, " 转发端口:", connectport)
		fmt.Println("[Info] 该规则已经添加过了！")
		return
	} else {
		Cmd("netsh interface portproxy add v4tov4 listenport=" + listenport + " listenaddress=" + connectaddress + " connectport=" + connectport + " connectaddress=" + connectaddress)
		c1 := Cmd("netsh interface portproxy show all")
		//fmt.Println(c1)
		if strings.Contains(c1, connectaddress) == true && strings.Contains(c1, connectport) {
			fmt.Println("[Info] listen_ip:", address, " 端口:", listenport)
			fmt.Println("[Info] forward_ip:", connectaddress, " 转发端口:", connectport)
			fmt.Println("[Info] success！")
		} else {
			fmt.Println("[Info] 添加失败！")
			os.Exit(1)
		}
	}
}
