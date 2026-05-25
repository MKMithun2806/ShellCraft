package listener

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

// ─── ANSI colors (self-contained, zero external deps) ───────────────────────

const (
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Cyan   = "\033[36m"
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

// ─── Connection tracking ────────────────────────────────────────────────────

type connMeta struct {
	id   int
	addr string
	time time.Time
	conn net.Conn
}

// Server holds the listener state and manages active connections.
type Server struct {
	port   int
	ln     net.Listener
	conns  []*connMeta
	mu     sync.Mutex
	active int
	nextID int
	stopCh chan struct{}
}

// Start launches the built-in listener on the given port.
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

	fmt.Println()
	fmt.Println(headerf("ShellCraft Built-in Listener"))
	fmt.Println(dimf(strings.Repeat("─", 55)))
	fmt.Printf("  %s\n", infof("Listening on :%d", port))
	fmt.Printf("  %s\n", dimf("Supports: bash, powershell, python, nc — any raw TCP reverse shell"))
	fmt.Printf("  %s\n", dimf("Press Ctrl+C to stop"))
	fmt.Println(dimf(strings.Repeat("─", 55)))

	// Signal handling for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Printf("\n%s\n", infof("Shutting down listener..."))
		s.ln.Close()
		s.mu.Lock()
		for _, c := range s.conns {
			c.conn.Close()
		}
		s.mu.Unlock()
		close(s.stopCh)
	}()

	go s.acceptLoop()
	s.menuLoop()
	return nil
}

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
		s.nextID++
		id := s.nextID
		meta := &connMeta{
			id:   id,
			addr: conn.RemoteAddr().String(),
			time: time.Now(),
			conn: conn,
		}
		s.conns = append(s.conns, meta)
		s.mu.Unlock()

		fmt.Printf("\n%s\n", successf("Incoming connection #%d from %s", id, meta.addr))
		fmt.Printf("  %s\n", infof("Type '%d' to interact with this connection", id))
	}
}

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
		fmt.Println(dimf(strings.Repeat("─", 55)))
		s.mu.Lock()
		for _, m := range s.conns {
			mark := " "
			if m.id == s.active {
				mark = col(">", Green)
			}
			fmt.Printf("  %s [%d] %-21s %s\n", mark, m.id, m.addr, dimf("(%s ago)", time.Since(m.time).Round(time.Second)))
		}
		s.mu.Unlock()
		fmt.Println(dimf(strings.Repeat("─", 55)))
		fmt.Printf("  %s\n", dimf("Enter connection # to interact, 'r' refresh, 'q' quit"))

		fmt.Printf("%s ", col("listener>", Blue))
		if !stdin.Scan() {
			return
		}
		input := strings.TrimSpace(stdin.Text())

		switch {
		case input == "q" || input == "quit" || input == "exit":
			fmt.Printf("%s\n", infof("Shutting down..."))
			s.ln.Close()
			s.mu.Lock()
			for _, c := range s.conns {
				c.conn.Close()
			}
			s.mu.Unlock()
			close(s.stopCh)
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

// interact opens an interactive session with the given connection.
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
	fmt.Printf("  %s\n", infof("Raw TCP mode — paste or type commands"))
	fmt.Printf("  %s\n", dimf("Ctrl+C or 'exit' to return to menu"))
	fmt.Println(dimf(strings.Repeat("─", 55)))

	// Try to enable raw terminal mode for a better shell experience
	restore, rawErr := rawTerminal()
	if rawErr == nil {
		defer restore()
	}

	errCh := make(chan error, 2)

	// Connection → stdout
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := meta.conn.Read(buf)
			if n > 0 {
				os.Stdout.Write(buf[:n])
			}
			if err != nil {
				errCh <- err
				return
			}
		}
	}()

	// Stdin → connection
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := os.Stdin.Read(buf)
			if n > 0 {
				data := buf[:n]
				// Check for exit on non-raw mode (line-oriented)
				if !isRaw.Load() {
					trimmed := strings.TrimSpace(string(data))
					if trimmed == "exit" || trimmed == "quit" {
						meta.conn.Write(data)
						errCh <- io.EOF
						return
					}
				}
				meta.conn.Write(data)
			}
			if err != nil {
				errCh <- err
				return
			}
		}
	}()

	// Wait for either side to close
	select {
	case <-errCh:
	case <-s.stopCh:
	}

	// Drain
	select {
	case <-errCh:
	default:
	}

	if rawErr == nil {
		restore()
	}
	fmt.Printf("\n%s\n", dimf("Connection #%d session ended", id))
}

// ─── Raw terminal (Unix only, fallback on other platforms) ──────────────────

var isRaw atomicBool

type atomicBool struct {
	v int32
}

func (b *atomicBool) Load() bool {
	return b.v != 0
}

func (b *atomicBool) Store(v bool) {
	if v {
		b.v = 1
	} else {
		b.v = 0
	}
}

func rawTerminal() (func(), error) {
	fd := int(os.Stdin.Fd())

	// Verify stdin is a terminal
	var termios syscall.Termios
	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), syscall.TCGETS, uintptr(unsafe.Pointer(&termios)), 0, 0, 0); err != 0 {
		return func() {}, fmt.Errorf("not a terminal: %v", err)
	}

	oldState := termios
	newState := termios
	newState.Iflag &^= syscall.IGNBRK | syscall.BRKINT | syscall.PARMRK | syscall.ISTRIP | syscall.INLCR | syscall.IGNCR | syscall.ICRNL | syscall.IXON
	newState.Oflag &^= syscall.OPOST
	newState.Lflag &^= syscall.ECHO | syscall.ECHONL | syscall.ICANON | syscall.ISIG | syscall.IEXTEN
	newState.Cflag &^= syscall.CSIZE | syscall.PARENB
	newState.Cflag |= syscall.CS8
	newState.Cc[syscall.VMIN] = 1
	newState.Cc[syscall.VTIME] = 0

	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), syscall.TCSETS, uintptr(unsafe.Pointer(&newState)), 0, 0, 0); err != 0 {
		return func() {}, fmt.Errorf("tcsetattr: %v", err)
	}

	isRaw.Store(true)

	return func() {
		isRaw.Store(false)
		syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), syscall.TCSETS, uintptr(unsafe.Pointer(&oldState)), 0, 0, 0)
	}, nil
}
