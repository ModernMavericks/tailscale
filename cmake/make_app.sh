#!/bin/sh
# make_app.sh <systray-binary> <out .app> <version>
# Wrap the tailscale-systray Go binary in a menu-bar-only (LSUIElement) .app. No ObjC -- the Go binary
# IS the app; the bundle just gives it an Info.plist so LaunchServices treats it as a menu-bar agent.
set -eu
BIN=$1; APP=$2; VER=$3
rm -rf "$APP"
mkdir -p "$APP/Contents/MacOS" "$APP/Contents/Resources"
install -m 0755 "$BIN" "$APP/Contents/MacOS/tailscale-systray"
cat > "$APP/Contents/Info.plist" <<PLIST
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>CFBundleName</key><string>Tailscale</string>
  <key>CFBundleDisplayName</key><string>Tailscale</string>
  <key>CFBundleIdentifier</key><string>dev.modernmavericks.tailscale-systray</string>
  <key>CFBundleExecutable</key><string>tailscale-systray</string>
  <key>CFBundleVersion</key><string>${VER}</string>
  <key>CFBundleShortVersionString</key><string>${VER}</string>
  <key>CFBundlePackageType</key><string>APPL</string>
  <key>LSUIElement</key><true/>
  <key>LSMinimumSystemVersion</key><string>10.9</string>
</dict>
</plist>
PLIST
echo "$APP"
