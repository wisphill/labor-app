package executor

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// ExecuteCommands nhận vào danh sách nhiều lệnh, ghép lại và thực thi qua Shell gốc của hệ điều hành.
// Trả về: (Output thu được, Lỗi nếu có)
func ExecuteCommands(commands ...string) (string, error) {
	if len(commands) == 0 {
		return "", nil
	}

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		// Trên Windows dùng & hoặc && để nối chuỗi lệnh
		joinedCmds := strings.Join(commands, " & ")
		cmd = exec.Command("cmd", "/C", joinedCmds)

	default: // macOS & Linux
		// Trên Unix-like dùng && để đảm bảo lệnh trước chạy thành công mới chạy lệnh sau
		joinedCmds := strings.Join(commands, " && ")
		cmd = exec.Command("sh", "-c", joinedCmds)
	}

	// Thực thi và thu thập đồng thời cả stdout và stderr
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("lỗi thực thi kịch bản [%v]: %w", err, err)
	}

	return string(output), nil
}

func main() {
	// Cách 1: Truyền trực tiếp danh sách tham số variadic
	out1, err := ExecuteCommands(
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
	out2, _ := ExecuteCommands(cmds...)
	fmt.Println(out2)
}
