package listener

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"golang.org/x/term"
)

// ─── ANSI colors (self-contained, mirrors cmd/shellcraft/color.go) ─────────

const (
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Bold   = "\033[1m"
	Reset  = "\033[0m"
)

func col(text, color string) string {
	return color + text + Reset
}

func errorf(format string, a ...interface{}) string {
	return col(fmt.Sprintf(format, a...), Red)
}

func successf(format string, a ...interface{}) string {
	return col(fmt.Sprintf(format, a...), Green)
}

func infof(format string, a ...interface{}) string {
	return col(fmt.Sprintf(format, a...), Yellow)
}

func headerf(format string, a ...interface{}) string {
	return Bold + Yellow + fmt.Sprintf(format, a...) + Reset
}

func dimf(format string, a ...interface{}) string {
	return col(fmt.Sprintf(format, a...), "\033[90m")
}

// ─── Connection metadata ──────────────────────────────────────────────────────

type connMeta struct {
	id   int
	addr string
	time time.Time
	conn net.Conn
}

// Server holds listener state and manages active connections.
type Server struct {
	port   int
	ln     net.Listener
	conns  []*connMeta
	mu     sync.Mutex
	active int
	nextID atomic.Int32
	stopCh chan struct{}
}

// Start launches the built-in TCP listener on the given port.
func Start(port int) error {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return fmt.Errorf("listen on :%d: %w", port, err)
	}

	s := &Server{
		port:   port,
		ln:     ln,
		stopCh: make(chan struct{}),
	}

	printBanner(port)

	// Signal handling for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Printf("\n%s\n", infof("Shutting down listener..."))
		s.shutdown()
	}()

	go s.acceptLoop()
	s.menuLoop()
	return nil
}

// ─── Banner ───────────────────────────────────────────────────────────────────

func printBanner(port int) {
	fmt.Println()
	fmt.Println(headerf("ShellCraft Built-in Listener"))
	fmt.Println(dimf(strings.Repeat("─", 60)))
	fmt.Printf("  %s\n", infof("Listening on :%d", port))
	fmt.Printf("  %s\n", dimf("Supports: bash, powershell, python, nc — any raw TCP reverse shell"))
	fmt.Printf("  %s\n", dimf("Interactive mode with PTY support, colored output"))
	fmt.Printf("  %s\n", dimf("Press Ctrl+C to stop the listener"))
	fmt.Println(dimf(strings.Repeat("─", 60)))
}

// ─── Connection acceptance loop ───────────────────────────────────────────────

func (s *Server) acceptLoop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			select {
			case <-s.stopCh:
			default:
			}
			return
		}

		s.mu.Lock()
		id := int(s.nextID.Add(1))
		meta := &connMeta{
			id:   id,
			addr: conn.RemoteAddr().String(),
			time: time.Now(),
			conn: conn,
		}
		s.conns = append(s.conns, meta)
		s.mu.Unlock()

		fmt.Printf("\n%s\n", successf("Incoming connection #%d from %s", id, meta.addr))
		fmt.Printf("  %s %s\n", infof("Connected:"), dimf(meta.time.Format("15:04:05")))
		fmt.Printf("  %s\n", infof("Type '%d' to interact with this connection", id))
	}
}

// ─── Menu loop ────────────────────────────────────────────────────────────────

func (s *Server) menuLoop() {
	stdin := bufio.NewScanner(os.Stdin)
	pending := false

	for {
		select {
		case <-s.stopCh:
			return
		default:
		}

		s.mu.Lock()
		count := len(s.conns)
		s.mu.Unlock()

		if count == 0 {
			if !pending {
				fmt.Printf("\r%s ", dimf("Waiting for connections..."))
				pending = true
			}
			time.Sleep(300 * time.Millisecond)
			continue
		}
		pending = false

		fmt.Printf("\n%s\n", headerf("Active Connections"))
		fmt.Println(dimf(strings.Repeat("─", 60)))
		s.mu.Lock()
		for _, m := range s.conns {
			mark := " "
			arrow := ""
			if m.id == s.active {
				mark = col("▸", Green)
				arrow = dimf(" ← active")
			}
			fmt.Printf("  %s [%d] %-21s %s%s\n", mark, m.id, m.addr, dimf("(%s ago)", time.Since(m.time).Round(time.Second)), arrow)
		}
		s.mu.Unlock()
		fmt.Println(dimf(strings.Repeat("─", 60)))
		fmt.Printf("  %s\n", dimf("Enter connection # to interact  |  'r' refresh  |  'q' quit"))

		fmt.Printf("%s ", col("listener>", Blue))
		if !stdin.Scan() {
			return
		}
		input := strings.TrimSpace(stdin.Text())

		switch {
		case input == "q" || input == "quit" || input == "exit":
			s.shutdown()
			return
		case input == "r" || input == "refresh":
			continue
		case input == "":
			continue
		default:
			id, err := strconv.Atoi(input)
			if err != nil {
				fmt.Printf("%s\n", errorf("Invalid input: %s", input))
				continue
			}
			s.interact(id)
		}
	}
}

