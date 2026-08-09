package server

import (
	"fmt"
	executor "labor-app/cmd/execute_commands"
	"os"
	"strings"
)

func TurnOffServer() {
	_, err := executor.ExecuteCommands(
		`ssh Windows@yuu "shutdown /s /t 0"`,
	)
	if err != nil {
		fmt.Printf("Error while processing shutting down %v", err)
	}
}

func TurnOnServer() {
	telegramBotToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	command := fmt.Sprintf(
		`curl -s -X POST "https://api.telegram.org/bot%s/sendMessage" -d chat_id="-5115557042" -d text="/wake"`,
		telegramBotToken,
	)
	_, err := executor.ExecuteCommands(command)
	if err != nil {
		fmt.Printf("Error while processing turning the server on %v", err)
	}
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
