package distccd

import (
	"net"
	"strings"
	"log"
	"regexp"
	"github.com/zmap/zgrab2"
)




func (f *Flags) Validate(args []string) error {
	// Implement the Validate method as needed
	return nil
}

type Flags struct {
	zgrab2.BaseFlags
	zgrab2.TLSFlags

	// Your other flags here, if needed
	Cmd string `long:"cmd" description:"Command to execute on the remote distccd server"`
}


// ScanResults represents the results of the distccd scan.
type ScanResults struct {
	Banner      string `json:"banner,omitempty"`
	ServiceName string `json:"service_name"`
	Version     string `json:"version"`
}

// Module implements the zgrab2.Module interface.
type Module struct{}

// Scanner implements the zgrab2.Scanner interface.
type Scanner struct {
	config *Flags
}
func (s *Scanner) GetTrigger() string {
	return ""
}
func (f *Flags) Help() string {
	// Implement the Help method as needed
	return ""
}


// Connection holds the state for a single connection to the distccd server.
type Connection struct {
	config  *Flags
	results ScanResults
	conn    net.Conn
}

// RegisterModule registers the distccd zgrab2 module.
func RegisterModule() {
	var module Module
	_, err := zgrab2.AddCommand("distccd", "distccd", module.Description(), 3632, &module)
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
	return "Grab a distccd banner"
}

// Protocol returns the protocol identifier for the scanner.
func (s *Scanner) Protocol() string {
	return "distccd"
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


func (s *Scanner) Scan(t zgrab2.ScanTarget) (zgrab2.ScanStatus, interface{}, error) {
	results := &ScanResults{}

	conn, err := t.Open(&s.config.BaseFlags)
	if err != nil {
		return zgrab2.TryGetScanStatus(err), results, err
	}
	defer conn.Close()

	// Send a probe to the distccd service
	probe := []byte("HELO\n") // Adjust the probe as needed
	_, err = conn.Write(probe)
	if err != nil {
		return zgrab2.SCAN_PROTOCOL_ERROR, results, err
	}

	// Read the response from the server
response := make([]byte, 1024) // Adjust buffer size as needed
_, err = conn.Read(response)

if err != nil {
    return zgrab2.TryGetScanStatus(err), results, err
}


	// Process the response to extract service and version information
	//banner := string(response[:n])
	// ...

// Call extractServiceAndVersion to get service name and version
serviceName, version := extractServiceAndVersion(results.Banner)

// Update the results with the obtained service name and version
results.ServiceName = serviceName
results.Version = version

// Return a successful scan status, updated results, and no errors
return zgrab2.SCAN_SUCCESS, results, nil


}
func extractServiceAndVersion(banner string) (string, string) {
	serviceName := "unknown"
	version := "unknown"

	// Split the banner into lines
	lines := strings.Split(banner, "\n")

	// Loop through lines to find service and version information
	for _, line := range lines {
		if strings.Contains(line, "distccd") {
			serviceName = "distccd"
			// Implement logic to extract the version using regular expressions
			versionPattern := regexp.MustCompile(`version ([\d.]+)`)
			match := versionPattern.FindStringSubmatch(line)
			if len(match) > 1 {
				version = match[1]
			}
			break // Exit loop after finding the service name and version
		}
	}

	return serviceName, version
}

