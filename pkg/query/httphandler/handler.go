package httphandler

import (
	"context"
	"net/http"

	"github.com/L4B0MB4/PRYVT/identification/pkg/query/httphandler/controller"
	"github.com/L4B0MB4/PRYVT/identification/pkg/query/httphandler/middleware"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type HttpHandler struct {
	httpServer     *http.Server
	router         *gin.Engine
	userController *controller.UserController
	authMiddleware *middleware.AuthMiddleware
}

func NewHttpHandler(c *controller.UserController, am *middleware.AuthMiddleware) *HttpHandler {
	r := gin.Default()
	srv := &http.Server{
		Addr:    "0.0.0.0" + ":" + "5517",
		Handler: r,
	}
	handler := &HttpHandler{
		router:         r,
		httpServer:     srv,
		userController: c,
		authMiddleware: am,
	}
	handler.RegisterRoutes()
	return handler
}

func (h *HttpHandler) RegisterRoutes() {
	h.router.POST("/:userId/token", h.userController.GetToken)
	h.router.Use(h.authMiddleware.AuthenticateMiddleware)
	{
		h.router.GET("/:userId", h.userController.GetUser)
		h.router.GET("/", h.userController.GetUsers)
	}
}

func (h *HttpHandler) Start() error {
	return h.httpServer.ListenAndServe()
}

func (h *HttpHandler) Stop() {
	err := h.httpServer.Shutdown(context.Background())
	if err != nil {
		log.Warn().Err(err).Msg("Error during reading response body")
	}
}
