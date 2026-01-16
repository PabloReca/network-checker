package backend

import (
	"net"
	"time"
)

func IsHostUp(ip string) bool {
	conn, err := net.DialTimeout("tcp", ip+":80", 1*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
