package dto

import (
	"github.com/tynrol/ITMO_IntelligentDataAnalysis/accessor-service/internal/model/domain"
	"time"
)

type Image struct {
	ID          string `json:"id"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	Description string `json:"description"`
	Urls        struct {
		Regular string `json:"regular"`
		Thumb   string `json:"thumb"`
	} `json:"urls"`
}

type PostImage struct {
	SessionId string `json:"session_id"`
	ImageId   string `json:"image_id"`
	ImageUrl  string `json:"image_url"`
	Type      string `json:"type"`
}

func ToDomain(image *Image) *domain.Image {
	return &domain.Image{
		ID:          image.ID,
		Width:       image.Width,
		Height:      image.Height,
		Description: image.Description,
		Url:         image.Urls.Thumb,
		CreatedAt:   time.Now(),
	}
}
