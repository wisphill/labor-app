#!/bin/bash

# LaborApp Installer Script
# This creates a professional PKG installer with custom UI

APP_NAME="LaborApp"
VERSION="1.0"
BUNDLE_ID="com.wisphill.laborapp"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
TARGET_DIR="$PROJECT_ROOT/target"
APP_BUNDLE="$PROJECT_ROOT/LaborApp.app"

# Create temporary build directory
TEMP_PKG_DIR=$(mktemp -d)
trap "rm -rf $TEMP_PKG_DIR" EXIT

echo "📦 Building $APP_NAME Installer..."

# Prepare the payload (what will be installed)
PAYLOAD_DIR="$TEMP_PKG_DIR/payload"
mkdir -p "$PAYLOAD_DIR/Applications"
cp -r "$APP_BUNDLE" "$PAYLOAD_DIR/Applications/"

# Create pre-install script (runs before installation)
mkdir -p "$TEMP_PKG_DIR/scripts"
cat > "$TEMP_PKG_DIR/scripts/preinstall" << 'EOF'
#!/bin/bash
# Pre-installation checks
echo "Preparing installation..."
EOF
chmod +x "$TEMP_PKG_DIR/scripts/preinstall"

# Create post-install script (runs after installation)
cat > "$TEMP_PKG_DIR/scripts/postinstall" << 'EOF'
#!/bin/bash
# Post-installation steps
echo "Installation complete!"
# Optional: Launch the app
# open /Applications/LaborApp.app
EOF
chmod +x "$TEMP_PKG_DIR/scripts/postinstall"

# Build the PKG
mkdir -p "$TARGET_DIR"
pkgbuild \
    --root "$PAYLOAD_DIR" \
    --scripts "$TEMP_PKG_DIR/scripts" \
    --identifier "$BUNDLE_ID" \
    --version "$VERSION" \
    --install-location / \
    "$TARGET_DIR/${APP_NAME}-Installer.pkg"

if [ $? -eq 0 ]; then
    echo "✅ Installer created successfully!"
    echo "📍 Location: $TARGET_DIR/${APP_NAME}-Installer.pkg"
    ls -lh "$TARGET_DIR/${APP_NAME}-Installer.pkg"
else
    echo "❌ Failed to create installer"
    exit 1
fi
