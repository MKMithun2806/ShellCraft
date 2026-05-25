#!/bin/bash

GOBIN=$(go env GOBIN)
if [ -z "$GOBIN" ]; then
    GOBIN=$(go env GOPATH)/bin
fi

if [ -f "$GOBIN/shellcraft" ]; then
    rm "$GOBIN/shellcraft"
    echo "[+] Removed ShellCraft binary from $GOBIN"
else
    echo "[!] ShellCraft binary not found in $GOBIN"
fi

CONFIG_DIR="$HOME/.shellcraft"
if [ -d "$CONFIG_DIR" ]; then
    rm -rf "$CONFIG_DIR"
    echo "[+] Removed ShellCraft config directory ($CONFIG_DIR)"
    echo "    Includes: templates, payload history, custom payloads"
fi
