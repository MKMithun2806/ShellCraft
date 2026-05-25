# ShellCraft

ShellCraft is a standalone, interactive CLI tool written in Go for generating local reverse shell payloads. It is designed for cybersecurity professionals and enthusiasts to quickly generate common and obfuscated payloads using only the Go standard library.

## Features

- **Interactive Workflow:** Guides you through inputting the attacker IP and port.
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

1. Clone the repository or download the source code.
2. Build the binary:
   ```bash
   go build -o shellcraft src/*.go
   ```

## Usage

Run the generated binary:
```bash
./shellcraft
```

Follow the interactive prompts:
1. Enter your **Attacker IP**.
2. Enter your **Attacker Port**.
3. Select the **Payload Type**.
4. Select the **Encoding/Wrapper**.

The tool will then output the `nc` listener command and the finalized payload.

## Disclaimer

This tool is intended for legal, authorized security testing and educational purposes only. Do not use it for malicious activities.
