package handlers

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/tynrol/ITMO_IntelligentDataAnalysis/detection-service/internal/gateways"
	"io"
	"log"
	"net/http"
)

type Handler struct {
	gateway *gateways.Gateway

	log *log.Logger
}

func NewHandler(
	gateway *gateways.Gateway,
	logger *log.Logger,
) *Handler {
	return &Handler{
		gateway: gateway,
		log:     logger,
	}
}

// TODO: need to be done through probe
func (h *Handler) Health(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, nil)
	return
}

func (h *Handler) Detect(c *gin.Context) {
	file, _, err := c.Request.FormFile("image")
	if err != nil {
		h.log.Println(err)
		c.IndentedJSON(411, err)
		return
	}
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		c.IndentedJSON(412, err)
		return
	}
	//MOCK FOR TESTING
	result, err := h.gateway.GetModelResult(buf.Bytes())
	if err != nil {
		c.IndentedJSON(413, err)
		return
	}

	//result := "cloudy"
	if _, err := io.Copy(buf, file); err != nil {
		c.IndentedJSON(414, err)
		return
	}
	c.IndentedJSON(http.StatusOK, result)
	return
}
