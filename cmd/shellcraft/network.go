package main

import (
	"fmt"
	"net"
	"strings"
)

func DetectLocalIP() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	// Priority list for interfaces
	priority := []string{"tun", "tap", "eth", "wlan", "en", "wl"}

	for _, p := range priority {
		for _, iface := range interfaces {
			if strings.HasPrefix(strings.ToLower(iface.Name), p) {
				addrs, err := iface.Addrs()
				if err != nil {
					continue
				}
				for _, addr := range addrs {
					if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
						if ipnet.IP.To4() != nil {
							return ipnet.IP.String(), nil
						}
					}
				}
			}
		}
	}

	// Fallback: search all
	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					return ipnet.IP.String(), nil
				}
			}
		}
	}

	return "", fmt.Errorf("could not detect a non-loopback IP address")
}
