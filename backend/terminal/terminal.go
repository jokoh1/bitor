package terminal

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"golang.org/x/crypto/ssh"
)

type TerminalMessage struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type TerminalSession struct {
	ws      *websocket.Conn
	client  *ssh.Client
	session *ssh.Session
	stdin   io.WriteCloser
	stdout  io.Reader
	stderr  io.Reader
	mutex   sync.Mutex
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins in development
	},
}

// HandleTerminalConnection handles WebSocket connections for terminal access
func HandleTerminalConnection(app *pocketbase.PocketBase) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get scan ID and IP from query params
		scanID := c.QueryParam("scanId")
		if scanID == "" {
			return fmt.Errorf("scan ID is required")
		}

		log.Printf("Handling terminal connection for scan ID: %s", scanID)

		// Get scan record to verify IP and get SSH key
		scan, err := app.Dao().FindRecordById("nuclei_scans", scanID)
		if err != nil {
			log.Printf("Failed to find scan: %v", err)
			return fmt.Errorf("failed to find scan: %v", err)
		}

		// Check scan status
		status := scan.GetString("status")
		if status != "Running" {
			log.Printf("Scan is not in Running state (current status: %s)", status)
			return fmt.Errorf("cannot connect to terminal: scan is not running (status: %s)", status)
		}

		ipAddress := scan.GetString("ip_address")
		if ipAddress == "" {
			log.Printf("Scan has no IP address")
			return fmt.Errorf("scan has no IP address")
		}

		log.Printf("Found IP address: %s", ipAddress)

		// Get SSH key from ansible collection
		sshKeys, err := app.Dao().FindRecordsByExpr("ansible")
		if err != nil {
			log.Printf("Failed to find SSH keys: %v", err)
			return fmt.Errorf("failed to find SSH keys: %v", err)
		}

		var privateKey string
		for _, record := range sshKeys {
			if key := record.GetString("ssh_private_key"); key != "" {
				privateKey = key
				break
			}
		}

		if privateKey == "" {
			log.Printf("No SSH private key found")
			return fmt.Errorf("no SSH private key found")
		}

		// Parse private key
		signer, err := ssh.ParsePrivateKey([]byte(privateKey))
		if err != nil {
			log.Printf("Failed to parse private key: %v", err)
			return fmt.Errorf("failed to parse private key: %v", err)
		}

		// Configure SSH client
		config := &ssh.ClientConfig{
			User: "ansible",
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(signer),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         10 * time.Second,
		}

		log.Printf("Attempting to connect to SSH server at %s:22 with user 'ansible'", ipAddress)

		// Connect to SSH server
		client, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", ipAddress), config)
		if err != nil {
			log.Printf("Failed to connect as ansible user: %v", err)
			return fmt.Errorf("failed to connect as ansible user: %v", err)
		}

		log.Printf("Successfully connected as ansible user")

		// Create new SSH session
		session, err := client.NewSession()
		if err != nil {
			client.Close()
			log.Printf("Failed to create session: %v", err)
			return fmt.Errorf("failed to create session: %v", err)
		}

		// Request pseudo terminal
		if err := session.RequestPty("xterm", 40, 80, ssh.TerminalModes{}); err != nil {
			session.Close()
			client.Close()
			log.Printf("Failed to request pty: %v", err)
			return fmt.Errorf("failed to request pty: %v", err)
		}

		// Get pipes for stdin/stdout
		stdin, err := session.StdinPipe()
		if err != nil {
			session.Close()
			client.Close()
			log.Printf("Failed to get stdin pipe: %v", err)
			return fmt.Errorf("failed to get stdin pipe: %v", err)
		}

		stdout, err := session.StdoutPipe()
		if err != nil {
			session.Close()
			client.Close()
			log.Printf("Failed to get stdout pipe: %v", err)
			return fmt.Errorf("failed to get stdout pipe: %v", err)
		}

		stderr, err := session.StderrPipe()
		if err != nil {
			session.Close()
			client.Close()
			log.Printf("Failed to get stderr pipe: %v", err)
			return fmt.Errorf("failed to get stderr pipe: %v", err)
		}

		// Start shell
		if err := session.Shell(); err != nil {
			session.Close()
			client.Close()
			log.Printf("Failed to start shell: %v", err)
			return fmt.Errorf("failed to start shell: %v", err)
		}

		log.Printf("Successfully started shell session")

		// Upgrade HTTP connection to WebSocket
		ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			session.Close()
			client.Close()
			log.Printf("Failed to upgrade connection: %v", err)
			return fmt.Errorf("failed to upgrade connection: %v", err)
		}

		log.Printf("Successfully upgraded to WebSocket connection")

		// Create terminal session
		term := &TerminalSession{
			ws:      ws,
			client:  client,
			session: session,
			stdin:   stdin,
			stdout:  stdout,
			stderr:  stderr,
		}

		// Handle WebSocket messages
		go term.handleWebSocket()
		go term.handleSSHOutput()

		return nil
	}
}

func (t *TerminalSession) handleWebSocket() {
	defer func() {
		log.Printf("Closing WebSocket connection")
		t.ws.Close()
		t.session.Close()
		t.client.Close()
	}()

	for {
		// Read message from WebSocket
		_, data, err := t.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("WebSocket read error: %v", err)
			}
			return
		}

		var msg TerminalMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}

		switch msg.Type {
		case "input":
			t.mutex.Lock()
			_, err := t.stdin.Write([]byte(msg.Data))
			t.mutex.Unlock()
			if err != nil {
				log.Printf("Failed to write to stdin: %v", err)
				return
			}
		case "resize":
			// TODO: Implement terminal resize
			log.Printf("Resize event received but not implemented")
		default:
			log.Printf("Unknown message type received: %s", msg.Type)
		}
	}
}

func (t *TerminalSession) handleSSHOutput() {
	defer func() {
		log.Printf("Closing SSH connection")
		t.ws.Close()
		t.session.Close()
		t.client.Close()
	}()

	// Create a buffer for reading
	buf := make([]byte, 4096)

	for {
		// Read from stdout
		n, err := t.stdout.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Printf("Failed to read from stdout: %v", err)
			} else {
				log.Printf("SSH connection closed (EOF)")
			}
			return
		}

		if n > 0 {
			msg := TerminalMessage{
				Type: "output",
				Data: string(buf[:n]),
			}

			data, err := json.Marshal(msg)
			if err != nil {
				log.Printf("Failed to marshal message: %v", err)
				continue
			}

			t.mutex.Lock()
			if err := t.ws.WriteMessage(websocket.TextMessage, data); err != nil {
				t.mutex.Unlock()
				log.Printf("Failed to write to WebSocket: %v", err)
				return
			}
			t.mutex.Unlock()
		}
	}
}
