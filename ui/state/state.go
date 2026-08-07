package state

import (
	"labor-app/ui/components"
	"sync"
	"time"

	"gioui.org/widget"
)

type HostAction int

const (
	HostActionTurnOn HostAction = iota
	HostActionShutdown
)

type AppState struct {
	NameInput   widget.Editor
	SelectedTab components.Tab
	BtnServer   widget.Clickable
	BtnWSL      widget.Clickable
}

type HostState struct {
	Name     string
	Address  string
	IsOnline bool
	PingRTT  time.Duration

	TerminalScripts struct {
		action   HostAction
		commands []string // commands for running the script
	}

	Mu sync.Mutex
}
