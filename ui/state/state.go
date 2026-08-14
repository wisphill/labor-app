package state

import (
	"context"
	"fmt"
	server "labor-app/cmd/host"
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
	Name         string
	Address      string
	IsOnline     bool
	PingRTT      time.Duration
	ServerSignal chan bool

	BtnPower widget.Clickable
	Wsls     []*WSLState
	Mu       sync.Mutex
}

type WSLState struct {
	Name     string
	BtnPower widget.Clickable
}

func (host *HostState) PingToServerLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			var wg sync.WaitGroup
			host.Mu.Lock()
			addr := host.Address
			host.Mu.Unlock()

			wg.Add(1)
			go func(h *HostState, address string) {
				defer wg.Done()
				online, rtt := server.PingOS(address, 1500*time.Millisecond)
				h.Mu.Lock()
				h.IsOnline = online
				h.PingRTT = rtt
				h.Mu.Unlock()
			}(host, addr)
			wg.Wait()

			time.Sleep(3 * time.Second)
		}
	}
}

func (host *HostState) FetchWSLNodesLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("yeah, return, no more loops!")
			return
		default:
			var wg sync.WaitGroup
			host.Mu.Lock()
			addr := host.Address
			host.Mu.Unlock()

			wg.Add(1)
			go func(h *HostState, address string) {
				defer wg.Done()
				wslNodes, err := server.GetRunningWSLNodes()
				h.Mu.Lock()
				h.Wsls = make([]*WSLState, 0)
				if err != nil {
					fmt.Println("Error while getting the WSL nodes")
					h.Mu.Unlock()
					return
				}

				for _, wslNode := range wslNodes {
					h.Wsls = append(h.Wsls, &WSLState{
						Name: wslNode,
					})
				}

				h.Mu.Unlock()
			}(host, addr)
			wg.Wait()
			time.Sleep(3 * time.Second)
		}
	}
}

func (host *HostState) HandleServerSignal(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("yeah, return, no more loops!")
			return
		case signal, ok := <-host.ServerSignal:
			if !ok {
				fmt.Println("ServerSignal closed")
				return
			}

			if signal == false {
				fmt.Println("Turn off the server nowwww")
				server.TurnOffServer()
			} else {
				fmt.Println("Turn on the server nowwww")
				server.TurnOnServer()
			}
		}
	}
}