// ─── Interactive session ──────────────────────────────────────────────────────

func (s *Server) interact(id int) {
	s.mu.Lock()
	var meta *connMeta
	for _, m := range s.conns {
		if m.id == id {
			meta = m
			break
		}
	}
	s.mu.Unlock()

	if meta == nil {
		fmt.Printf("%s\n", errorf("Connection #%d not found", id))
		return
	}

	s.active = id
	fmt.Printf("\n%s\n", headerf("Connected to %s [%d]", meta.addr, id))
	fmt.Printf("  %s\n", infof("Remote shell ready — all input is forwarded transparently"))
	fmt.Printf("  %s\n", dimf("Ctrl+] to return to menu  |  Ctrl+C sends SIGINT to remote"))
	fmt.Printf("  %s\n", dimf("Arrow keys, Tab, Ctrl+Z all work — raw PTY mode"))
	fmt.Println(dimf(strings.Repeat("─", 60)))

	// Attempt raw terminal mode via golang.org/x/term
	oldState, rawErr := term.MakeRaw(int(os.Stdin.Fd()))
	if rawErr == nil {
		s.interactRaw(meta, oldState)
	} else {
		fmt.Printf("  %s\n", infof("PTY not available, falling back to line mode"))
		s.interactLine(meta)
	}
}

// interactRaw provides a full PTY-aware interactive session.
func (s *Server) interactRaw(meta *connMeta, oldState *term.State) {
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	// Platform-specific signal handling (SIGWINCH on Unix, no-op on Windows)
	stopSig := setupWinchHandler()
	if stopSig != nil {
		defer stopSig()
	}

	done := make(chan struct{})
	var closeOnce sync.Once
	signalDone := func() { closeOnce.Do(func() { close(done) }) }

	// Connection → stdout (transparent pass-through — no color wrappers
	// so remote shell escape sequences and prompts are not corrupted)
	go func() {
		defer signalDone()
		buf := make([]byte, 4096)
		for {
			n, err := meta.conn.Read(buf)
			if n > 0 {
				os.Stdout.Write(buf[:n])
			}
			if err != nil {
				return
			}
		}
	}()

	// Stdin → connection (with escape detection: Ctrl+])
	go func() {
		defer signalDone()
		// Close the TCP conn on escape so the remote goroutine unblocks
		connClosed := false
		closeConn := func() {
			if !connClosed {
				connClosed = true
				meta.conn.Close()
			}
		}
		defer closeConn()

		buf := make([]byte, 4096)
		for {
			n, err := os.Stdin.Read(buf)
			if n > 0 {
				data := buf[:n]
				// Check for escape character Ctrl+] (0x1d)
				if idx := bytes.IndexByte(data, 0x1d); idx >= 0 {
					// Forward bytes before the escape
					if idx > 0 {
						meta.conn.Write(data[:idx])
					}
					return
				}
				meta.conn.Write(data)
			}
			if err != nil {
				return
			}
		}
	}()

	// Block until session ends
	<-done

	// Restore terminal before printing exit message
	term.Restore(int(os.Stdin.Fd()), oldState)

	fmt.Printf("\n%s\n", dimf("Connection #%d session ended (returning to menu)", meta.id))
	fmt.Println(dimf(strings.Repeat("─", 60)))
}

// interactLine is a plain line-oriented mode fallback when PTY is unavailable.
func (s *Server) interactLine(meta *connMeta) {
	done := make(chan struct{})
	var closeOnce sync.Once
	signalDone := func() { closeOnce.Do(func() { close(done) }) }

	// Connection → stdout
	go func() {
		defer signalDone()
		buf := make([]byte, 4096)
		for {
			n, err := meta.conn.Read(buf)
			if n > 0 {
				os.Stdout.Write(buf[:n])
			}
			if err != nil {
				return
			}
		}
	}()

	// Stdin → connection (line by line, check for exit/quit)
	go func() {
		defer signalDone()
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			trimmed := strings.TrimSpace(line)
			if trimmed == "exit" || trimmed == "quit" {
				meta.conn.Write([]byte(line + "\n"))
				return
			}
			meta.conn.Write([]byte(line + "\n"))
		}
	}()

	<-done

	fmt.Printf("\n%s\n", dimf("Connection #%d session ended", meta.id))
}

// ─── Graceful shutdown ────────────────────────────────────────────────────────

func (s *Server) shutdown() {
	s.ln.Close()
	s.mu.Lock()
	for _, c := range s.conns {
		c.conn.Close()
	}
	s.mu.Unlock()
	select {
	case <-s.stopCh:
		// Already closed
	default:
		close(s.stopCh)
	}
}
