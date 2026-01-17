package backend

import (
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

func IsHostUp(ip string) bool {
	pinger, err := probing.NewPinger(ip)
	if err != nil {
		return false
	}
	pinger.SetPrivileged(false)
	pinger.Count = 1
	pinger.Timeout = 2 * time.Second

	err = pinger.Run()
	if err != nil {
		return false
	}

	stats := pinger.Statistics()
	return stats.PacketLoss == 0
}
