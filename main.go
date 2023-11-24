package main

import (
	"log"
	"os"
	"strconv"

	webserver "activation_flow/internal/http"
	"activation_flow/internal/security"
	"activation_flow/pkg/email"
	"activation_flow/pkg/livechat"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	go StartServer()

	quit := make(chan os.Signal)
	<-quit
}

func StartServer() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	config := livechat.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectUrl:  os.Getenv("REDIRECT_URL"),
	}

	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	emailConfig := email.Config{
		SmtpServer:   os.Getenv("SMTP_SERVER"),
		SmtpPort:     port,
		SenderEmail:  os.Getenv("SENDER_EMAIL"),
		SenderPasswd: os.Getenv("SENDER_PASSWD"),
	}

	client := livechat.NewClient(config)
	emailClient := email.NewClient(emailConfig)
	userRepository := security.NewInMemoryUserRepository()
	tokenRepository := security.NewInMemoryTokenRepository()

	validator := security.NewValidator(userRepository, client)
	srv := security.NewActivationService(client, userRepository, tokenRepository, emailClient)
	handlers := webserver.NewHandler(client, validator, srv)

	e := echo.New()
	e.GET("/", handlers.Redirect)
	e.GET("/activate", handlers.Activation)
	e.GET("/dashboard", handlers.Dashboard)

	// todo: configurable port
	e.Logger.Fatal(e.Start(":8181"))
}
