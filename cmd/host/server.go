package server

import (
	"fmt"
	executor "labor-app/cmd/execute_commands"
	"os"
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
