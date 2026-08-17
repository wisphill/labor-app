APP_NAME := labor-app
TARGET_DIR := target
BINARY := $(TARGET_DIR)/$(APP_NAME)
APP_BUNDLE := LaborApp.app
APP_BUNDLE_MACOS := $(APP_BUNDLE)/Contents/MacOS

.PHONY: run build clean app install

run:
	CGO_ENABLED=1 go run .

build:
	mkdir -p $(TARGET_DIR)
	CGO_ENABLED=1 go build -o $(BINARY) .

app: build
	mkdir -p $(APP_BUNDLE_MACOS)
	mkdir -p $(APP_BUNDLE)/Contents
	cp $(BINARY) $(APP_BUNDLE_MACOS)/$(APP_NAME)
	cp bundle/Info.plist $(APP_BUNDLE)/Contents/Info.plist
	@echo "✓ LaborApp.app built successfully"

install: app
	mkdir -p ~/Applications
	cp -r $(APP_BUNDLE) ~/Applications/$(APP_BUNDLE)
	@echo "✓ LaborApp.app installed to ~/Applications"

start: build
	./$(BINARY)

clean:
	rm -f $(APP_NAME)
	go clean
