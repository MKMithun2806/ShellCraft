package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/MKMithun2806/ShellCraft/internal/listener"
)

type ListenerType struct {
	Name    string
	Command string
	Detect  func() bool
}

func GetListeners() []ListenerType {
	return []ListenerType{
		{
			Name:    "Built-in (Go listener)",
			Command: "",
			Detect:  func() bool { return true },
		},
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
		fmt.Printf("%s\n", errorf("Usage: shellcraft listen <port> [listener-type]"))
		fmt.Println("    Listener types: built-in, nc, socat, rlwrap, ncat")
		fmt.Println("    Defaults to the built-in Go listener (no external tools needed).")
		return
	}

	port := args[0]
	p, err := strconv.Atoi(port)
	if err != nil {
		fmt.Printf("%s\n", errorf("Invalid port: %s", port))
		return
	}

	listenerName := "built-in"
	if len(args) > 1 {
		listenerName = strings.ToLower(args[1])
	}

	// Built-in listener
	if listenerName == "built-in" || listenerName == "builtin" || listenerName == "go" || listenerName == "internal" {
		fmt.Printf("%s\n", successf("Starting built-in listener on port %d", p))
		fmt.Printf("  %s\n", infof("No external tools required — pure Go TCP listener"))
		fmt.Println()
		if err := listener.Start(p); err != nil {
			fmt.Printf("%s\n", errorf("Listener error: %v", err))
		}
		return
	}

	// External tool listener
	listeners := GetListeners()
	var selected *ListenerType
	for _, l := range listeners {
		if strings.Contains(strings.ToLower(l.Name), listenerName) {
			selected = &l
			break
		}
	}

	if selected == nil {
		fmt.Printf("%s\n", errorf("Unknown listener type '%s'. Available: built-in, nc, socat, rlwrap, ncat", listenerName))
		return
	}

	cmdStr := fmt.Sprintf(selected.Command, port)
	fmt.Printf("%s\n", infof("Starting %s listener on port %s...", selected.Name, port))
	fmt.Printf("%s\n", col(cmdStr, Cyan))

	parts := strings.Fields(cmdStr)
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		fmt.Printf("%s\n", errorf("Error running %s: %v", selected.Name, err))
	}
}

func SuggestListeners(port int) string {
	portStr := fmt.Sprintf("%d", port)
	var suggestions []string
	suggestions = append(suggestions, col("  shellcraft listen "+portStr+"   (built-in Go listener)", Cyan))
	for _, l := range GetListeners() {
		if l.Detect() && l.Command != "" {
			suggestions = append(suggestions, col(fmt.Sprintf("  "+l.Command, portStr), Cyan))
		}
	}
	return strings.Join(suggestions, "\n")
}

func PickListenerAndStart(port int) {
	portStr := fmt.Sprintf("%d", port)
	available := GetListeners()
	var names []string
	for _, l := range available {
		if l.Detect() {
			names = append(names, l.Name)
		}
	}
	if len(names) == 0 {
		names = append(names, "nc -lvnp PORT (default fallback)")
	}
	names = append(names, "Cancel")

	idx := SelectOption("Select Listener Type", names)
	if idx >= len(names)-1 {
		fmt.Printf("%s\n", infof("Listener cancelled."))
		return
	}

	lt := available[idx]

	// Built-in listener selected
	if strings.Contains(strings.ToLower(lt.Name), "built-in") {
		fmt.Printf("%s\n", successf("Starting built-in listener on port %d", port))
		fmt.Printf("  %s\n", infof("No external tools required"))
		if err := listener.Start(port); err != nil {
			fmt.Printf("%s\n", errorf("Listener error: %v", err))
		}
		return
	}

	// External tool listener
	cmdStr := fmt.Sprintf(lt.Command, portStr)
	fmt.Printf("%s\n", col(cmdStr, Cyan))
	confirm := GetInput("Start this listener? (y/n)")
	if strings.ToLower(confirm) != "y" && strings.ToLower(confirm) != "yes" {
		fmt.Printf("%s\n", infof("Listener cancelled."))
		return
	}

	fmt.Printf("%s\n", infof("Starting %s listener on port %s...", lt.Name, portStr))
	parts := strings.Fields(cmdStr)
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		fmt.Printf("%s\n", errorf("Error running %s: %v", lt.Name, err))
	}
}
