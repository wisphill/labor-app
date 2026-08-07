package server_check

import (
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"time"
)

var timeRegex = regexp.MustCompile(`(?i)time[=<]\s*([\d.]+)\s*ms`)

// PingOS sử dụng binary ping của HĐH: Không cần Root, Không cần mở Port
func PingOS(host string, timeout time.Duration) (bool, time.Duration) {
	var cmd *exec.Cmd
	timeoutMs := strconv.Itoa(int(timeout.Milliseconds()))

	switch runtime.GOOS {
	case "windows":
		// -n 1 (1 gói), -w timeout (ms)
		cmd = exec.Command("ping", "-n", "1", "-w", timeoutMs, host)
	case "darwin":
		// macOS: -c 1 (1 gói), -W timeout (ms)
		cmd = exec.Command("ping", "-c", "1", "-W", timeoutMs, host)
	default:
		// Linux: -c 1 (1 gói), -W timeout (giây)
		timeoutSec := strconv.Itoa(int(timeout.Seconds()))
		if timeoutSec == "0" {
			timeoutSec = "1"
		}
		cmd = exec.Command("ping", "-c", "1", "-W", timeoutSec, host)
	}

	out, err := cmd.Output()
	if err != nil {
		return false, 0
	}

	matches := timeRegex.FindStringSubmatch(string(out))
	if len(matches) > 1 {
		msFloat, err := strconv.ParseFloat(matches[1], 64)
		if err == nil {
			return true, time.Duration(msFloat * float64(time.Millisecond))
		}
	}

	return true, 0
}
