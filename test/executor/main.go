package main

import (
	"fmt"
	executor "labor-app/cmd/execute_commands"
)

func main() {
	// Cách 1: Truyền trực tiếp danh sách tham số variadic
	out1, err := executor.ExecuteCommands(
		"echo === STEP 1 ===",
		"uptime",
		"echo === STEP 2 ===",
	)
	if err != nil {
		fmt.Println("Lỗi:", err)
	}
	fmt.Println(out1)

	// Cách 2: Truyền từ struct TerminalScript (commands []string) trong app Gio UI của bạn
	cmds := []string{
		"echo 'Shutting down service...'",
		"echo 'Done!'",
	}
	out2, _ := executor.ExecuteCommands(cmds...)
	fmt.Println(out2)
}
