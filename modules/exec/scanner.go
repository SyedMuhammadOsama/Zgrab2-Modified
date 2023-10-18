package exec

import (
	"context"
	"io"
	"net"
	"regexp"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/zmap/zgrab2"
)

// ScanResults is the output of the scan.
type ScanResults struct {
	Response string `json:"response,omitempty"`
	Service  string `json:"service,omitempty"`
	Version  string `json:"version,omitempty"`
}

// Flags holds the command-line flags for the EXEC module.
type Flags struct {
	zgrab2.BaseFlags
}

// Module implements the zgrab2.Module interface.
type Module struct {
}

// Scanner implements the zgrab2.Scanner interface.
type Scanner struct {
	config *Flags
}

// Connection holds the state for a single connection to the EXEC service.
type Connection struct {
	config  *Flags
	results ScanResults
	conn    net.Conn
}

// RegisterModule registers the EXEC ZGrab2 module.
func RegisterModule() {
	var module Module
	_, err := zgrab2.AddCommand("exec", "EXEC", module.Description(), 512, &module)
	if err != nil {
		log.Fatal(err)
	}
}

// NewFlags returns the default flags object.
func (m *Module) NewFlags() interface{} {
	return new(Flags)
}

// NewScanner creates a new Scanner instance.
func (m *Module) NewScanner() zgrab2.Scanner {
	return new(Scanner)
}

// Description returns an overview of this module.
func (m *Module) Description() string {
	return "Grab an EXEC banner"
}

// Validate checks if the command-line flags are valid.
func (f *Flags) Validate(args []string) error {
	return nil
}

// Help returns the help string for the EXEC module.
func (f *Flags) Help() string {
	return ""
}

// Protocol returns the protocol identifier for the scanner.
func (s *Scanner) Protocol() string {
	return "exec"
}

// Init initializes the Scanner instance with the command-line flags.
func (s *Scanner) Init(flags zgrab2.ScanFlags) error {
	f, _ := flags.(*Flags)
	s.config = f
	return nil
}

func (s *Scanner) InitPerSender(senderID int) error {
	return nil
}


// GetName returns the configured name for the Scanner.
func (s *Scanner) GetName() string {
	return s.config.Name
}

// GetTrigger returns the Trigger defined in the Flags.
func (scanner *Scanner) GetTrigger() string {
	return scanner.config.Trigger
}

// Scan performs the configured scan on the EXEC service.
// Scan performs the exec scan.
func (scanner *Scanner) Scan(t zgrab2.ScanTarget) (zgrab2.ScanStatus, interface{}, error) {
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := t.Open(&scanner.config.BaseFlags)
	if err != nil {
		return zgrab2.TryGetScanStatus(err), nil, err
	}
	defer conn.Close()

	results := ScanResults{}
	execConn := Connection{conn: conn, config: scanner.config, results: results}

	// Customize your EXEC protocol handshake and response handling here
	execRequest := "YOUR_EXEC_REQUEST"
	_, err = io.WriteString(execConn.conn, execRequest+"\n")
	if err != nil {
		return zgrab2.SCAN_PROTOCOL_ERROR, nil, err
	}

	responseBuffer := make([]byte, 4096)
	n, err := execConn.conn.Read(responseBuffer)
	if err != nil {
		return zgrab2.SCAN_PROTOCOL_ERROR, nil, err
	}

	execConn.results.Response = strings.TrimSpace(string(responseBuffer[:n]))

	// Parse the response to extract service and version information
	service, version := extractServiceAndVersion(execConn.results.Response)

	// Set the extracted service and version in ScanResults
	execConn.results.Service = service
	execConn.results.Version = version

	return zgrab2.SCAN_SUCCESS, &execConn.results, nil
}

// extractServiceAndVersion extracts service and version from a response using regular expressions.
func extractServiceAndVersion(response string) (string, string) {
	// Define regular expressions for matching service and version
	servicePattern := regexp.MustCompile(`Service: ([^\s]+)`)
	versionPattern := regexp.MustCompile(`Version: ([^\s]+)`)

	// Find matches for service and version
	serviceMatches := servicePattern.FindStringSubmatch(response)
	versionMatches := versionPattern.FindStringSubmatch(response)

	// Extract service and version if matches are found
	var service, version string
	if len(serviceMatches) == 2 {
		service = serviceMatches[1]
	}
	if len(versionMatches) == 2 {
		version = versionMatches[1]
	}

	return service, version
}
