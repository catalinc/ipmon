package lib

import (
	"strconv"
	"strings"
)

// Diffs compute a report with the differences between two configurations
func Diffs(crt *NetConfig, prev *NetConfig) string {
	sb := strings.Builder{}
	if crt.Hostname != prev.Hostname {
		sb.WriteString("Hostname changed: ")
		sb.WriteString(prev.Hostname)
		sb.WriteString(" -> ")
		sb.WriteString(crt.Hostname)
		sb.WriteString("\n")
	}
	if crt.IPCount() != prev.IPCount() {
		sb.WriteString("IP count changed: ")
		sb.WriteString(strconv.Itoa(prev.IPCount()))
		sb.WriteString(" -> ")
		sb.WriteString(strconv.Itoa(crt.IPCount()))
		sb.WriteString("\n")
	}
	for _, ip := range crt.IPs {
		if !Contains(prev.IPs, ip) {
			sb.WriteString("New IP: ")
			sb.WriteString(ip)
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

// Report builds the summary for network configuration changes
func Report(crt *NetConfig, prev *NetConfig) string {
	sb := strings.Builder{}
	if prev != nil {
		sb.WriteString("Changes summary:\n")
		sb.WriteString(Diffs(crt, prev))
		sb.WriteString("\n")
	}
	sb.WriteString("Current configuration:\n")
	sb.WriteString(crt.String())
	return sb.String()
}
