package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type ListenerType struct {
	Name    string
	Command string
	Detect  func() bool
}

func GetListeners() []ListenerType {
	return []ListenerType{
		{
			Name:    "netcat (nc)",
			Command: "nc -lvnp %s",
			Detect:  func() bool { return exec.Command("nc", "-h").Run() == nil },
		},
		{
			Name:    "socat",
			Command: "socat -d -d TCP-LISTEN:%s STDOUT",
			Detect:  func() bool { _, err := exec.LookPath("socat"); return err == nil },
		},
		{
			Name:    "rlwrap + nc (interactive TTY)",
			Command: "rlwrap nc -lvnp %s",
			Detect:  func() bool { _, err := exec.LookPath("rlwrap"); return err == nil },
		},
		{
			Name:    "ncat (Nmap)",
			Command: "ncat -lvnp %s",
			Detect:  func() bool { _, err := exec.LookPath("ncat"); return err == nil },
		},
	}
}

func handleListen(args []string) {
	if len(args) < 1 {
		fmt.Println("[!] Usage: shellcraft listen <port> [listener-type]")
		fmt.Println("    Listener types: nc, socat, rlwrap, ncat")
		return
	}

	port := args[0]
	listenerName := "nc"
	if len(args) > 1 {
		listenerName = strings.ToLower(args[1])
	}

	listeners := GetListeners()
	var selected *ListenerType
	for _, l := range listeners {
		if strings.Contains(strings.ToLower(l.Name), listenerName) {
			selected = &l
			break
		}
	}

	if selected == nil {
		fmt.Printf("[!] Unknown listener type '%s'. Available: nc, socat, rlwrap, ncat\n", listenerName)
		return
	}

	cmdStr := fmt.Sprintf(selected.Command, port)
	fmt.Printf("[*] Starting %s listener on port %s...\n", selected.Name, port)
	fmt.Printf("[*] Command: %s\n", cmdStr)
	fmt.Println(strings.Repeat("-", 50))

	parts := strings.Fields(cmdStr)
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		fmt.Printf("[!] Error running %s: %v\n", selected.Name, err)
	}
}

func SuggestListeners(port int) string {
	portStr := fmt.Sprintf("%d", port)
	var suggestions []string
	for _, l := range GetListeners() {
		if l.Detect() {
			suggestions = append(suggestions, fmt.Sprintf("  %s", fmt.Sprintf(l.Command, portStr)))
		}
	}
	suggestions = append(suggestions, fmt.Sprintf("  nc -lvnp %s", portStr))
	return strings.Join(suggestions, "\n")
}
