package lib

import (
	"encoding/json"
	"net"
	"os"
	"strings"
)

// NetConfig holds the hostname and the list of IP addresses for the current host
type NetConfig struct {
	Hostname string   `json:"hostname"`
	IPs      []string `json:"ips"`
}

// String returns the string representation of the network configuration
func (n *NetConfig) String() string {
	sb := strings.Builder{}
	sb.WriteString("ServerHost: ")
	sb.WriteString(n.Hostname)
	sb.WriteString("\nIPs:\n")
	sb.WriteString(strings.Join(n.IPs, "\n"))
	return sb.String()
}

// NewNetConfig creates a new configuration initialized from current network settings
func NewNetConfig() (*NetConfig, error) {
	n := &NetConfig{}
	hostname, err := os.Hostname()
	if err != nil {
		return n, err
	}
	n.Hostname = hostname
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return n, err
	}
	for _, addr := range addresses {
		n.IPs = append(n.IPs, addr.String())
	}
	return n, nil
}

// NewNetConfigFromFile loads the configuration from given file
func NewNetConfigFromFile(path string) (*NetConfig, error) {
	n := &NetConfig{}
	err := n.Load(path)
	return n, err
}

// IPCount returns the number of IPs in the configuration
func (n *NetConfig) IPCount() int {
	return len(n.IPs)
}

// Load the hostname and the IP list from given file
func (n *NetConfig) Load(path string) error {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return err
	}
	dec := json.NewDecoder(file)
	return dec.Decode(n)
}

// Save the hostname and the IP list to the given file
func (n *NetConfig) Save(path string) error {
	file, err := os.Create(path)
	defer file.Close()
	if err != nil {
		return err
	}
	enc := json.NewEncoder(file)
	return enc.Encode(n)
}

// IsChanged returns true if the two configurations differ
func (n *NetConfig) IsChanged(other *NetConfig) bool {
	if n.Hostname != other.Hostname {
		return true
	}
	if n.IPCount() != other.IPCount() {
		return true
	}
	for _, ip := range n.IPs {
		if !contains(other.IPs, ip) {
			return true
		}
	}
	return false
}

func contains(lines []string, s string) bool {
	for _, line := range lines {
		if line == s {
			return true
		}
	}
	return false
}
