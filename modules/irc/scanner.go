package irc

import (
	"fmt"
	"net"
	log "github.com/sirupsen/logrus"

	"github.com/zmap/zgrab2"
)

// ScanResults is the output of the IRC scan.
type ScanResults struct {
	Banner      string   `json:"banner,omitempty"`
	Joined      bool     `json:"joined,omitempty"`
	JoinChannel string   `json:"join_channel,omitempty"`
	Responses   []string `json:"responses,omitempty"`
}

// Flags are the IRC-specific command-line flags.
type Flags struct {
	zgrab2.BaseFlags
}

// Module implements the zgrab2.Module interface for IRC scanning.
type Module struct {
}

// Scanner implements the zgrab2.Scanner interface for IRC scanning.
type Scanner struct {
	config *Flags
}

// Connection holds the state for a single connection to the IRC server.
type Connection struct {
	conn    net.Conn
	config  *Flags
	results ScanResults
}

// RegisterModule registers the IRC zgrab2 module.
func RegisterModule() {
	var module Module
	_, err := zgrab2.AddCommand("irc", "IRC", module.Description(), 6667, &module)
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
	return "Grab an IRC banner, join a channel, and retrieve detailed responses"
}

// Validate flags
func (f *Flags) Validate(args []string) (err error) {
	return
}

// Help returns this module's help string.
func (f *Flags) Help() string {
	return ""
}

// Protocol returns the protocol identifier for the scanner.
func (s *Scanner) Protocol() string {
	return "irc"
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

// Scan performs the configured scan on the IRC server.
func (s *Scanner) Scan(t zgrab2.ScanTarget) (status zgrab2.ScanStatus, result interface{}, thrown error) {
	var err error
	conn, err := t.Open(&s.config.BaseFlags)
	if err != nil {
		return zgrab2.TryGetScanStatus(err), nil, err
	}
	cn := conn
	defer func() {
		cn.Close()
	}()

	results := ScanResults{}
	irc := Connection{conn: cn, config: s.config, results: results}

	// Read the IRC server's banner
	buffer := make([]byte, 1024)
	n, err := irc.conn.Read(buffer)
	if err != nil {
		return zgrab2.TryGetScanStatus(err), &irc.results, err
	}
	irc.results.Banner = string(buffer[:n])

	// Optionally, you can perform additional IRC operations here.
	// For example, joining an IRC channel:
	// Uncomment the following lines to join a channel and customize the channel name.
	channelName := "mychannel"
	err = irc.JoinChannel(channelName)
	if err != nil {
	    return zgrab2.SCAN_APPLICATION_ERROR, &irc.results, err
	}
	irc.results.Joined = true
	irc.results.JoinChannel = channelName

	// Perform additional IRC operations and retrieve responses.
	// Customize the IRC commands and responses as needed.
	// Example:
	responses, err := irc.SendAndReceive([]string{"NICK bot", "USER bot 0 * :bot"})
	if err != nil {
	    return zgrab2.SCAN_APPLICATION_ERROR, &irc.results, err
	}
	irc.results.Responses = responses

	return zgrab2.SCAN_SUCCESS, &irc.results, nil
}

// JoinChannel sends a JOIN command to join an IRC channel.
func (irc *Connection) JoinChannel(channel string) error {
	command := fmt.Sprintf("JOIN %s", channel)
	_, err := irc.sendCommand(command)
	if err != nil {
		return err
	}
	return nil
}

// SendAndReceive sends a list of IRC commands and reads the responses.
func (irc *Connection) SendAndReceive(commands []string) ([]string, error) {
	var responses []string
	for _, cmd := range commands {
		response, err := irc.sendCommand(cmd)
		if err != nil {
			return responses, err
		}
		responses = append(responses, response)
	}
	return responses, nil
}

// sendCommand sends a command to the IRC server and reads the response.
func (irc *Connection) sendCommand(command string) (string, error) {
	command = fmt.Sprintf("%s\r\n", command)
	_, err := irc.conn.Write([]byte(command))
	if err != nil {
		return "", err
	}

	// Read the response from the IRC server
	buffer := make([]byte, 1024)
	n, err := irc.conn.Read(buffer)
	if err != nil {
		return "", err
	}
	response := string(buffer[:n])
	return response, nil
}

