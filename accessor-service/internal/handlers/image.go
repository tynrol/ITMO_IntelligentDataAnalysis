package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/tynrol/ITMO_IntelligentDataAnalysis/accessor-service/internal/gateways"
	"github.com/tynrol/ITMO_IntelligentDataAnalysis/accessor-service/internal/model/domain"
	"github.com/tynrol/ITMO_IntelligentDataAnalysis/accessor-service/internal/model/dto"
	"github.com/tynrol/ITMO_IntelligentDataAnalysis/accessor-service/internal/repositories"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
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
	// 1 in 10 chance to give honey image
	r := rand.Intn(10)
	if r == 1 {
		image, err := h.imageRepo.GetHoney(c)
		if err != nil {
			c.IndentedJSON(404, nil)
			return
		}
		h.log.Printf("Return honey image with id %s", image.ID)
		c.IndentedJSON(http.StatusOK, image)
		return
	}

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

	if request.Type == "WRONG" {
		c.IndentedJSON(200, nil)
		return
	}

	user, _ := h.userRepo.GetBySessionId(c, request.SessionId)
	h.log.Printf("USER ID: %s", user.SessionID)
	if !user.IsValid() {
		user = domain.User{
			SessionID: request.SessionId,
			Lied:      0,
		}
		err := h.userRepo.Create(c, user)
		if err != nil {
			c.IndentedJSON(409, err)
			return
		}
	}

	uploadPath, err := h.constructPath(request.Type, request.ImageId)
	if err != nil {
		c.IndentedJSON(412, err)
		return
	}

	image, _ := h.imageRepo.GetById(c, request.ImageId)
	if image.IsValid() {
		//then honey image
		h.log.Printf("CHECKING HONEY IMAGE %s", image.ID)
		if image.Type != request.Type {
			user.Lied = user.Lied + 1
			err := h.userRepo.Update(c, user)
			if err != nil {
				c.IndentedJSON(410, err)
				return
			}
			err = h.imageRepo.UpdatePathById(c, request.ImageId, uploadPath, request.Type, user.SessionID)
			if err != nil {
				c.IndentedJSON(414, err)
				return
			}

			c.IndentedJSON(http.StatusOK, nil)
			return
		}
	}

	fileBody, err := h.gateway.GetPhoto(request.ImageUrl)
	if err != nil {
		c.IndentedJSON(411, err)
		return
	}

	dir := filepath.Dir(uploadPath)
	if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
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
	h.log.Printf("TYPE %s, SESSION ID %s", request.Type, user.SessionID)
	err = h.imageRepo.UpdatePathById(c, request.ImageId, uploadPath, request.Type, user.SessionID)
	if err != nil {
		h.log.Printf("SOME SHITTY ERROR %s", err)
		c.IndentedJSON(414, err)
		return
	}

	c.IndentedJSON(http.StatusOK, nil)
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
	case "RAINY":
		path = h.datasetsPath + date + "/RAINY/" + imageId + ".jpeg"
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
