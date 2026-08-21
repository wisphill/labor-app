# PiggyBank Installer Guide

## Available Installer Options

### 1. **PKG Installer** (Recommended - "Next Next" UI)
Create a traditional macOS installer with a graphical UI:
```bash
make pkg
```
- ✅ Professional installer UI
- ✅ "Next Next" workflow
- ✅ Installs to `/Applications` (user only)
- ✅ Fast installation
- 📍 Output: `target/PiggyBank-Installer.pkg`

**How to use:**
1. Run `make pkg`
2. Open `target/PiggyBank-Installer.pkg`
3. Follow the installer prompts
4. App installs to `/Applications/PiggyBank.app`

---

### 2. **DMG Installer** (Standard macOS Distribution)
Create a disk image for distribution:
```bash
make dmg
```
- ✅ Standard macOS distribution method
- ✅ Drag & drop installation
- ✅ Professional look
- 📍 Output: `target/PiggyBank.dmg`

**How to use:**
1. Run `make dmg`
2. Open `target/PiggyBank.dmg`
3. Drag `PiggyBank.app` to Applications folder

---

### 3. **Direct Installation** (For Testing)
Install directly to your Applications:
```bash
make install
```
- ✅ Quick local installation
- 📍 Installs to: `~/Applications/PiggyBank.app`

---

## Quick Start
```bash
# Build and create PKG installer
make pkg

# Or for DMG distribution
make dmg
```

Both are user-only installations - no system-wide changes, fast and clean!
