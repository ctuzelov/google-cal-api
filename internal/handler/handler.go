package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
)

type ClientHandler struct {
	client  *http.Client
	service *calendar.Service
	log     *slog.Logger
	config  *oauth2.Config
}

func New(config *oauth2.Config, logger *slog.Logger) *ClientHandler {
	return &ClientHandler{
		log:    logger,
		config: config,
	}
}

func (h *ClientHandler) InitRoutes() *gin.Engine {
	router := gin.Default()

	router.GET("/token", h.SetToken)
	router.GET("/ping", h.AuthLink)

	router.GET("/event/all", h.IsAuthMiddleware, h.GetEvents)
	router.GET("/event/:id", h.IsAuthMiddleware, h.GetEvent)
	router.POST("/event/create", h.IsAuthMiddleware, h.CreateEvent)
	router.DELETE("/event/:id", h.IsAuthMiddleware, h.DeleteEvent)
	router.GET("/event/zoom", h.IsAuthMiddleware, h.GetFirstZoomMeeting)

	router.GET("/logout", h.IsAuthMiddleware, h.Logout)

	return router
}
