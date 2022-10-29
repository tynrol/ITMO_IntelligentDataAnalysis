package gateways

import (
	"encoding/json"
	"github.com/tynrol/ITMO_IntelligentDataAnalysis/accessor-service/internal/domain"
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

func (g *Gateway) GetRandomPhoto() (*domain.ImageResp, error) {
	var image = &domain.ImageResp{}

	url := "https://api.unsplash.com/photos/random"
	authHeader := "Client-ID " + g.token

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Set("Authorization", authHeader)

	q := req.URL.Query()
	q.Add("query", domain.RandomWeather())
	req.URL.RawQuery = q.Encode()
	g.log.Printf("Doing a req with url: %s", req.URL.String())

	res, _ := g.client.Do(req)

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, image)
	if err != nil {
		return nil, err
	}

	return image, nil
}
