package rpcbind

import (
	"net"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/zmap/zgrab2"
)

// ScanResults is the output of the RPCBIND scan.
type ScanResults struct {
	Banner string `json:"banner,omitempty"`
}

// Flags are the RPCBIND-specific command-line flags.
type Flags struct {
	zgrab2.BaseFlags
	Verbose bool `long:"verbose" description:"More verbose logging, include debug fields in the scan results"`
}

// Module implements the zgrab2.Module interface for RPCBIND scanning.
type Module struct{}

// Scanner implements the zgrab2.Scanner interface and holds the state for a single scan.
type Scanner struct {
	config *Flags
}

// Connection holds the state for a single connection to the RPCBIND service.
type Connection struct {
	config  *Flags
	results ScanResults
	conn    net.Conn
}

// RegisterModule registers the RPCBIND zgrab2 module.
func RegisterModule() {
	var module Module
	_, err := zgrab2.AddCommand("rpcbind", "RPCBIND", module.Description(), 111, &module)
	if err != nil {
		log.Fatal(err)
	}
}

// NewFlags returns the default flags object to be filled in with the command-line arguments.
func (m *Module) NewFlags() interface{} {
	return new(Flags)
}

// NewScanner returns a new Scanner instance.
func (m *Module) NewScanner() zgrab2.Scanner {
	return new(Scanner)
}

// Description returns an overview of this module.
func (m *Module) Description() string {
	return "Grab an RPCBIND banner"
}

// Validate flags
func (f *Flags) Validate(args []string) error {
	return nil
}

// Help returns this module's help string.
func (f *Flags) Help() string {
	return ""
}

// Protocol returns the protocol identifier for the scanner.
func (s *Scanner) Protocol() string {
	return "rpcbind"
}

// Init initializes the Scanner instance with the flags from the command line.
func (s *Scanner) Init(flags zgrab2.ScanFlags) error {
	f, _ := flags.(*Flags)
	s.config = f
	return nil
}

// InitPerSender does nothing in this module.
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

// Scan performs the RPCBIND banner grabbing and scanning.
func (s *Scanner) Scan(t zgrab2.ScanTarget) (status zgrab2.ScanStatus, result interface{}, thrown error) {
    var err error
    conn, err := t.Open(&s.config.BaseFlags)
    if err != nil {
        return zgrab2.TryGetScanStatus(err), nil, err
    }
    defer conn.Close()

    results := ScanResults{}
    rpcbindConn := Connection{conn: conn, config: s.config, results: results}

    // Read the banner from the RPCBIND service
    _, err = rpcbindConn.GetBanner()
    if err != nil {
        return zgrab2.TryGetScanStatus(err), nil, err
    }

    return zgrab2.SCAN_SUCCESS, &rpcbindConn.results, nil
}

// GetBanner reads the banner sent by the RPCBIND service.
func (rpc *Connection) GetBanner() (string, error) {
	// Customize your RPCBIND banner grabbing logic here
	// You can read data from the connection and extract the banner
	// The following code is a placeholder and should be adapted to RPCBIND's behavior
	buffer := make([]byte, 4096)
	n, err := rpc.conn.Read(buffer)
	if err != nil {
		return "", err
	}

	banner := strings.TrimSpace(string(buffer[:n]))
	rpc.results.Banner = banner

	return banner, nil
}
