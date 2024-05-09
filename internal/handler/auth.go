package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ctuzelov/google-cal-api/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func (h *ClientHandler) AuthLink(c *gin.Context) {
	authURL := h.config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	c.JSON(http.StatusOK, fmt.Sprintf("Go to the following link in your browser: %v", authURL))
}

func (h *ClientHandler) SetToken(c *gin.Context) {
	code := c.Query("code")

	// TODO: validate FormValue("state") to protect against CSRF attacks
	tokens, err := utils.GenerateToken(code, h.config)
	if err != nil {
		h.log.Error("failed to generate token", err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	h.client = h.config.Client(context.TODO(), tokens)
	h.service, err = calendar.NewService(context.TODO(), option.WithHTTPClient(h.client))

	if err != nil {
		h.log.Error("failed to create calendar service", err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	MustSaveToken(tokens.AccessToken)

	c.JSON(http.StatusOK, tokens)
}

func (h *ClientHandler) Logout(c *gin.Context) {
	MustChangeToken()
}
