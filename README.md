# ShellCraft

[![Go Version](https://img.shields.io/github/go-mod/go-version/MKMithun2806/ShellCraft)](https://go.dev/)
[![License](https://img.shields.io/github/license/MKMithun2806/ShellCraft)](LICENSE)
[![Latest Release](https://img.shields.io/github/v/release/MKMithun2806/ShellCraft)](https://github.com/MKMithun2806/ShellCraft/releases)

ShellCraft is a standalone, interactive CLI tool written in Go for generating local reverse shell payloads. It is designed for cybersecurity professionals and enthusiasts to quickly generate common and obfuscated payloads using only the Go standard library.

## Features

- **Multi-Mode Generation:** Use the interactive guided workflow or quick CLI flags for one-liner generation.
- **Listener Helper:** Quickly start a netcat listener with the `listen` subcommand.
- **Template System:** Save and load custom payload configurations to speed up your workflow.
- **Clipboard Support:** Copy generated payloads directly to your clipboard with one click.
- **Expanded Payload Library:**
  - **Linux:** Bash, NC FIFO, Python, PHP, Ruby, Perl.
  - **Windows:** PowerShell (AMSI Bypass).
  - **macOS:** Zsh Native.
- **Encoding & Wrappers:** Raw, URL Encoded, and Base64/Pipeline.
- **Zero External Dependencies:** Built entirely with the Go standard library.
- **Clean Exit:** Gracefully handles `Ctrl+C`.

## Installation

Ensure you have Go installed on your system.

### Option 1: Pre-built Binaries (Recommended)

You can download the latest pre-built binaries for Linux, Windows, and macOS directly from the [Releases](https://github.com/MKMithun2806/ShellCraft/releases) page.

1. Download the binary for your platform.
2. Make it executable (Linux/macOS): `chmod +x shellcraft-*`
3. Run it: `./shellcraft-*`

### Option 2: Quick Install (via script)

Run the following command to download and run the installer:
```bash
curl -sSL https://raw.githubusercontent.com/MKMithun2806/ShellCraft/main/install.sh | bash
```

### Option 3: Using go install

```bash
go install github.com/MKMithun2806/ShellCraft/cmd/shellcraft@latest
```

### Option 4: Manual Build

1. Clone the repository:
   ```bash
   git clone https://github.com/MKMithun2806/ShellCraft.git
   cd ShellCraft
   ```
2. Build the binary:
   ```bash
   go build -ldflags "-X main.Version=1.0.0" -o shellcraft ./cmd/shellcraft
   ```

## Usage

### Interactive Mode
Simply run the tool and follow the prompts:
```bash
shellcraft
```

### One-liner Mode (Non-interactive)
Generate payloads instantly using CLI flags:
```bash
shellcraft -i 10.10.10.10 -p 4444 -t python -e base64
```

**Generation Flags:**
- `-i`: Attacker IP
- `-p`: Attacker Port
- `-t`: Payload Type (bash, nc, ps, zsh, python, php, ruby, perl)
- `-e`: Encoding (raw, url, b64)

**Template Management Flags:**
- `-list`: List all saved templates.
- `-load <name>`: Load and run a saved template.
- `-save <name>`: Save the current CLI flags as a template.

### Listener Helper
Quickly start a netcat listener:
```bash
shellcraft listen 4444
```

### Template System
After generating a payload in interactive mode, select **"Save as Template"**. The tool will prompt you to load saved templates the next time you run it interactively.

## Uninstallation

To remove ShellCraft:
```bash
curl -sSL https://raw.githubusercontent.com/MKMithun2806/ShellCraft/main/uninstall.sh | bash
```

## Disclaimer

This tool is intended for legal, authorized security testing and educational purposes only. Do not use it for malicious activities.
