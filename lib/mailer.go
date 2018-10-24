package lib

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/mail"
	"net/smtp"
	"os"
	"strconv"
)

// MailConfig represents the static configuration for mailer
type MailConfig struct {
	From           string `json:"from"`
	To             string `json:"to"`
	ServerHost     string `json:"serverHost"`
	ServerPort     int    `json:"serverPort"`
	ServerUser     string `json:"serverUser"`
	ServerPassword string `json:"serverPassword"`
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
func SendMailSSL(cfg *MailConfig, subject string, body string) error {
	from := mail.Address{Address: cfg.From}
	to := mail.Address{Address: cfg.To}

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subject

	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Connect to the SMTP Server
	auth := smtp.PlainAuth("", cfg.ServerUser, cfg.ServerPassword, cfg.ServerHost)

	// TLS config
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         cfg.ServerHost,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	serverHostPort := cfg.ServerHost + ":" + strconv.Itoa(cfg.ServerPort)
	conn, err := tls.Dial("tcp", serverHostPort, tlsConfig)
	if err != nil {
		return err
	}

	c, err := smtp.NewClient(conn, cfg.ServerHost)
	if err != nil {
		log.Panic(err)
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		return err
	}

	// To && From
	if err = c.Mail(from.Address); err != nil {
		return err
	}

	if err = c.Rcpt(to.Address); err != nil {
		return err
	}

	// Data
	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	c.Quit()

	return nil
}
