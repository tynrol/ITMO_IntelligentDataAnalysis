package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/tynrol/ITMO_IntelligentDataAnalysis/accessor-service/internal/gateways"
	"github.com/tynrol/ITMO_IntelligentDataAnalysis/accessor-service/internal/model/dto"
	"github.com/tynrol/ITMO_IntelligentDataAnalysis/accessor-service/internal/repositories"
	"log"
	"net/http"
	"os"
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
		c.IndentedJSON(404, err)
		return
	}

	err = h.repo.Create(ctx, *dto.ToDomain(image))
	if err != nil {
		c.IndentedJSON(500, err)
		return
	}

	c.IndentedJSON(http.StatusOK, image)
	return
}

func (h *Handler) PostPhoto(c *gin.Context) {
	ctx := context.Background()

	var request dto.PostImage

	if err := c.BindJSON(&request); err != nil {
		c.IndentedJSON(410, "Wrong weather type")
		return
	}

	fileBody, err := h.gateway.GetPhotoPhoto(request.ImageUrl)
	if err != nil {
		c.IndentedJSON(411, err)
		return
	}

	date := time.Now().Format("2006-01-02")
	var uploadedPath string

	switch request.Type {
	case "SUNNY":
		uploadedPath = "./datasets/dataset" + date + "/SUNNY/" + request.ImageId + ".jpeg"
		break
	case "CLOUDY":
		uploadedPath = "./datasets/dataset" + date + "/CLOUDY/" + request.ImageId + ".jpeg"
		break
	case "RAIN":
		uploadedPath = "./datasets/dataset" + date + "/RAIN/" + request.ImageId + ".jpeg"
		break
	case "SUNRISE":
		uploadedPath = "./datasets/dataset" + date + "/SUNRISE/" + request.ImageId + ".jpeg"
		break
	case "SUNSET":
		uploadedPath = "./datasets/dataset" + date + "/SUNRISE/" + request.ImageId + ".jpeg"
		break
	case "WRONG":
		c.IndentedJSON(http.StatusOK, "Non relevant picture")
		return
	default:
		c.IndentedJSON(412, "Wrong weather type")
		return
	}
	h.log.Printf("Upload path at %s", uploadedPath)

	f, err := os.Create(uploadedPath)
	if err != nil {
		c.IndentedJSON(413, err)
		return
	}
	_, err = f.Write(fileBody)
	if err != nil {
		c.IndentedJSON(413, err)
		f.Close()
		return
	}
	err = f.Close()
	if err != nil {
		c.IndentedJSON(413, err)
		return
	}

	err = h.repo.UpdatePathById(ctx, request.ImageId, uploadedPath)
	if err != nil {
		c.IndentedJSON(413, err)
		return
	}

	c.IndentedJSON(http.StatusOK, nil)
	return
}
