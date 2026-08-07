package state

import (
	"labor-app/ui/components"

	"gioui.org/widget"
)

type AppState struct {
	NameInput   widget.Editor
	SelectedTab components.Tab
	BtnServer   widget.Clickable
	BtnWSL      widget.Clickable
}
