<div align="center">

# ShellCraft

**Reverse shell payload generator — interactive, obfuscated, ready to go.**

[![Go Version](https://img.shields.io/github/go-mod/go-version/MKMithun2806/ShellCraft?style=flat-square&logo=go)](https://go.dev/)
[![License](https://img.shields.io/github/license/MKMithun2806/ShellCraft?style=flat-square)](/LICENSE)
[![Latest Release](https://img.shields.io/github/v/release/MKMithun2806/ShellCraft?style=flat-square&logo=github)](https://github.com/MKMithun2806/ShellCraft/releases)
[![Build](https://img.shields.io/github/actions/workflow/status/MKMithun2806/ShellCraft/release.yml?style=flat-square)](https://github.com/MKMithun2806/ShellCraft/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/MKMithun2806/ShellCraft?style=flat-square)](https://goreportcard.com/report/github.com/MKMithun2806/ShellCraft)

</div>

---

## ✨ Features

<table>
<tr>
<td width="50%">

**Generation**
- Interactive guided workflow & CLI one-liners
- Auto IP detection (`-i auto`)
- 10+ payload types across Linux, Windows, macOS
- Obfuscation engine (Levels 1–3) for PowerShell & Python
- Encoding wrappers: Raw, URL, Base64

**Post-Generation**
- AV/EDR bypass suggestions per payload type
- Shell upgrade commands (PTY, socat, stty)
- HTTP delivery methods (curl, wget, IEX, certutil)
- C2 framework stagers (Covenant, Sliver, Empire, Metasploit)
- Clipboard copy (xclip/xsel/pbcopy/clip)

</td>
<td width="50%">

**🎛️ Listeners**
- **Built-in Go listener** (default, no external tools required)
  - PTY-aware raw terminal via `golang.org/x/term`
  - Arrow keys, Tab, Ctrl+C, Ctrl+Z pass through transparently
  - SIGWINCH resize handling
  - Multi-connection management with interactive switching menu
  - `Ctrl+]` to suspend session and return to connection menu
  - A proper `nc -lvnp` / ncat replacement — just `shellcraft listen 4444`
- External tools still fully supported: `nc`, `socat`, `rlwrap`, `ncat`
- Auto-detects available tools, colorized suggestions

**💾 Persistence**
- Payload history (last 20, with timestamps)
- Template system (save/load configs)
- Custom user-defined payloads (persisted to disk)

**🔧 Internals**
- Zero external dependencies — pure Go stdlib
- Clean `Ctrl+C` handling
- Cross-platform (Linux, macOS, Windows)

</td>
</tr>
</table>

---

## 📦 Installation

<details open>
<summary><b>Option 1 — Pre-built binaries</b> ⭐ <i>(recommended)</i></summary>

Download from the [Releases page](https://github.com/MKMithun2806/ShellCraft/releases), then:

```bash
chmod +x shellcraft-*
./shellcraft-*
```

</details>

<details>
<summary><b>Option 2 — Quick install script</b></summary>

```bash
curl -sSL https://raw.githubusercontent.com/MKMithun2806/ShellCraft/main/install.sh | bash
```

</details>

<details>
<summary><b>Option 3 — <code>go install</code></b></summary>

```bash
go install github.com/MKMithun2806/ShellCraft/cmd/shellcraft@main
```

</details>

<details>
<summary><b>Option 4 — Manual build</b></summary>

```bash
git clone https://github.com/MKMithun2806/ShellCraft.git
cd ShellCraft
go build -ldflags "-X main.Version=1.0.0" -o shellcraft ./cmd/shellcraft
```

</details>

---

## 🎮 Usage

### Interactive mode

```bash
shellcraft
```

Follow the prompts to configure IP, port, payload type, encoding, and obfuscation.

### One-liner mode

```bash
# Basic
shellcraft -i 10.10.10.10 -p 4444 -t python -e base64

# With auto IP and obfuscation
shellcraft -i auto -p 4444 -t powershell -e b64 --obfs 2 --suggest
```

<details>
<summary><b>📋 Full CLI reference</b></summary>

#### Generation flags

| Flag | Description | Example |
|------|-------------|---------|
| `-i` | Attacker IP (`auto` for auto-detection) | `-i 10.0.0.1` |
| `-p` | Attacker port | `-p 4444` |
| `-t` | Payload type | `-t powershell` |
| `-e` | Encoding: `raw`, `url`, `b64` | `-e b64` |
| `-obfs` | Obfuscation level 0–3 (PS & Python) | `-obfs 2` |

#### Template flags

| Flag | Description |
|------|-------------|
| `-list` | List saved templates |
| `-load <name>` | Load and run a template |
| `-save <name>` | Save current config as template |

#### Extra output flags

| Flag | Description |
|------|-------------|
| `-suggest` | Show AV/EDR bypass & shell upgrade tips |
| `-delivery` | Show HTTP delivery methods |
| `-c2` | Show C2 framework stagers |

#### Custom payload flags

| Flag | Description |
|------|-------------|
| `-add-custom name:code[:type:os]` | Add a custom payload |
| `-list-custom` | List custom payloads |
| `-delete-custom <index>` | Delete by index |

</details>

### Subcommands

<details open>
<summary><b>🎧 Listener</b></summary>

The default listener is the **built-in Go TCP listener** — a full-featured `nc -lvnp` replacement:

```bash
# Built-in listener (default) — no external tools needed
shellcraft listen 4444
shellcraft listen -p 4444
shellcraft listener --port 9001
```

Inside an interactive session:
- **Remote output** is colored in cyan for visibility
- **Arrow keys, Tab, Ctrl+C, Ctrl+Z** work naturally (raw PTY mode)
- **Ctrl+]** suspends the session and returns to the connection menu
- Type `exit` on the remote shell to close the connection automatically

External tools are also fully supported for when you prefer them:

```bash
shellcraft listen 4444 nc
shellcraft listen 4444 socat
shellcraft listen 4444 rlwrap
shellcraft listen 4444 ncat
```

</details>

<details>
<summary><b>Payload history</b></summary>

```bash
shellcraft history
```

Shows the last 20 generated payloads with timestamps, IP, port, type, and encoder.

</details>

<details>
<summary><b>Custom payloads</b></summary>

```bash
shellcraft custom-payload list
shellcraft custom-payload add            # interactive
shellcraft custom-payload delete <index>
```

Custom payloads are automatically merged into the interactive payload selection menu.

</details>

<details>
<summary><b>ℹ️ Version</b></summary>

```bash
shellcraft version
```

</details>

### Post-generation menu

After generating a payload interactively, the menu lets you:

| # | Action |
|---|--------|
| 1 | Copy payload to clipboard |
| 2 | Save as template |
| 3 | Show AV/EDR bypass & shell upgrade tips |
| 4 | Show HTTP delivery methods |
| 5 | Show C2 framework integration |
| 6 | Add a custom payload |
| 7 | Start listener (built-in Go listener, defaults to PTY raw mode) |
| 8 | Exit |

---

## 🧪 Payload reference

| OS | Payload | Notes |
|----|---------|-------|
| 🐧 Linux | Bash, NC FIFO, Python, PHP, Ruby, Perl | — |
| 🪟 Windows | PowerShell (AMSI Bypass), CMD, MSHTA (VBScript), Certutil Stager | Includes signature evasion |
| 🍏 macOS | Zsh Native | Native `ztcp` |

### Obfuscation levels

| Level | PowerShell | Python |
|-------|------------|--------|
| 0 | None | None |
| 1 | Tick fragmentation keywords | Base64 exec wrapper |
| 2 | Variable name randomization | Base64 exec wrapper |
| 3 | Variable randomization + string concatenation | Base64 exec wrapper |

---

## 🗑️ Uninstalll Script (Removes EVERYTHING shellcraft)

```bash
curl -sSL https://raw.githubusercontent.com/MKMithun2806/ShellCraft/main/uninstall.sh | bash
```

---

## ⚠️ Disclaimer

> This tool is intended for **legal, authorized security testing and educational purposes only**.  
> Do not use it for malicious activities. The authors are not responsible for misuse.

---

<div align="center">
Made with Go &nbsp;·&nbsp; No external dependencies
</div>
