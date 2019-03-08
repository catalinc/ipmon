package main

import (
	"github.com/catalinc/ipmon/lib"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestGetNetConfig(t *testing.T) {
	n, err := lib.GetNetConfig()
	if err != nil {
		t.Errorf("Unable to read current network configuration: %v", err)
	}
	if n.IPCount() < 1 {
		t.Errorf("Got %d IPs", n.IPCount())
	}
	if n.Hostname == "" {
		t.Errorf("Got no hostname")
	}
}

func TestSaveLoadNetConfig(t *testing.T) {
	tmpFile, err := createTempFile()
	if err != nil {
		t.Errorf("Unable to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile)

	n := &lib.NetConfig{Hostname: "example", IPs: []string{"127.0.0.1", "192.168.0.105"}}
	err = n.Save(tmpFile)
	if err != nil {
		t.Errorf("Unable to save network configuration: %v", err)
	}

	saved, err := lib.LoadNetConfig(tmpFile)
	if err != nil {
		t.Errorf("Unable to read saved network configuration: %v", err)
	}
	if n.IsChanged(saved) {
		t.Errorf("Saved configuration is different: initial: %v saved: %v", n, saved)
	}
}

func TestNetConfigIsChanged(t *testing.T) {
	testData := []*lib.NetConfig{
		{"example", []string{"127.0.0.1", "192.168.0.3"}},
		{"example", []string{"127.0.0.1", "192.168.0.19", "10.0.0.1"}},
		{"example", []string{"127.0.0.1"}},
		{"example2", []string{"127.0.0.1", "192.168.0.19"}},
		{"example", []string{}},
		{"", []string{"127.0.0.1", "192.168.0.19"}},
		{},
	}

	n := &lib.NetConfig{Hostname: "example", IPs: []string{"127.0.0.1", "192.168.0.19"}}
	for _, nc := range testData {
		if !n.IsChanged(nc) {
			t.Errorf("%v should be seen as changed from %v", n, nc)
		}
	}
}

func TestLoadMailConfig(t *testing.T) {
	cfg, err := lib.LoadMailConfig("mail.json")
	if err != nil {
		t.Errorf("Unable to load mailer configuration from file: %v", err)
	}
	expectedCfg := &lib.MailConfig{
		From:           "user@example.net",
		To:             "someone@example.net",
		ServerHost:     "mx.example.net",
		ServerPort:     25,
		ServerUser:     "user@example.net",
		ServerPassword: "secret",
	}
	if !reflect.DeepEqual(cfg, expectedCfg) {
		t.Errorf("Expected %v got %v", expectedCfg, cfg)
	}
}

func TestContains(t *testing.T) {
	strings := []string{"aa", "bb", "cc"}
	in := "aa"
	out := "dd"
	if !lib.Contains(strings, in) {
		t.Errorf("%v should contain %v", strings, in)
	}
	if lib.Contains(strings, "dd") {
		t.Errorf("%v should not contain %v", strings, out)
	}
}

func TestDiffs(t *testing.T) {
	crtConf := &lib.NetConfig{Hostname: "spiderman", IPs: []string{"127.0.0.1", "192.168.0.17"}}
	prevConf := &lib.NetConfig{Hostname: "batman", IPs: []string{"127.0.0.1", "192.168.0.18", "10.0.0.1"}}

	diffs := lib.Diffs(crtConf, prevConf)
	expectedDiffs := "Hostname changed: batman -> spiderman\nIP count changed: 3 -> 2\nNew IP: 192.168.0.17\n"
	if diffs != expectedDiffs {
		t.Errorf("Expected %v got %v", expectedDiffs, diffs)
	}
}

func TestReport(t *testing.T) {
	crtConf := &lib.NetConfig{Hostname: "spiderman", IPs: []string{"127.0.0.1", "192.168.0.17"}}
	prevConf := &lib.NetConfig{Hostname: "batman", IPs: []string{"127.0.0.1", "192.168.0.18", "10.0.0.1"}}

	report := lib.Report(crtConf, prevConf)
	expectedReport := "Changes summary:\n" +
		"Hostname changed: batman -> spiderman\n" +
		"IP count changed: 3 -> 2\nNew IP: 192.168.0.17\n\n" +
		"Current configuration:\n" + crtConf.String()
	if report != expectedReport {
		t.Errorf("Expected %v got %v", expectedReport, report)
	}

	report = lib.Report(crtConf, nil)
	expectedReport = "Current configuration:\n" + crtConf.String()
	if report != expectedReport {
		t.Errorf("Expected %v got %v", expectedReport, report)
	}
}

func createTempFile() (string, error) {
	tmpFile, err := ioutil.TempFile("", "ipmon.test")
	if err != nil {
		return "", err
	}
	tmpFile.Close()
	return tmpFile.Name(), nil
}
