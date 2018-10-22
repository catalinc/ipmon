package lib

import (
	"encoding/json"
	"net/smtp"
	"os"
	"strconv"
)

// MailConfig represents the static configuration for mailer
type MailConfig struct {
	From       string   `json:"from"`
	Password   string   `json:"password"`
	Recipients []string `json:"recipients"`
	Hostname   string   `json:"hostname"`
	Port       int      `json:"port"`
}

// NewMailConfig creates mail configuration from file
func NewMailConfig(path string) (*MailConfig, error) {
	m := &MailConfig{}
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return m, err
	}
	dec := json.NewDecoder(file)
	err = dec.Decode(m)
	return m, err
}

// SendMail sends the given text according to configuration
func SendMail(cfg *MailConfig, text string) error {
	auth := smtp.PlainAuth("", cfg.From, cfg.Password, cfg.Hostname)
	msg := []byte(text)
	return smtp.SendMail(cfg.Hostname+":"+strconv.Itoa(cfg.Port), auth, cfg.From, cfg.Recipients, msg)
}
