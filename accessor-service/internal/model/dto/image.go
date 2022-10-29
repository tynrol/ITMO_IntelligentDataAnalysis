package dto

import (
	"github.com/tynrol/ITMO_IntelligentDataAnalysis/accessor-service/internal/model/domain"
	"time"
)

type Image struct {
	ID             string `json:"id"`
	Width          int    `json:"width"`
	Height         int    `json:"height"`
	Color          string `json:"color"`
	Description    string `json:"description"`
	AltDescription string `json:"alt_description"`
	Urls           struct {
		Raw     string `json:"raw"`
		Full    string `json:"full"`
		Regular string `json:"regular"`
		Small   string `json:"small"`
		Thumb   string `json:"thumb"`
		SmallS3 string `json:"small_s3"`
	} `json:"urls"`
	Links struct {
		Self             string `json:"self"`
		HTML             string `json:"html"`
		Download         string `json:"download"`
		DownloadLocation string `json:"download_location"`
	} `json:"links"`
}

func DtoToDomain(image *Image) *domain.Image {
	return &domain.Image{
		ID:          image.ID,
		Width:       image.Width,
		Height:      image.Height,
		Description: image.Description,
		Url:         image.Urls.Thumb,
		CreatedAt:   time.Now(),
	}

}
