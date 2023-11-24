package security

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"time"
)

type ActivationService struct {
	extUsers    ExternalUserProvider
	users       UserRepository
	tokens      ActivationTokenRepository
	emailClient EmailClient
}

func NewActivationService(extUsers ExternalUserProvider, users UserRepository, tokens ActivationTokenRepository, emailClient EmailClient) *ActivationService {
	return &ActivationService{
		extUsers:    extUsers,
		users:       users,
		tokens:      tokens,
		emailClient: emailClient,
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

type EmailClient interface {
	Send(recipientEmail string, message string) error
}

func (s *ActivationService) SendActivationMessage(token *Token) error {
	user, err := s.extUsers.GetUser(*token)
	if err != nil {
		return err
	}

	activationToken, err := generateRandomToken()
	if err != nil {
		log.Println("Error connecting to the SMTP server:", err)
	}

	at := ActivationToken{
		Token:  activationToken,
		UserID: user.AccountID,
		SentAt: time.Now(),
	}
	err = s.tokens.SaveToken(at)
	if err != nil {
		return err
	}

	err = s.sendActivationEmail(user.Email, activationToken)
	if err != nil {
		return err
	}

	s.users.Save(*user)

	log.Println("Email sent successfully to", user.Email)

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

func (s *ActivationService) sendActivationEmail(email string, token string) error {
	subject := "Aktywuj swoje konto"
	body := "This is a test email sent from Golang to MailHog. Go to http://localhost:8181/activate?token=" + token + " to activate your account."

	// Compose the email message
	message := "Subject: " + subject + "\r\n" +
		"\r\n" +
		body

	return s.emailClient.Send(email, message)
}

func generateRandomToken() (string, error) {
	length := 16

	numBytes := (length * 3) / 4 // 4 bytes encode to 3 characters in base64

	tokenBytes := make([]byte, numBytes)

	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}

	token := base64.URLEncoding.EncodeToString(tokenBytes)

	// Trim the token to the desired length
	if len(token) > length {
		token = token[:length]
	}

	return token, nil
}
