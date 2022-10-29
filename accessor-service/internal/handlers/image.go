package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/tynrol/ITMO_IntelligentDataAnalysis/accessor-service/internal/gateways"
	"github.com/tynrol/ITMO_IntelligentDataAnalysis/accessor-service/internal/model/dto"
	"github.com/tynrol/ITMO_IntelligentDataAnalysis/accessor-service/internal/repositories"
	"log"
	"net/http"
	"time"
)

type Handler struct {
	gateway *gateways.Gateway
	repo    *repositories.Repo

	log *log.Logger
}

func NewHandler(gateway *gateways.Gateway, repo *repositories.Repo, logger *log.Logger) *Handler {
	return &Handler{
		gateway: gateway,
		repo:    repo,
		log:     logger,
	}
}

func (h *Handler) Health(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, nil)
	return
}

func (h *Handler) GetRandPhoto(c *gin.Context) {
	ctx := context.Background()

	image, err := h.gateway.GetRandomPhoto()
	if err != nil {
		c.IndentedJSON(404, nil)
	}

	err = h.repo.Create(ctx, *dto.DtoToDomain(image))
	if err != nil {
		c.IndentedJSON(500, nil)
	}

	c.IndentedJSON(http.StatusOK, image)
	return
}

func (h *Handler) PostPhoto(c *gin.Context) {
	ctx := context.Background()

	id := c.Query("id")
	weather := c.Query("weather")

	multipartFileHeader, err := c.FormFile("file")
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, nil)
	}

	date := time.Now().Format("2006-01-02")
	var uploadedPath string

	switch weather {
	case "SUNNY":
		uploadedPath = "/home/tynrol/Code/GolandProjects/ITMO_IntelligentDataAnalysis/accessor-service/datasets/dataset" + date + "/SUNNY/" + id + ".jpeg"
		break
	case "CLOUDY":
		uploadedPath = "/home/tynrol/Code/GolandProjects/ITMO_IntelligentDataAnalysis/accessor-service/datasets/dataset" + date + "/CLOUDY/" + id + ".jpeg"
		break
	case "RAIN":
		uploadedPath = "/home/tynrol/Code/GolandProjects/ITMO_IntelligentDataAnalysis/accessor-service/datasets/dataset" + date + "/RAIN/" + id + ".jpeg"
		break
	case "SUNRISE":
		uploadedPath = "/home/tynrol/Code/GolandProjects/ITMO_IntelligentDataAnalysis/accessor-service/datasets/dataset" + date + "/SUNRISE/" + id + ".jpeg"
		break
	case "SUNSET":
		uploadedPath = "/home/tynrol/Code/GolandProjects/ITMO_IntelligentDataAnalysis/accessor-service/datasets/dataset" + date + "/SUNRISE/" + id + ".jpeg"
		break
	default:
		c.IndentedJSON(http.StatusBadRequest, nil)
		return
	}
	h.log.Printf("Upload path at %s", uploadedPath)

	//transcation
	err = c.SaveUploadedFile(multipartFileHeader, uploadedPath)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}
	h.repo.UpdatePathById(ctx, id, uploadedPath)

	c.IndentedJSON(http.StatusOK, nil)
	return
}
