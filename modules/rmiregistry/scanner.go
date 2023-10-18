package rmiregistry

import (
	"fmt"
	"net"
	"strings"
    	"net/rpc"
	"regexp"
	"time"
	log "github.com/sirupsen/logrus"
	"github.com/zmap/zgrab2"
)

// ScanResults is the output of the RMI Registry scan.
type ScanResults struct {
	Banner string `json:"banner,omitempty"`
}

// Flags are the RMI Registry-specific command-line flags.
type Flags struct {
	zgrab2.BaseFlags
	Verbose bool `long:"verbose" description:"More verbose logging, include debug fields in the scan results"`
}

// Module implements the zgrab2.Module interface for RMI Registry.
type Module struct {
}

// Scanner implements the zgrab2.Scanner interface and holds the state for a single scan.
type Scanner struct {
	config *Flags
}

// Connection holds the state for a single connection to the RMI Registry server.
type Connection struct {
	config  *Flags
	results ScanResults
	conn    net.Conn
	buffer  [10000]byte
}

// RegisterModule registers the rmiregistry zgrab2 module.
func RegisterModule() {
	var module Module
	_, err := zgrab2.AddCommand("rmiregistry", "RMI Registry", module.Description(), 1099, &module)
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
	return "Grab an RMI Registry banner"
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
	return "rmiregistry"
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

// rmiEndRegex matches zero or more lines followed by a numeric RMI Registry status code and linebreak, e.g., "RMI Registry 1099\r\n"
var rmiEndRegex = regexp.MustCompile(`^(?:.*\r?\n)*RMI Registry ([0-9]+)\r?\n$`)

// isOKResponse returns true if and only if the given response code indicates success (e.g., 1099)
func (rmi *Connection) isOKResponse(retCode string) bool {
	return strings.HasPrefix(retCode, "1099")
}

// readResponse reads an RMI Registry response chunk from the server.
// It returns the full response, as well as the status code alone.
func (rmi *Connection) readResponse() (string, string, error) {
	respLen, err := zgrab2.ReadUntilRegex(rmi.conn, rmi.buffer[:], rmiEndRegex)
	if err != nil {
		return "", "", err
	}
	ret := string(rmi.buffer[0:respLen])
	retCode := rmiEndRegex.FindStringSubmatch(ret)[1]
	return ret, retCode, nil
}

// GetRMIBanner reads the data sent by the server immediately after connecting.
// Returns true if and only if the server returns a success status code (1099).
func (rmi *Connection) GetRMIBanner() (bool, error) {
	banner, retCode, err := rmi.readResponse()
	if err != nil {
		return false, err
	}
	rmi.results.Banner = banner
	return rmi.isOKResponse(retCode), nil
}

// Scan performs the configured scan on the RMI Registry server, as follows:
// * Read the banner into results.Banner (if it is "RMI Registry 1099", it's considered a success).
type RMIArgs struct {
    Arg1 string
    Arg2 int
}

// Define a sample structure for RMI method results.


type YourRMIResult struct {
    // Define fields specific to your method's result.
    ResultField1 string
    ResultField2 int
}


func (s *Scanner) Scan(target zgrab2.ScanTarget) (status zgrab2.ScanStatus, result interface{}, thrown error) {
    var err error

    // Set the timeout duration to a longer value (e.g., 10 seconds).
    timeoutDuration := 50 * time.Second

    // Create a dialer with the specified timeout.
    dialer := &net.Dialer{
        Timeout: timeoutDuration,
    }

    // Build the address for the connection using the target's IP and Port.
    address := fmt.Sprintf("%s:%d", target.IP, target.Port)

    // Establish a network connection with the specified timeout.
    conn, err := dialer.Dial("tcp", address)
    if err != nil {
        fmt.Println("Error opening connection:", err)
        return zgrab2.TryGetScanStatus(err), nil, err
    }
    cn := conn
    defer func() {
        cn.Close()
    }()

    // Assuming the "java-rmi GNU Classpath grmiregistry" service is running on the target, we can create an RPC client.
    client := rpc.NewClient(cn)

    // Replace with your specific details:
    // - RMI Method Name (Assuming you want to invoke a method on the registry itself).
    rmiMethodName := "YourRMIMethodName"

    // - RMI Method Arguments (if any).
     rmiMethodArgs := RMIArgs{
        Arg1: "SampleArgument1",
        Arg2: 42,
    }
    // - RMI Result Type: Use the specific result type you defined.
    var rmiResult YourRMIResult

    // Call the RMI method on the specified object (registry).
    err = client.Call(rmiMethodName, rmiMethodArgs, &rmiResult)
    if err != nil {
        fmt.Println("Error calling RMI method:", err)
        return zgrab2.SCAN_UNKNOWN_ERROR, nil, err
    }

    // Process the RMI result here.

    return zgrab2.SCAN_SUCCESS, rmiResult, nil
}

