package utils

import (
	"net"
	"strings"
)

//获取本机对外ip地址
func GetOutboundIP() (ip string, err error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return
	}
	defer conn.Close()
	localAdder := conn.LocalAddr().(*net.UDPAddr)
	//fmt.Println(localAdder.String())
	ip = strings.Split(localAdder.IP.String(), ":")[0]
	return
}
