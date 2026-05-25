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
		case "history":
			handleHistory()
			return
		case "custom-payload":
			handleCustomPayload(os.Args[2:])
			return
		case "version", "--version", "-v":
			fmt.Printf("ShellCraft v%s\n", Version)
			return
		}
	}

	ipFlag := flag.String("i", "", "Attacker IP (use 'auto' for auto-detection)")
	portFlag := flag.Int("p", 0, "Attacker Port")
	typeFlag := flag.String("t", "", "Payload Type (bash, nc, ps, zsh, python, php, ruby, perl)")
	encFlag := flag.String("e", "raw", "Encoding (raw, url, b64)")
	obfsFlag := flag.Int("obfs", 0, "Obfuscation Level (0-3)")
	
	listFlag := flag.Bool("list", false, "List all saved templates")
	loadFlag := flag.String("load", "", "Load and run a saved template")
	saveFlag := flag.String("save", "", "Save current CLI config as a template")

	deliveryFlag := flag.Bool("delivery", false, "Show HTTP delivery methods")
	c2Flag := flag.Bool("c2", false, "Show C2 framework integration commands")
	suggestFlag := flag.Bool("suggest", false, "Show AV/EDR bypass and shell upgrade suggestions")

	addCustomFlag := flag.String("add-custom", "", "Add a custom payload (name:code:type:os)")
	listCustomFlag := flag.Bool("list-custom", false, "List all custom payloads")
	deleteCustomFlag := flag.Int("delete-custom", -1, "Delete a custom payload by index")
	flag.Parse()

	showSuggest = *suggestFlag
	showDelivery = *deliveryFlag
	showC2 = *c2Flag

	if *listCustomFlag {
		ListCustomPayloads()
		return
	}
	if *deleteCustomFlag >= 0 {
		if err := DeleteCustomPayload(*deleteCustomFlag - 1); err != nil {
			fmt.Printf("[!] Error deleting custom payload: %v\n", err)
		} else {
			fmt.Println("[+] Custom payload deleted.")
		}
		return
	}
	if *addCustomFlag != "" {
		parts := strings.SplitN(*addCustomFlag, ":", 4)
		if len(parts) < 2 {
			fmt.Println("[!] Usage: --add-custom name:code[:type:os]")
			return
		}
		name, code := parts[0], parts[1]
		pType, osType := "custom", "linux"
		if len(parts) > 2 {
			pType = parts[2]
		}
		if len(parts) > 3 {
			osType = parts[3]
		}
		if err := AddCustomPayload(name, code, pType, osType); err != nil {
			fmt.Printf("[!] Error: %v\n", err)
		} else {
			fmt.Printf("[+] Custom payload '%s' saved!\n", name)
		}
		return
	}

	if *ipFlag == "auto" {
		autoIP, err := DetectLocalIP()
		if err != nil {
			fmt.Printf("[!] IP Detection failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("[*] Auto-detected IP: %s\n", autoIP)
		*ipFlag = autoIP
	}

	// 1. Handle Template Listing
	if *listFlag {
		templates, err := ListTemplates()
		if err != nil {
			fmt.Printf("[!] Error loading templates: %v\n", err)
			return
		}
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
		templates, err := ListTemplates()
		if err != nil {
			fmt.Printf("[!] Error loading templates: %v\n", err)
			return
		}
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
		handleNonInteractive(found.IP, found.Port, found.Type, found.Encoder, *obfsFlag)
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
		handleNonInteractive(*ipFlag, *portFlag, *typeFlag, *encFlag, *obfsFlag)
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

	fmt.Print(banner)
	fmt.Printf("      Version: %s\n", Version)

	var ip string
	var port int
	var selectedPayload Payload
	var selectedEncoder Encoder

	// Template Selection
	templates, err := ListTemplates()
	if err != nil {
		fmt.Printf("[!] Warning: Could not load templates: %v\n", err)
	}
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
		ip = GetValidatedIP("Attacker IP (or 'auto')")
		if ip == "auto" {
			autoIP, _ := DetectLocalIP()
			ip = autoIP
			fmt.Printf("[*] Auto-detected IP: %s\n", ip)
		}

		// 2. Get Attacker Port
		port = GetValidatedPort("Attacker Port")

		// 3. Select Payload Type
		payloads := MergeCustomPayloads(GeneratePayloads(ip, port))
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

	// Apply Obfuscation
	payloadCode := selectedPayload.Code
	if strings.Contains(strings.ToLower(selectedPayload.Name), "powershell") || strings.Contains(strings.ToLower(selectedPayload.Name), "python") {
		options := []string{"None (Level 0)", "Basic (Level 1)", "Medium (Level 2)", "Advanced (Level 3)"}
		obfsLevel := SelectOption("Select Obfuscation Level", options)
		payloadCode = Obfuscate(payloadCode, obfsLevel, selectedPayload.Name)
	}

	// 5. Generate and Print
	finalPayload := selectedEncoder.Wrap(payloadCode)

	// Save to History
	SaveToHistory(HistoryEntry{
		IP:      ip,
		Port:    port,
		Type:    selectedPayload.Name,
		Encoder: selectedEncoder.Name,
		Payload: finalPayload,
	})

	fmt.Printf("\n" + strings.Repeat("=", 60))
	fmt.Printf("\n[+] Listener Commands:\n%s\n", SuggestListeners(port))
	fmt.Printf("\n[+] Finalized Payload (%s - %s):\n\n%s\n", selectedPayload.Name, selectedEncoder.Name, finalPayload)

	// Post-Generation Actions
	options := []string{
		"Copy Payload to Clipboard",
		"Save as Template",
		"Show AV/EDR Bypass & Shell Upgrade Tips",
		"Show HTTP Delivery Methods",
		"Show C2 Framework Integration",
		"Add Custom Payload",
		"Exit",
	}
	choice := SelectOption("Post-Generation Actions", options)

	switch choice {
	case 0:
		if err := copyToClipboard(finalPayload); err != nil {
			fmt.Printf("[!] Clipboard copy failed: %v\n", err)
			outPath := "payload.txt"
			if writeErr := os.WriteFile(outPath, []byte(finalPayload), 0644); writeErr == nil {
				fmt.Printf("[+] Payload written to %s instead.\n", outPath)
			}
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
	case 2:
		PrintSuggestions(selectedPayload.Name, port)
	case 3:
		PrintDeliveryMethods(ip, port, finalPayload)
	case 4:
		PrintC2Integrations(ip, port)
	case 5:
		AddCustomPayloadInteractive()
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
