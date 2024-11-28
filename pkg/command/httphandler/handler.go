package httphandler

import (
	"context"
	"net/http"

	"github.com/L4B0MB4/PRYVT/identification/pkg/command/httphandler/controller"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type HttpHandler struct {
	httpServer     *http.Server
	router         *gin.Engine
	userController *controller.UserController
}

func NewHttpHandler(c *controller.UserController) *HttpHandler {
	r := gin.Default()
	srv := &http.Server{
		Addr:    "0.0.0.0" + ":" + "5516",
		Handler: r,
	}
	handler := &HttpHandler{
		router:         r,
		httpServer:     srv,
		userController: c,
	}

	handler.RegisterRoutes()

	return handler
}

func (h *HttpHandler) RegisterRoutes() {
	h.router.POST("users/:userId/changeName", h.userController.ChangeDisplayName)
	h.router.POST("users/00000000-0000-0000-0000-000000000000/create", h.userController.CreateUser)
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
