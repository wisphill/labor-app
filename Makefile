APP_NAME := labor-app
TARGET_DIR := target
BINARY := $(TARGET_DIR)/$(APP_NAME)
APP_BUNDLE := LaborApp.app
APP_BUNDLE_MACOS := $(APP_BUNDLE)/Contents/MacOS

.PHONY: run build clean app install dmg pkg

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

# Create DMG installer (standard macOS distribution)
dmg: app
	@echo "📦 Creating DMG installer..."
	mkdir -p $(TARGET_DIR)
	rm -f $(TARGET_DIR)/LaborApp.dmg
	mkdir -p /tmp/laborapp-dmg
	cp -r $(APP_BUNDLE) /tmp/laborapp-dmg/
	ln -s /Applications /tmp/laborapp-dmg/Applications 2>/dev/null || true
	hdiutil create -volname "LaborApp" -srcfolder /tmp/laborapp-dmg -ov -format UDZO $(TARGET_DIR)/LaborApp.dmg
	rm -rf /tmp/laborapp-dmg
	@echo "✓ DMG installer created: $(TARGET_DIR)/LaborApp.dmg"

# Create PKG installer (traditional "next next" installer with UI)
pkg: app
	@echo "📦 Creating PKG installer..."
	mkdir -p $(TARGET_DIR)/LaborApp-pkg/Applications
	cp -r $(APP_BUNDLE) $(TARGET_DIR)/LaborApp-pkg/Applications/
	pkgbuild --root $(TARGET_DIR)/LaborApp-pkg \
		--identifier com.wisphill.laborapp \
		--version 1.0 \
		--install-location / \
		$(TARGET_DIR)/LaborApp-Installer.pkg
	rm -rf $(TARGET_DIR)/LaborApp-pkg
	@echo "✓ PKG installer created: $(TARGET_DIR)/LaborApp-Installer.pkg"
	@echo "👉 Double-click to install with UI!"
