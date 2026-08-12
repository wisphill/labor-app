This is the flow to change the behavior of the Cocoa Macos NSWindow

```bash
main.go
│
├── C.installWindowCentering()
│       │
│       └── Cocoa / Objective-C runtime
│              │
│              └── hook NSWindow.orderFront
│
├── new(app.Window)
│       │
│       ▼
│      Gio
│       │
│       ▼
│   create NSWindow
│
└── app.Main()
        │
        ▼
   Gio want to show window
        │
        ▼
   NSWindow.orderFront()
        │
        ▼
   [HOOK]
        │
        ├── get NSScreen
        ├── get window size
        ├── calculate center
        └── setFrameOrigin()
        │
        ▼
   original orderFront()
        │
        ▼
   Window appears
```
