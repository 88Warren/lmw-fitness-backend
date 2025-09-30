package email

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

func SendEmail(fromAddress, toAddress, subject, body, replyToAddress, smtpPassword string) error {
	log.Printf("Preparing to send email from %s to %s with subject: %s", fromAddress, toAddress, subject)

	m := gomail.NewMessage()
	m.SetHeader("From", fromAddress)
	m.SetHeader("To", toAddress)
	m.SetHeader("Subject", subject)

	if replyToAddress != "" {
		m.SetHeader("Reply-To", replyToAddress)
	}

	m.SetBody("text/html", body)

	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		log.Printf("Error converting SMTP_PORT to int: %v", err)
		return fmt.Errorf("invalid SMTP_PORT configuration: %w", err)
	}

	log.Printf("Creating SMTP dialer with host: %s, port: %d, username: %s",
		os.Getenv("SMTP_HOST"), port, os.Getenv("SMTP_USERNAME"))

	d := gomail.NewDialer(
		os.Getenv("SMTP_HOST"),
		port,
		os.Getenv("SMTP_USERNAME"),
		smtpPassword,
	)

	// Configure TLS
	d.TLSConfig = &tls.Config{
		ServerName:         os.Getenv("SMTP_HOST"),
		InsecureSkipVerify: false, // Keep secure for production
	}

	log.Printf("Attempting to dial and send email...")
	if err := d.DialAndSend(m); err != nil {
		log.Printf("Send email failed with detailed error: %v", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("Email sent successfully to %s", toAddress)
	return nil
}
