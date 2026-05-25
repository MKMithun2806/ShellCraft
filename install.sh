#!/usr/bin/env bash
set -e

echo "[*] Installing ShellCraft..."

# Check Go exists
if ! command -v go >/dev/null 2>&1; then
    echo "[!] Go not found."
    echo "Install Go first: https://go.dev/dl/"
    exit 1
fi

# Clean install from GitHub using standard Go layout
echo "[*] Running: go install github.com/MKMithun2806/ShellCraft/cmd/shellcraft@latest"
go install github.com/MKMithun2806/ShellCraft/cmd/shellcraft@latest

# Determine GOBIN
GOBIN=$(go env GOBIN)
if [ -z "$GOBIN" ]; then
    GOPATH=$(go env GOPATH)
    GOBIN="$GOPATH/bin"
fi

# Add Go bin to PATH if missing
if ! echo "$PATH" | grep -q "$GOBIN"; then
    SHELL_RC="$HOME/.bashrc"
    # Detect Zsh
    if [[ "$SHELL" == */zsh ]] || [ -f "$HOME/.zshrc" ]; then
        SHELL_RC="$HOME/.zshrc"
    fi
    
    echo "" >> "$SHELL_RC"
    echo "# Go Binaries" >> "$SHELL_RC"
    echo "export PATH=\"\$PATH:$GOBIN\"" >> "$SHELL_RC"
    export PATH="$PATH:$GOBIN"
    echo "[*] Added $GOBIN to $SHELL_RC"
    echo "[*] Please run: source $SHELL_RC"
fi

if command -v shellcraft >/dev/null 2>&1; then
    echo "[+] ShellCraft installed successfully!"
    echo "[+] Try: shellcraft --version"
else
    echo "[!] Installation finished, but 'shellcraft' command not found in PATH."
    echo "[!] You may need to restart your terminal or source your config."
fi
