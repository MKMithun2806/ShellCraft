#!/usr/bin/env bash
set -e

echo "[*] Installing ShellCraft..."

# Check Go exists
if ! command -v go >/dev/null 2>&1; then
    echo "[!] Go not found."
    echo "Install Go first: https://go.dev/dl/"
    exit 1
fi

# Install from GitHub
go install github.com/MKMithun2806/ShellCraft@latest

# Add Go bin to PATH if missing
if ! echo "$PATH" | grep -q "$HOME/go/bin"; then
    echo 'export PATH="$PATH:$HOME/go/bin"' >> "$HOME/.bashrc"
    export PATH="$PATH:$HOME/go/bin"
fi

echo "[+] ShellCraft installed!"
echo "[+] Try: shellcraft"
