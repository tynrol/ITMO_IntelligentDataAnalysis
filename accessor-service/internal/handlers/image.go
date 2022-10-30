package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/tynrol/ITMO_IntelligentDataAnalysis/accessor-service/internal/gateways"
	"github.com/tynrol/ITMO_IntelligentDataAnalysis/accessor-service/internal/model/domain"
	"github.com/tynrol/ITMO_IntelligentDataAnalysis/accessor-service/internal/model/dto"
	"github.com/tynrol/ITMO_IntelligentDataAnalysis/accessor-service/internal/repositories"
	"log"
	"net/http"
	"os"
	"time"
)

type Handler struct {
	gateway   *gateways.Gateway
	imageRepo *repositories.ImageRepo
	userRepo  *repositories.UserRepo

	datasetsPath string

	log *log.Logger
}

func NewHandler(
	gateway *gateways.Gateway,
	userRepo *repositories.UserRepo,
	imageRepo *repositories.ImageRepo,
	path string,
	logger *log.Logger,
) *Handler {
	return &Handler{
		gateway:      gateway,
		userRepo:     userRepo,
		imageRepo:    imageRepo,
		datasetsPath: path,
		log:          logger,
	}
}

// TODO: need to be done through probe
func (h *Handler) Health(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, nil)
	return
}

func (h *Handler) GetRandPhoto(c *gin.Context) {
	image, err := h.gateway.GetRandomPhoto()
	if err != nil {
		c.IndentedJSON(404, err)
		return
	}

	err = h.imageRepo.Create(c, *dto.ToDomain(image))
	if err != nil {
		c.IndentedJSON(500, err)
		return
	}

	c.IndentedJSON(http.StatusOK, image)
	return
}

func (h *Handler) PostPhoto(c *gin.Context) {
	var request dto.PostImage

	if err := c.BindJSON(&request); err != nil {
		c.IndentedJSON(400, "Cannot unmarshal req")
		return
	}

	//user, err := h.userRepo.GetBySessionId(c, request.SessionId);
	//if
	if image, err := h.imageRepo.GetById(c, request.ImageId); image.IsValid() && err != nil {
		newUser := domain.User{
			SessionID: request.SessionId,
			IsLying:   image.Type != request.Type,
		}
		err := h.userRepo.Create(c, newUser)
		if err != nil {
			c.IndentedJSON(410, err)
			return
		}
	}

	fileBody, err := h.gateway.GetPhoto(request.ImageUrl)
	if err != nil {
		c.IndentedJSON(411, err)
		return
	}

	uploadPath, err := h.constructPath(request.Type, request.ImageId)
	if err != nil {
		c.IndentedJSON(412, err)
		return
	}
	if uploadPath == "" {
		c.IndentedJSON(200, nil)
		return
	}

	f, err := os.Create(uploadPath)
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

	err = h.imageRepo.UpdatePathById(c, request.ImageId, uploadPath)
	if err != nil {
		c.IndentedJSON(413, err)
		return
	}

	c.IndentedJSON(http.StatusOK, nil)
	return
}

func (h *Handler) GetHoney(c *gin.Context) {
	image, err := h.imageRepo.GetHoney(c)
	if err != nil {
		c.IndentedJSON(404, err)
		return
	}

	if !image.IsValid() {
		randImage, err := h.gateway.GetRandomPhoto()
		if err != nil {
			c.IndentedJSON(404, err)
			return
		}
		image = *dto.ToDomain(randImage)
	}

	c.IndentedJSON(http.StatusOK, image)
	return
}

func (h *Handler) constructPath(weatherType string, imageId string) (path string, err error) {
	const op = "ImageHandler_constructPath"
	date := time.Now().Format("2006-01-02")

	switch weatherType {
	case "SUNNY":
		path = h.datasetsPath + date + "/SUNNY/" + imageId + ".jpeg"
		break
	case "CLOUDY":
		path = h.datasetsPath + date + "/CLOUDY/" + imageId + ".jpeg"
		break
	case "RAIN":
		path = h.datasetsPath + date + "/RAIN/" + imageId + ".jpeg"
		break
	case "SUNRISE":
		path = h.datasetsPath + date + "/SUNRISE/" + imageId + ".jpeg"
		break
	default:
		return "", errors.Wrap(errors.New("Incorrect weather type"), op)
	}
	h.log.Printf("Upload path at %s", path)

	return path, nil
}
