package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"magento-consumer-service/domain"

	"github.com/imroc/req"
	"github.com/joho/godotenv"
)

type Request interface {
	Send(string, string, []byte) (*req.Resp, error)
	GetToken() error
	Post(string, []byte) (*req.Resp, error)
	Patch(string, []byte) (*req.Resp, error)
	Delete(string) (*req.Resp, error)
}

type request struct {
	BaseURL string
	Token   string
}

func NewRequest(baseUrl string) Request {
	return &request{
		BaseURL: baseUrl,
	}
}

func (r *request) GetToken() error {
	e := godotenv.Load() //Load .env file
	if e != nil {
		fmt.Print(e)
	}

	creds := &domain.Credential{
		Username: os.Getenv("MAGENTO_USERNAME"),
		Password: os.Getenv("MAGENTO_PASSWORD"),
	}

	credentials, err := json.Marshal(creds)
	if err != nil {
		log.Println("Error encoding credentials : " + err.Error())
		return err
	}

	req, err := r.Send("TOKEN", "/integration/admin/token", credentials)
	if err != nil {
		return err
	}
	type tokenRaw string

	var getToken tokenRaw
	req.ToJSON(&getToken)

	_ = request{
		Token: string(getToken),
	}

	return nil
}

func (r *request) Send(method string, url string, body []byte) (*req.Resp, error) {
	header := req.Header{
		"Content-Type": "application/json",
	}

	tokenHeader := req.Header{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + r.Token,
	}

	switch method {
	case "TOKEN":
		reqs, err := req.Post(r.BaseURL+url, header, bytes.NewBuffer(body))
		return reqs, err
	case "GET":
		reqs, err := req.Get(r.BaseURL+url, header, bytes.NewBuffer(body))
		return reqs, err
	case "POST":
		reqs, err := req.Post(r.BaseURL+url, tokenHeader, bytes.NewBuffer(body))
		return reqs, err
	case "PATCH":
		reqs, err := req.Patch(r.BaseURL+url, tokenHeader, bytes.NewBuffer(body))
		return reqs, err
	case "DELETE":
		reqs, err := req.Delete(r.BaseURL+url, tokenHeader)
		return reqs, err
	}

	return nil, nil
}

func (r *request) Post(url string, body []byte) (*req.Resp, error) {
	return r.Send("POST", url, body)
}

func (r *request) Patch(url string, body []byte) (*req.Resp, error) {
	return r.Send("PATCH", url, body)
}

func (r *request) Delete(url string) (*req.Resp, error) {
	return r.Send("DELETE", url, []byte(""))
}
