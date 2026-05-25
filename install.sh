#!/usr/bin/env bash
set -e

echo "[*] Installing ShellCraft..."

# Check Go exists
if ! command -v go >/dev/null 2>&1; then
    echo "[!] Go not found."
    echo "Install Go first: https://go.dev/dl/"
    exit 1
fi

# Install from GitHub - Now works cleanly at the root
echo "[*] Running: go install github.com/MKMithun2806/ShellCraft@latest"
go install github.com/MKMithun2806/ShellCraft@latest

# Add Go bin to PATH if missing
GOBIN=$(go env GOBIN)
if [ -z "$GOBIN" ]; then
    GOBIN=$(go env GOPATH)/bin
fi

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
