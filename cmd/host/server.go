package server

import (
	"fmt"
	executor "labor-app/cmd/execute_commands"
	"net"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var timeRegex = regexp.MustCompile(`(?i)time[=<]\s*([\d.]+)\s*ms`)

// resolve host to the ipv4
func resolveHost(host string) string {
	host = strings.TrimSpace(host)

	if ip := net.ParseIP(host); ip != nil {
		return host
	}

	ips, err := net.LookupIP(host)
	if err != nil {
		return host
	}

	// get the IP v4
	for _, ip := range ips {
		if ipv4 := ip.To4(); ipv4 != nil {
			return ipv4.String()
		}
	}

	return host
}

func TurnOffServer() {
	_, err := executor.ExecuteCommands(
		`ssh Windows@yuu "shutdown /s /t 0"`,
	)
	if err != nil {
		fmt.Printf("Error while processing shutting down %v", err)
	}
}

func TurnOnServer() error {
	telegramBotToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if telegramBotToken == "" {
		return fmt.Errorf("Cannot find the telegram bot token. Configure at .labor_app/config")
	}
	command := fmt.Sprintf(
		`curl -s -X POST "https://api.telegram.org/bot%s/sendMessage" -d chat_id="-5115557042" -d text="/wake"`,
		telegramBotToken,
	)
	_, err := executor.ExecuteCommands(command)
	if err != nil {
		fmt.Printf("Error while processing turning the server on %v", err)
		return err
	}

	return nil
}

func GetRunningWSLNodes() ([]string, error) {
	output, err := executor.ExecuteCommands(`ssh Windows@yuu "wsl -l -v" | iconv -f UTF-16LE -t UTF-8 | sed '1d; s/^\* //'`)
	if err != nil {
		fmt.Printf("Error while getting the WSL nodes %v", err)
		return nil, err

	}

	var runningWSLNodes []string
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		// Ubuntu-24.04    Running    2
		if fields[1] == "Running" {
			runningWSLNodes = append(runningWSLNodes, fields[0])
		}
	}

	return runningWSLNodes, nil
}

func PingOS(host string, timeout time.Duration) (bool, time.Duration) {
	var cmd *exec.Cmd
	timeoutMs := strconv.Itoa(int(timeout.Milliseconds()))

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("ping", "-n", "1", "-w", timeoutMs, host)
	case "darwin":
		cmd = exec.Command("ping", "-c", "1", "-W", timeoutMs, host)
	default:
		timeoutSec := strconv.Itoa(int(timeout.Seconds()))
		if timeoutSec == "0" {
			timeoutSec = "1"
		}
		cmd = exec.Command("ping", "-c", "1", "-W", timeoutSec, host)
	}

	out, err := cmd.CombinedOutput()
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
