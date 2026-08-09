package main

import (
	"fmt"
	executor "labor-app/cmd/execute_commands"
)

func main() {
	out1, err := executor.ExecuteCommands(
		"echo === STEP 1 ===",
		"uptime",
		"echo === STEP 2 ===",
	)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println(out1)

	cmds := []string{
		"echo 'Shutting down service...'",
		"echo 'Done!'",
	}
	out2, _ := executor.ExecuteCommands(cmds...)
	fmt.Println(out2)
}
