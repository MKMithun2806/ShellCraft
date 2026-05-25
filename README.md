# ShellCraft

[![Go Version](https://img.shields.io/github/go-mod/go-version/MKMithun2806/ShellCraft)](https://go.dev/)
[![License](https://img.shields.io/github/license/MKMithun2806/ShellCraft)](LICENSE)
[![Latest Release](https://img.shields.io/github/v/release/MKMithun2806/ShellCraft)](https://github.com/MKMithun2806/ShellCraft/releases)

ShellCraft is a standalone, interactive CLI tool written in Go for generating local reverse shell payloads. It is designed for cybersecurity professionals and enthusiasts to quickly generate common and obfuscated payloads using only the Go standard library.

## Features

- **Multi-Mode Generation:** Interactive guided workflow or quick CLI flags for one-liner generation.
- **Auto IP Detection:** Use `-i auto` to auto-detect your local IP (tun0, eth0, wlan0).
- **Expanded Payload Library:**
  - **Linux:** Bash, NC FIFO, Python, PHP, Ruby, Perl.
  - **Windows:** PowerShell (AMSI Bypass), CMD PowerShell one-liner, MSHTA (VBScript), Certutil Stager.
  - **macOS:** Zsh Native.
- **Obfuscation Engine (Levels 1-3):** Variable randomization, tick fragmentation, base64 encoding for PowerShell and Python payloads.
- **Encoding & Wrappers:** Raw, URL Encoded, and Base64/Pipeline.
- **Better Listener:** Start a listener with nc, socat, rlwrap, or ncat — auto-detects available tools.
- **Listener Commands:** After generating a payload, get the matching listener command for any available listener tool.
- **Template System:** Save and load payload configurations to speed up workflow.
- **Custom Payloads:** Add your own payloads via CLI or interactively — persisted to disk.
- **Payload History:** Last 20 generated payloads automatically saved with timestamps; view via `history` subcommand.
- **AV/EDR Bypass Suggestions:** Per-payload-type evasion tips (AMSI bypass, PowerShell logging bypass, PyArmor, bash string evasion).
- **Shell Upgrade Commands:** PTY spawn, script, socat, and stty magic suggestions for upgrading a basic shell.
- **HTTP Delivery Methods:** Generate curl, wget, PowerShell IEX, certutil, and base64 delivery commands.
- **C2 Framework Integration:** One-click stagers for Covenant, Sliver, Empire, and Metasploit.
- **Clipboard Support:** Copy generated payloads directly to your clipboard (xclip/xsel/pbcopy/clip).
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
shellcraft -i auto -p 4444 -t powershell -e b64 --obfs 2 --suggest
```

**Generation Flags:**
- `-i`: Attacker IP (use `auto` for auto-detection).
- `-p`: Attacker Port.
- `-t`: Payload Type (bash, nc, powershell, zsh, python, php, ruby, perl).
- `-e`: Encoding (raw, url, b64).
- `-obfs`: Obfuscation level 0-3 (PowerShell & Python only).

**Template Management Flags:**
- `-list`: List all saved templates.
- `-load <name>`: Load and run a saved template.
- `-save <name>`: Save the current CLI flags as a template.

**Extra Output Flags:**
- `-suggest`: Show AV/EDR bypass tips and shell upgrade commands.
- `-delivery`: Show HTTP delivery methods (curl, wget, IEX, etc.).
- `-c2`: Show C2 framework integration stagers.

**Custom Payload Flags:**
- `-add-custom name:code[:type:os]`: Add a custom payload from CLI.
- `-list-custom`: List all custom payloads.
- `-delete-custom <index>`: Delete a custom payload by index.

### Subcommands

**Listener:**
```bash
shellcraft listen 4444         # netcat (auto-detected)
shellcraft listen 4444 socat   # socat
shellcraft listen 4444 rlwrap  # rlwrap + nc (interactive TTY)
shellcraft listen 4444 ncat    # ncat (Nmap)
```

**Payload History:**
```bash
shellcraft history
```
Displays the last 20 generated payloads with timestamps.

**Custom Payloads:**
```bash
shellcraft custom-payload list            # list all custom payloads
shellcraft custom-payload add             # interactive add
shellcraft custom-payload delete <index>  # delete by index
```

**Version:**
```bash
shellcraft version
```

### Interactive Mode Features

- **Template Loading:** If saved templates exist, you'll be prompted to load one on startup.
- **Custom Payloads:** Custom payloads are merged into the payload selection list.
- **Obfuscation:** After selecting a PowerShell or Python payload, choose an obfuscation level.
- **Post-Generation Menu:** After a payload is generated, choose from:
  - Copy to clipboard
  - Save as template
  - Show AV/EDR bypass & shell upgrade tips
  - Show HTTP delivery methods
  - Show C2 framework integration
  - Add a custom payload

## Uninstallation

To remove ShellCraft:
```bash
curl -sSL https://raw.githubusercontent.com/MKMithun2806/ShellCraft/main/uninstall.sh | bash
```

## Disclaimer

This tool is intended for legal, authorized security testing and educational purposes only. Do not use it for malicious activities.
