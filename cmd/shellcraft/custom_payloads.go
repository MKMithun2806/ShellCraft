package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type CustomPayloadEntry struct {
	Name    string `json:"name"`
	Code    string `json:"code"`
	Type    string `json:"type"`
	OSType  string `json:"os_type"`
}

func getCustomPayloadsFile() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot get home dir: %w", err)
	}
	dir := filepath.Join(home, ".shellcraft")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("cannot create config dir: %w", err)
	}
	return filepath.Join(dir, "custom_payloads.json"), nil
}

func LoadCustomPayloads() ([]CustomPayloadEntry, error) {
	path, err := getCustomPayloadsFile()
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return []CustomPayloadEntry{}, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var payloads []CustomPayloadEntry
	if err := json.Unmarshal(data, &payloads); err != nil {
		return nil, err
	}
	return payloads, nil
}

func SaveCustomPayloads(payloads []CustomPayloadEntry) error {
	path, err := getCustomPayloadsFile()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(payloads, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func AddCustomPayloadInteractive() {
	payloads, _ := LoadCustomPayloads()

	name := GetInput("Custom Payload Name")
	code := GetInput("Payload Command/Code")
	pType := GetInput("Type hint (e.g., bash, python, powershell)")
	osType := GetInput("OS (linux/windows/macos)")

	payloads = append(payloads, CustomPayloadEntry{
		Name:   name,
		Code:   code,
		Type:   pType,
		OSType: osType,
	})

	if err := SaveCustomPayloads(payloads); err != nil {
		fmt.Printf("%s\n", errorf("Error saving custom payload: %v", err))
	} else {
		fmt.Printf("%s\n", successf("Custom payload '%s' saved!", name))
	}
}

func AddCustomPayload(name, code, pType, osType string) error {
	payloads, _ := LoadCustomPayloads()
	payloads = append(payloads, CustomPayloadEntry{
		Name:   name,
		Code:   code,
		Type:   pType,
		OSType: osType,
	})
	return SaveCustomPayloads(payloads)
}

func MergeCustomPayloads(base []Payload) []Payload {
	custom, err := LoadCustomPayloads()
	if err != nil || len(custom) == 0 {
		return base
	}
	for _, c := range custom {
		base = append(base, Payload{
			Name: fmt.Sprintf("Custom: %s (%s)", c.Name, c.OSType),
			Code: c.Code,
		})
	}
	return base
}

func ListCustomPayloads() {
	payloads, err := LoadCustomPayloads()
	if err != nil {
		fmt.Printf("%s\n", errorf("Error loading custom payloads: %v", err))
		return
	}
	if len(payloads) == 0 {
		fmt.Printf("%s\n", errorf("No custom payloads found."))
		return
	}
	fmt.Printf("\n%s\n", headerf("Custom Payloads"))
	for i, p := range payloads {
		fmt.Printf("[%d] %-20s | Type: %s | OS: %s\n", i+1, p.Name, p.Type, p.OSType)
	}
}

func DeleteCustomPayload(idx int) error {
	payloads, err := LoadCustomPayloads()
	if err != nil {
		return err
	}
	if idx < 0 || idx >= len(payloads) {
		return fmt.Errorf("invalid index %d", idx)
	}
	payloads = append(payloads[:idx], payloads[idx+1:]...)
	return SaveCustomPayloads(payloads)
}
