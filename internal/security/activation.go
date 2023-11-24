package security

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/smtp"
	"time"
)

type ActivationService struct {
	extUsers ExternalUserProvider
	users    UserRepository
	tokens   ActivationTokenRepository
}

func NewActivationService(extUsers ExternalUserProvider, users UserRepository, tokens ActivationTokenRepository) *ActivationService {
	return &ActivationService{
		extUsers: extUsers,
		users:    users,
		tokens:   tokens,
	}
}

type ActivationToken struct {
	Token  string
	UserID string
	SentAt time.Time
}

type ActivationTokenRepository interface {
	SaveToken(t ActivationToken) error
	FindActivationToken(token string) (*ActivationToken, error)
}

func (s *ActivationService) SendActivationMessage(token *Token) error {
	user, err := s.extUsers.GetUser(*token)
	if err != nil {
		return err
	}

	activationToken, err := generateRandomToken()
	if err != nil {
		fmt.Println("Error connecting to the SMTP server:", err)
	}

	// MailHog SMTP server address and port
	smtpServer := "localhost" // You may need to change this if MailHog is running on a different host
	smtpPort := 1025          // Default MailHog SMTP port

	// Sender's email and password (not required for MailHog)
	senderEmail := "your.email@example.com"
	senderPassword := "your_password"

	// Recipient's email address
	recipientEmail := user.Email

	at := ActivationToken{
		Token:  activationToken,
		UserID: user.AccountID,
		SentAt: time.Now(),
	}
	err = s.tokens.SaveToken(at)
	if err != nil {
		return err
	}

	// Email subject and body
	subject := "Aktywuj swoje konto"
	body := "This is a test email sent from Golang to MailHog. Go to http://localhost:8181/activate?token=" + activationToken + " to activate your account."

	// Compose the email message
	message := "Subject: " + subject + "\r\n" +
		"\r\n" +
		body

	// Set up the authentication credentials
	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpServer)

	// Connect to the MailHog SMTP server
	client, err := smtp.Dial(fmt.Sprintf("%s:%d", smtpServer, smtpPort))
	if err != nil {
		fmt.Println("Error connecting to the SMTP server:", err)
		return err
	}
	defer client.Close()

	// Authenticate with the server (not required for MailHog)
	if err := client.Auth(auth); err != nil {
		fmt.Println("Error authenticating with the SMTP server:", err)
		return err
	}

	// Set the sender and recipient
	if err := client.Mail(senderEmail); err != nil {
		fmt.Println("Error setting sender:", err)
		return err
	}
	if err := client.Rcpt(recipientEmail); err != nil {
		fmt.Println("Error setting recipient:", err)
		return err
	}

	// Send the email message
	wc, err := client.Data()
	if err != nil {
		fmt.Println("Error opening data connection:", err)
		return err
	}
	_, err = wc.Write([]byte(message))
	if err != nil {
		fmt.Println("Error writing email body:", err)
		return err
	}
	err = wc.Close()
	if err != nil {
		fmt.Println("Error closing data connection:", err)
		return err
	}

	s.users.Save(*user)

	fmt.Println("Email sent successfully to", recipientEmail)

	return nil
}

func (s *ActivationService) Activate(token string) error {
	at, err := s.tokens.FindActivationToken(token)
	if err != nil {
		return err
	}

	currentTime := time.Now()
	timeDifference := currentTime.Sub(at.SentAt)
	if timeDifference.Hours() > 1 {
		return err
	}

	user, err := s.users.GetUserByID(at.UserID)
	if err != nil {
		return err
	}

	user.EmailVerified = true
	s.users.Save(*user)

	return nil
}

func generateRandomToken() (string, error) {
	length := 16

	// Calculate the number of bytes needed to generate the token
	numBytes := (length * 3) / 4 // 4 bytes encode to 3 characters in base64

	// Create a byte slice to store the random data
	tokenBytes := make([]byte, numBytes)

	// Read random data into the byte slice
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}

	// Encode the random data to base64 to create the token
	token := base64.URLEncoding.EncodeToString(tokenBytes)

	// Trim the token to the desired length
	if len(token) > length {
		token = token[:length]
	}

	return token, nil
}
