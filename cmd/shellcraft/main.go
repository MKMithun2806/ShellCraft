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
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "listen":
			handleListen(os.Args[2:])
			return
		case "version", "--version", "-v":
			fmt.Printf("ShellCraft v%s\n", Version)
			return
		}
	}

	ipFlag := flag.String("i", "", "Attacker IP")
	portFlag := flag.Int("p", 0, "Attacker Port")
	typeFlag := flag.String("t", "", "Payload Type (bash, nc, ps, zsh, python, php, ruby, perl)")
	encFlag := flag.String("e", "raw", "Encoding (raw, url, b64)")
	
	listFlag := flag.Bool("list", false, "List all saved templates")
	loadFlag := flag.String("load", "", "Load and run a saved template")
	saveFlag := flag.String("save", "", "Save current CLI config as a template")
	flag.Parse()

	// 1. Handle Template Listing
	if *listFlag {
		templates, _ := ListTemplates()
		if len(templates) == 0 {
			fmt.Println("[!] No templates found.")
			return
		}
		fmt.Println("\n--- Saved Templates ---")
		for _, t := range templates {
			fmt.Printf("- %-15s (%s:%d | %s | %s)\n", t.Name, t.IP, t.Port, t.Type, t.Encoder)
		}
		return
	}

	// 2. Handle Template Loading (Non-interactive)
	if *loadFlag != "" {
		templates, _ := ListTemplates()
		var found *Template
		for _, t := range templates {
			if strings.EqualFold(t.Name, *loadFlag) {
				found = &t
				break
			}
		}
		if found == nil {
			fmt.Printf("[!] Template '%s' not found.\n", *loadFlag)
			return
		}
		handleNonInteractive(found.IP, found.Port, found.Type, found.Encoder)
		return
	}

	// 3. Handle Non-interactive Mode (Direct Flags)
	if *ipFlag != "" && *portFlag != 0 && *typeFlag != "" {
		if *saveFlag != "" {
			err := SaveTemplate(Template{
				Name:    *saveFlag,
				IP:      *ipFlag,
				Port:    *portFlag,
				Type:    *typeFlag,
				Encoder: *encFlag,
			})
			if err != nil {
				fmt.Printf("[!] Error saving template: %v\n", err)
			} else {
				fmt.Printf("[+] Template '%s' saved!\n", *saveFlag)
			}
		}
		handleNonInteractive(*ipFlag, *portFlag, *typeFlag, *encFlag)
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

	var ip string
	var port int
	var selectedPayload Payload
	var selectedEncoder Encoder

	// Template Selection
	templates, _ := ListTemplates()
	if len(templates) > 0 {
		fmt.Println("[*] Saved Templates found.")
		options := []string{"Use a Template", "Continue without Template"}
		if SelectOption("Template Menu", options) == 0 {
			tNames := make([]string, len(templates))
			for i, t := range templates {
				tNames[i] = t.Name
			}
			tIdx := SelectOption("Select a Template", tNames)
			t := templates[tIdx]
			ip = t.IP
			port = t.Port
			
			// Auto-select payload
			payloads := GeneratePayloads(ip, port)
			for _, p := range payloads {
				if strings.Contains(strings.ToLower(p.Name), strings.ToLower(t.Type)) {
					selectedPayload = p
					break
				}
			}
			
			// Auto-select encoder
			encoders := GetEncoders()
			for _, e := range encoders {
				if strings.Contains(strings.ToLower(e.Name), strings.ToLower(t.Encoder)) {
					selectedEncoder = e
					break
				}
			}
			fmt.Printf("[+] Loaded Template: %s\n", t.Name)
		}
	}

	if ip == "" {
		// 1. Get Attacker IP
		ip = GetValidatedIP("Attacker IP")

		// 2. Get Attacker Port
		port = GetValidatedPort("Attacker Port")

		// 3. Select Payload Type
		payloads := GeneratePayloads(ip, port)
		payloadNames := make([]string, len(payloads))
		for i, p := range payloads {
			payloadNames[i] = p.Name
		}
		payloadIdx := SelectOption("Select Payload Type", payloadNames)
		selectedPayload = payloads[payloadIdx]

		// 4. Select Encoding/Wrapper
		encoders := GetEncoders()
		encoderNames := make([]string, len(encoders))
		for i, e := range encoders {
			encoderNames[i] = e.Name
		}
		encoderIdx := SelectOption("Select Encoding/Wrapper", encoderNames)
		selectedEncoder = encoders[encoderIdx]
	}

	// 5. Generate and Print
	finalPayload := selectedEncoder.Wrap(selectedPayload.Code)

	fmt.Printf("\n" + strings.Repeat("=", 60))
	fmt.Printf("\n[+] Listener Command:\n    nc -lvnp %d\n", port)
	fmt.Printf("\n[+] Finalized Payload (%s - %s):\n\n%s\n", selectedPayload.Name, selectedEncoder.Name, finalPayload)
	fmt.Printf("\n" + strings.Repeat("=", 60) + "\n")

	// Post-Generation Actions
	options := []string{"Copy Payload to Clipboard", "Save as Template", "Exit"}
	choice := SelectOption("Post-Generation Actions", options)

	switch choice {
	case 0:
		if err := copyToClipboard(finalPayload); err != nil {
			fmt.Printf("[!] Failed to copy to clipboard: %v\n", err)
		} else {
			fmt.Println("[+] Payload copied to clipboard!")
		}
	case 1:
		name := GetInput("Enter Template Name")
		err := SaveTemplate(Template{
			Name:    name,
			IP:      ip,
			Port:    port,
			Type:    selectedPayload.Name,
			Encoder: selectedEncoder.Name,
		})
		if err != nil {
			fmt.Printf("[!] Error saving template: %v\n", err)
		} else {
			fmt.Println("[+] Template saved successfully!")
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
