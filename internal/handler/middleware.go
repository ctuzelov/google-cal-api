package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *ClientHandler) IsAuthMiddleware(c *gin.Context) {
	token, err := c.Cookie("token")

	savedToken := MustLoadToken()
	if err != nil || token != savedToken {
		c.JSON(http.StatusUnauthorized, errors.New("user is not unauthorized"))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	h.log.Info("user is authorized")

	c.Next()
}
