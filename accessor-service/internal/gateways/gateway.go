package gateways

import (
	"encoding/json"
	"github.com/tynrol/ITMO_IntelligentDataAnalysis/accessor-service/internal/model/dto"
	"io"
	"log"
	"net/http"
)

type Gateway struct {
	client *http.Client
	token  string

	log *log.Logger
}

func NewGateway(client http.Client, token string, log *log.Logger) *Gateway {
	return &Gateway{
		client: &client,
		token:  token,
		log:    log,
	}
}

func (g *Gateway) GetRandomPhoto() (image *dto.Image, err error) {
	url := "https://api.unsplash.com/photos/random"
	authHeader := "Client-ID " + g.token

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Set("Authorization", authHeader)

	q := req.URL.Query()
	q.Add("query", dto.RandomWeather())
	req.URL.RawQuery = q.Encode()
	g.log.Printf("Doing a req with url: %s", req.URL.String())

	res, _ := g.client.Do(req)

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return image, err
	}

	err = json.Unmarshal(body, &image)
	if err != nil {
		return image, err
	}

	return image, err
}

func (g *Gateway) GetPhotoPhoto(url string) (body []byte, err error) {
	req, _ := http.NewRequest("GET", url, nil)
	res, _ := g.client.Do(req)

	defer res.Body.Close()
	body, err = io.ReadAll(res.Body)
	if err != nil {
		return body, err
	}

	return body, err
}
