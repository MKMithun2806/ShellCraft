#!/usr/bin/env bash
set -e

echo "[*] Installing ShellCraft..."

# Check Go exists
if ! command -v go >/dev/null 2>&1; then
    echo "[!] Go not found."
    echo "Install Go first: https://go.dev/dl/"
    exit 1
fi

# Install from GitHub - fixed path to src
echo "[*] Running: go install github.com/MKMithun2806/ShellCraft/src@latest"
go install github.com/MKMithun2806/ShellCraft/src@latest

# The binary name will be 'src' because of the path, we should rename it to 'shellcraft' 
# or use a different install method. Actually, better to use the local build if in repo, 
# but for remote install, the package name in 'src/main.go' is 'main'.
# If 'go install' is used on a subfolder, the binary is named after the subfolder.

GOBIN=$(go env GOBIN)
if [ -z "$GOBIN" ]; then
    GOBIN=$(go env GOPATH)/bin
fi

if [ -f "$GOBIN/src" ]; then
    mv "$GOBIN/src" "$GOBIN/shellcraft"
    echo "[+] Renamed binary to shellcraft"
fi

# Add Go bin to PATH if missing
if ! echo "$PATH" | grep -q "$GOBIN"; then
    SHELL_RC="$HOME/.bashrc"
    if [[ "$SHELL" == */zsh ]]; then
        SHELL_RC="$HOME/.zshrc"
    fi
    echo "export PATH=\"\$PATH:$GOBIN\"" >> "$SHELL_RC"
    export PATH="$PATH:$GOBIN"
    echo "[*] Added $GOBIN to $SHELL_RC"
fi

echo "[+] ShellCraft installed!"
echo "[+] Try: shellcraft"
