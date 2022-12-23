package gateways

import (
	"bytes"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
)

type Gateway struct {
	client *http.Client

	log *log.Logger
}

func NewGateway(client http.Client, log *log.Logger) *Gateway {
	return &Gateway{
		client: &client,
		log:    log,
	}
}

// GetModelResult is an ABSOLUTE PISS SHIT, will remove when get my nvidia card later
func (g *Gateway) GetModelResult(image []byte) (res string, err error) {
	url := "http://localhost:33334/get_class"

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	fileWriter, err := bodyWriter.CreateFormFile("image", "test.jpg")
	if err != nil {
		return "", err
	}

	fileWriter.Write(image)
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(url, contentType, bodyBuf)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(respData[1 : len(respData)-1]), err
}
