//go:build darwin

package darwin

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

const Label = "com.myapp.laboratory"

func plistPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error while getting home dir")
	}

	fmt.Println(filepath.Join(
		home,
		"Library",
		"LaunchAgents",
		Label+".plist",
	))

	return filepath.Join(
		home,
		"Library",
		"LaunchAgents",
		Label+".plist",
	), nil
}

func EnableStartAtLogin() error {
	executable, err := os.Executable()
	if err != nil {
		return err
	}

	path, err := plistPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	uid := strconv.Itoa(os.Getuid())
	domain := "gui/" + uid

	plist := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC
 "-//Apple//DTD PLIST 1.0//EN"
 "http://www.apple.com/DTDs/PropertyList-1.0.dtd">

<plist version="1.0">
<dict>

	<key>Label</key>
	<string>%s</string>

	<key>ProgramArguments</key>
	<array>
		<string>%s</string>
	</array>

	<key>RunAtLoad</key>
	<true/>

</dict>
</plist>
`, Label, executable)

	if err := os.WriteFile(path, []byte(plist), 0644); err != nil {
		fmt.Printf("Error while writing plist files")
		return err
	}

	// Remove old registration if it exists.
	err = exec.Command(
		"launchctl",
		"bootout",
		domain,
		path,
	).Run()
	if err != nil {
		fmt.Printf("Error while removing old registration %v\n", err)
	}

	// Register the new LaunchAgent.
	cmd := exec.Command(
		"launchctl",
		"bootstrap",
		domain,
		path,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(
			"launchctl bootstrap: %w: %s",
			err,
			output,
		)
	}

	return nil
}

func DisableStartAtLogin() error {
	path, err := plistPath()
	if err != nil {
		return err
	}

	uid := strconv.Itoa(os.Getuid())
	domain := "gui/" + uid

	_ = exec.Command(
		"launchctl",
		"bootout",
		domain,
		path,
	).Run()

	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}

func IsStartAtLoginEnabled() bool {
	path, err := plistPath()
	if err != nil {
		return false
	}
	_, err = os.Stat(path)
	return err == nil
}
