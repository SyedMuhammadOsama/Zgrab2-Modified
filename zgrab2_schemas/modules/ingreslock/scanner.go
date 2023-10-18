package ingreslock

import (
	"net"
	"strings"
	"regexp"

	log "github.com/sirupsen/logrus"
	"github.com/zmap/zgrab2"
)

// ScanResults is the output of the scan.
type ScanResults struct {
	Banner        string `json:"banner,omitempty"`
}

// Flags are the specific command-line flags for the ingreslock module.
type Flags struct {
	zgrab2.BaseFlags
}

// Module implements the zgrab2.Module interface.
type Module struct{}

// Scanner implements the zgrab2.Scanner interface.
type Scanner struct {
	config *Flags
}

// Connection holds the state for a single connection to the ingreslock service.
type Connection struct {
	config  *Flags
	results ScanResults
	conn    net.Conn
}

// RegisterModule registers the ingreslock zgrab2 module.
func RegisterModule() {
	var module Module
	_, err := zgrab2.AddCommand("ingreslock", "Ingreslock", module.Description(), 1524, &module)
	if err != nil {
		log.Fatal(err)
	}
}

// NewFlags returns the default flags object for the ingreslock module.
func (m *Module) NewFlags() interface{} {
	return new(Flags)
}

// NewScanner returns a new Scanner instance.
func (m *Module) NewScanner() zgrab2.Scanner {
	return new(Scanner)
}

// Description returns an overview of this module.
func (m *Module) Description() string {
	return "Grab an ingreslock banner and execute commands"
}

// Validate validates the ingreslock module's flags.
func (f *Flags) Validate(args []string) error {
	return nil
}

// Help returns the help string for the ingreslock module's flags.
func (f *Flags) Help() string {
	return ""
}

// Protocol returns the protocol identifier for the scanner.
func (s *Scanner) Protocol() string {
	return "ingreslock"
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

// ingreslockEndRegex matches zero or more lines followed by a numeric status code.
var ingreslockEndRegex = regexp.MustCompile(`^.*([0-9]+)$`)

// readResponse reads the banner response from the ingreslock service.
func (ingres *Connection) readResponse() (string, error) {
	buffer := make([]byte, 4096) // Create a buffer to read the response
	n, err := ingres.conn.Read(buffer)
	if err != nil {
		return "", err
	}
	return string(buffer[:n]), nil
}

// GetIngreslockBannerAndCommands performs the banner grab on port 1524/tcp (ingreslock)
// and executes commands 'ls' and 'whoami'.
func (ingres *Connection) GetIngreslockBannerAndCommands() error {
	response, err := ingres.readResponse()
	if err != nil {
		return err
	}
	ingres.results.Banner = strings.TrimSpace(response)


	return nil
}

// Scan performs the configured scan on the ingreslock service (port 1524/tcp).
func (s *Scanner) Scan(t zgrab2.ScanTarget) (status zgrab2.ScanStatus, result interface{}, thrown error) {
	conn, err := t.Open(&s.config.BaseFlags)
	if err != nil {
		return zgrab2.TryGetScanStatus(err), nil, err
	}
	cn := conn
	defer func() {
		cn.Close()
	}()

	results := ScanResults{}
	ingres := Connection{conn: cn, config: s.config, results: results}

	if err := ingres.GetIngreslockBannerAndCommands(); err != nil {
		return zgrab2.TryGetScanStatus(err), &ingres.results, err
	}

	return zgrab2.SCAN_SUCCESS, &ingres.results, nil
}
