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

func GetValidatedIP(prompt string) string {
	for {
		ipStr := GetInput(prompt)
		if ipStr == "auto" {
			return ipStr
		}
		ip := net.ParseIP(ipStr)
		if ip == nil {
			fmt.Println("[!] Invalid IP address format. Please enter a valid IPv4 or IPv6 address.")
			continue
		}
		if ip.IsUnspecified() {
			fmt.Println("[!] 0.0.0.0 is not a valid attacker IP. Use 'auto' or a real interface IP.")
			continue
		}
		if ip.IsLoopback() && ipStr != "127.0.0.1" && ipStr != "::1" {
			fmt.Println("[!] Loopback addresses (127.x.x.x) won't work for remote shells. Use a real network IP or 'auto'.")
			continue
		}
		if ip.IsMulticast() {
			fmt.Println("[!] Multicast addresses cannot be used for reverse shells.")
			continue
		}
		return ipStr
	}
}

func GetValidatedPort(prompt string) int {
	for {
		portStr := GetInput(prompt)
		port, err := strconv.Atoi(portStr)
		if err != nil {
			fmt.Println("[!] Port must be a number.")
			continue
		}
		if port < 1 || port > 65535 {
			fmt.Println("[!] Port must be between 1 and 65535.")
			continue
		}
		if port < 1024 {
			fmt.Println("[!] Ports below 1024 require root. Consider using 1024-65535.")
		}
		return port
	}
}

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
