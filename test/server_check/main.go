package main

import (
	"fmt"
	"labor-app/cmd/server_check"
	"time"
)

func main() {
	hosts := []string{"Yuu.local", "phil.local"}

	for _, host := range hosts {
		serverOnline, ping := server_check.PingOS(host, 5*time.Second)
		if serverOnline {
			fmt.Printf("✅ %s đang ONLINE, ping: %s \n", host, ping)
		} else {
			fmt.Printf("❌ %s đang OFFLINE (hoặc chặn Ping)\n", host)
		}
	}
}
