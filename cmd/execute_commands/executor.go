package executor

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// use the OS shell
func ExecuteCommands(commands ...string) (string, error) {
	if len(commands) == 0 {
		return "", nil
	}

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		joinedCmds := strings.Join(commands, " & ")
		cmd = exec.Command("cmd", "/C", joinedCmds)
	default:
		joinedCmds := strings.Join(commands, " && ")
		cmd = exec.Command("sh", "-c", joinedCmds)
	}

	// Execute and get the stdout and stderr
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("Error while running commands [%v]: %w", err, err)
	}

	return string(output), nil
}
