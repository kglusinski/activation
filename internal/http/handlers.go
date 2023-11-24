package http

import (
	"log"
	"net/http"

	"activation_flow/internal/security"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	provider  security.TokenProvider
	validator EmailVerificationValidator
	service   *security.ActivationService
}

func NewHandler(provider security.TokenProvider, validator EmailVerificationValidator, service *security.ActivationService) *Handler {
	return &Handler{provider, validator, service}
}

type EmailVerificationValidator interface {
	HasEmailVerified(t *security.Token) bool
}

func (h *Handler) Redirect(c echo.Context) error {
	// get code from url
	code := c.QueryParam("code")

	// get auth token from GetToken()
	accessToken, err := h.provider.GetToken(code)
	if err != nil {
		return c.String(http.StatusOK, err.Error())
	}

	if !h.validator.HasEmailVerified(accessToken) {
		h.service.SendActivationMessage(accessToken)

		c.File("public/activation_sent.html")
	}

	// if yes, redirect to dashboard
	return c.Redirect(http.StatusMovedPermanently, "/dashboard")
}

func (h *Handler) Activation(c echo.Context) error {
	token := c.QueryParam("token")
	err := h.service.Activate(token)
	if err != nil {
		return c.String(http.StatusOK, err.Error())
	}

	return c.Redirect(http.StatusTemporaryRedirect, "/dashboard")
}

func (h *Handler) Dashboard(c echo.Context) error {
	log.Println("dashboard")
	return c.File("public/dashboard.html")
}
