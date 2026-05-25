package main

import (
	"fmt"
	"strings"
)

func handleCustomPayload(args []string) {
	if len(args) < 1 {
		fmt.Printf("%s\n", errorf("Usage: shellcraft custom-payload <subcommand>"))
		fmt.Println("    Subcommands: list, add, delete <index>")
		return
	}
	switch args[0] {
	case "list":
		ListCustomPayloads()
	case "add":
		AddCustomPayloadInteractive()
	case "delete":
		if len(args) < 2 {
			fmt.Printf("%s\n", errorf("Usage: shellcraft custom-payload delete <index>"))
			return
		}
		idx := 0
		fmt.Sscanf(args[1], "%d", &idx)
		if err := DeleteCustomPayload(idx - 1); err != nil {
			fmt.Printf("%s\n", errorf("Error: %v", err))
		} else {
			fmt.Printf("%s\n", successf("Custom payload deleted."))
		}
	default:
		fmt.Printf("%s\n", errorf("Unknown subcommand '%s'. Use: list, add, delete", args[0]))
	}
}

func handleHistory() {
	history, err := LoadHistory()
	if err != nil {
		fmt.Printf("%s\n", errorf("Error loading history: %v", err))
		return
	}
	if len(history) == 0 {
		fmt.Printf("%s\n", infof("No history found."))
		return
	}
	fmt.Printf("\n%s\n", headerf("Payload History (Last 20)"))
	for i, entry := range history {
		ts := entry.Timestamp
		if len(ts) > 19 {
			ts = ts[:19]
		}
		fmt.Printf("[%d] %s | %s:%d | %s | %s\n", i+1, ts, entry.IP, entry.Port, entry.Type, entry.Encoder)
		fmt.Printf("    %s\n\n", col(entry.Payload, Cyan))
	}
}

var (
	showSuggest bool
	showDelivery bool
	showC2 bool
)

func handleNonInteractive(ip string, port int, pType string, enc string, obfs int) {
	payloads := MergeCustomPayloads(GeneratePayloads(ip, port))
	var selectedPayload *Payload
	
	pType = strings.ToLower(pType)
	for _, p := range payloads {
		if strings.Contains(strings.ToLower(p.Name), pType) {
			selectedPayload = &p
			break
		}
	}

	if selectedPayload == nil {
		fmt.Printf("%s\n", errorf("Unknown payload type: %s", pType))
		return
	}

	// Apply Obfuscation
	code := Obfuscate(selectedPayload.Code, obfs, selectedPayload.Name)

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

	finalPayload := selectedEncoder.Wrap(code)
	fmt.Println(finalPayload)

	// Save to History
	SaveToHistory(HistoryEntry{
		IP:      ip,
		Port:    port,
		Type:    selectedPayload.Name,
		Encoder: selectedEncoder.Name,
		Payload: finalPayload,
	})

	if showSuggest || showDelivery || showC2 {
		if showSuggest {
			PrintSuggestions(selectedPayload.Name, port)
		}
		if showDelivery {
			PrintDeliveryMethods(ip, port, finalPayload)
		}
		if showC2 {
			PrintC2Integrations(ip, port)
		}
	}
}
