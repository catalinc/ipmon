package main

import (
	"github.com/catalinc/ipmon/lib"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestNewNetConfig(t *testing.T) {
	n, err := lib.NewNetConfig()
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

func TestNetConfigSaveAndLoad(t *testing.T) {
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

	saved, err := lib.NewNetConfigFromFile(tmpFile)
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

func TestNewMailerConfig(t *testing.T) {
	cfg, err := lib.NewMailConfig("mail.json")
	if err != nil {
		t.Errorf("Unable to load mailer configuration from file: %v", err)
	}
	expectedCfg := lib.MailConfig{
		From:       "user@example.net",
		Password:   "secret",
		Recipients: []string{"someone@example.net"},
		Hostname:   "mx.example.net",
		Port:       25,
	}
	if !reflect.DeepEqual(*cfg, expectedCfg) {
		t.Errorf("Expected %v got %v", expectedCfg, cfg)
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
