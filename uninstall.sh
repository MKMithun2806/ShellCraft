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
