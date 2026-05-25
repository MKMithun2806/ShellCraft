package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

var (
	inputScanner = bufio.NewScanner(os.Stdin)
)

// GetInput prompts the user for a value and validates it.
func GetInput(prompt string) string {
	for {
		fmt.Printf("[?] %s: ", prompt)
		if !inputScanner.Scan() {
			os.Exit(0)
		}
		input := strings.TrimSpace(inputScanner.Text())
		if input != "" {
			return input
		}
	}
}

// GetValidatedIP prompts for an IP and validates it.
func GetValidatedIP(prompt string) string {
	for {
		ipStr := GetInput(prompt)
		if net.ParseIP(ipStr) != nil {
			return ipStr
		}
		fmt.Println("[!] Invalid IP address. Please enter a valid IPv4 or IPv6 address.")
	}
}

// GetValidatedPort prompts for a port and validates it.
func GetValidatedPort(prompt string) int {
	for {
		portStr := GetInput(prompt)
		port, err := strconv.Atoi(portStr)
		if err == nil && port > 0 && port <= 65535 {
			return port
		}
		fmt.Println("[!] Invalid port. Please enter a number between 1 and 65535.")
	}
}

// SelectOption prompts the user to select from a list of options.
func SelectOption(prompt string, options []string) int {
	for {
		fmt.Printf("\n--- %s ---\n", prompt)
		for i, opt := range options {
			fmt.Printf("[%d] %s\n", i+1, opt)
		}

		choiceStr := GetInput("Select an option")
		choice, err := strconv.Atoi(choiceStr)
		if err == nil && choice >= 1 && choice <= len(options) {
			return choice - 1
		}
		fmt.Printf("[!] Invalid choice. Please enter a number between 1 and %d.\n", len(options))
	}
}
