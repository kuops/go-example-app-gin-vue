package ip

import (
	"github.com/kuops/go-example-app/server/pkg/log"
	"net"
)

func GetIPAddress() string {
	var ip string
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		log.Error("get ip address failed")
	}

	defer conn.Close()

	if conn != nil {
		ip = conn.LocalAddr().(*net.UDPAddr).IP.String()
	}

	return ip
}