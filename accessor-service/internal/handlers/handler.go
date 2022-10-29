package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/tynrol/ITMO_IntelligentDataAnalysis/accessor-service/internal/gateways"
	"log"
	"net/http"
)

type Handler struct {
	gateway *gateways.Gateway

	log *log.Logger
}

func NewHandler(gateway *gateways.Gateway, logger *log.Logger) *Handler {
	return &Handler{
		gateway: gateway,
		log:     logger,
	}
}

func (h *Handler) Health(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, nil)
	return
}

func (h *Handler) GetRandPhoto(c *gin.Context) {
	image, err := h.gateway.GetRandomPhoto()
	if err != nil {
		c.IndentedJSON(http.StatusForbidden, nil)
	}

	//save id and url to db
	//transfer to service
	c.IndentedJSON(http.StatusOK, image)
	return
}

func (h *Handler) PostAccessorResult(c *gin.Context) {

}
