package main

import (
	"os"

	webserver "activation_flow/internal/http"
	"activation_flow/internal/security"
	"activation_flow/pkg/livechat"
	"github.com/labstack/echo/v4"
)

func main() {
	// Start the server
	go StartServer()

	// graceful shutdown
	quit := make(chan os.Signal)
	<-quit
}

func StartServer() {
	config := livechat.Config{
		ClientID:     "f9e9e9d9ea2a1696e031c9eb23f13f26",
		ClientSecret: "03868ed41f6243984fb7a234fd12fef9fdd0ee2e",
		RedirectUrl:  "http://localhost:8181",
	}

	client := livechat.NewClient(config)
	userRepository := security.NewInMemoryUserRepository()
	tokenRepository := security.NewInMemoryTokenRepository()

	validator := security.NewValidator(userRepository, client)
	srv := security.NewActivationService(client, userRepository, tokenRepository)
	handlers := webserver.NewHandler(client, validator, srv)

	e := echo.New()
	e.File("/redirect", "public/index.html")
	// e.GET("/test", func(c echo.Context) error {
	// 	code := c.QueryParam("code")
	// 	log.Println("code is:", code)
	//
	// 	token, err := client.GetToken(code)
	// 	if err != nil {
	// 		log.Println("error:", err)
	// 		return c.String(http.StatusOK, err.Error())
	// 	}
	//
	// 	user, err := client.GetUserInfo(token.AccessToken)
	// 	if err != nil {
	// 		log.Println("error:", err)
	// 		return c.String(http.StatusOK, err.Error())
	// 	}
	//
	// 	return c.String(http.StatusOK, user.Email)
	// })
	e.GET("/", handlers.Redirect)
	e.GET("/activate", handlers.Activation)
	e.File("/activation/sent", "public/activation_sent.html")
	e.GET("/dashboard", handlers.Dashboard)
	e.Logger.Fatal(e.Start(":8181"))
}
