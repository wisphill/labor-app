package state

import (
	"sync"
	"time"

	"gioui.org/widget"
)

type HostAction int

const (
	HostActionTurnOn HostAction = iota
	HostActionShutdown
)

type TerminalScript struct {
	Action   HostAction
	Commands []string // commands for running the script
}

type HostState struct {
	Name     string
	Address  string
	IsOnline bool
	PingRTT  time.Duration

	BtnPower widget.Clickable
	Wsls     []*WSLState
	Mu       sync.Mutex
}

type WSLState struct {
	Name     string
	BtnPower widget.Clickable
}
