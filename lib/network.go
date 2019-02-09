package lib

import (
	"encoding/json"
	"net"
	"os"
	"strconv"
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
	sb.WriteString("Hostname: ")
	sb.WriteString(n.Hostname)
	sb.WriteString("\nIPs:\n")
	sb.WriteString(strings.Join(n.IPs, "\n"))
	return sb.String()
}

// GetCurrentNetConfig reads current network settings
func GetCurrentNetConfig() (*NetConfig, error) {
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

// LoadNetConfig loads network settings from file
func LoadNetConfig(path string) (*NetConfig, error) {
	n := &NetConfig{}
	err := n.Load(path)
	return n, err
}

// IPCount returns the number of IPs in the configuration
func (n *NetConfig) IPCount() int {
	return len(n.IPs)
}

// Load hostname and IP list from file
func (n *NetConfig) Load(path string) error {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return err
	}
	dec := json.NewDecoder(file)
	return dec.Decode(n)
}

// Save hostname and IP list to file
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

// Diffs compute a report with the differences between two configurations
func (n *NetConfig) Diffs(other *NetConfig) string {
	sb := strings.Builder{}
	if n.Hostname != other.Hostname {
		sb.WriteString("Hostname changed: ")
		sb.WriteString(other.Hostname)
		sb.WriteString(" -> ")
		sb.WriteString(n.Hostname)
		sb.WriteString("\n")
	}
	if n.IPCount() != other.IPCount() {
		sb.WriteString("IP count changed: ")
		sb.WriteString(strconv.Itoa(other.IPCount()))
		sb.WriteString(" -> ")
		sb.WriteString(strconv.Itoa(n.IPCount()))
		sb.WriteString("\n")
	}
	for _, ip := range n.IPs {
		if !contains(other.IPs, ip) {
			sb.WriteString("New IP: ")
			sb.WriteString(ip)
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

func contains(lines []string, s string) bool {
	for _, line := range lines {
		if line == s {
			return true
		}
	}
	return false
}
