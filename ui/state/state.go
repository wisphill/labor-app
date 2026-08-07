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

type TerminalScript struct {
	Action   HostAction
	Commands []string // commands for running the script
}

type HostState struct {
	Name     string
	Address  string
	IsOnline bool
	PingRTT  time.Duration

	// Scripts cho 2 hành động Bật & Tắt
	TurnOnScript   TerminalScript
	ShutdownScript TerminalScript

	// State cho Gio UI Widget
	BtnTurnOn   widget.Clickable
	BtnShutdown widget.Clickable

	Mu sync.Mutex
}
