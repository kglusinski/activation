package email

import (
	"fmt"
	"log"
	"net/smtp"
)

type Config struct {
	SmtpServer   string
	SmtpPort     int
	SenderEmail  string
	SenderPasswd string
}

type Client struct {
	config Config
}

func NewClient(config Config) *Client {
	return &Client{
		config: config,
	}
}

func (c *Client) Send(recipientEmail string, message string) error {
	auth := smtp.PlainAuth("", c.config.SenderEmail, c.config.SenderPasswd, c.config.SmtpServer)

	client, err := smtp.Dial(fmt.Sprintf("%s:%d", c.config.SmtpServer, c.config.SmtpPort))
	if err != nil {
		log.Println("Error connecting to the SMTP server:", err)
		return err
	}
	defer client.Close()

	if err := client.Auth(auth); err != nil {
		log.Println("Error authenticating with the SMTP server:", err)
		return err
	}

	if err := client.Mail(c.config.SenderEmail); err != nil {
		log.Println("Error setting sender:", err)
		return err
	}
	if err := client.Rcpt(recipientEmail); err != nil {
		log.Println("Error setting recipient:", err)
		return err
	}

	wc, err := client.Data()
	if err != nil {
		log.Println("Error opening data connection:", err)
		return err
	}
	_, err = wc.Write([]byte(message))
	if err != nil {
		log.Println("Error writing email body:", err)
		return err
	}
	err = wc.Close()
	if err != nil {
		log.Println("Error closing data connection:", err)
		return err
	}

	return nil
}
