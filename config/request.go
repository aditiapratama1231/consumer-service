package config

import (
	"bytes"
	"net/http"
)

type Request interface {
	Post(string, []byte) (*http.Request, error)
	Patch(string, []byte) (*http.Request, error)
	Delete(string) (*http.Request, error)
}

type request struct {
	BaseURL string
}

func NewRequest(baseUrl string) request {
	return request{
		BaseURL: baseUrl,
	}
}

func (r *request) Post(url string, body []byte) (*http.Request, error) {
	return http.NewRequest("POST", r.BaseURL+url, bytes.NewBuffer(body))
}

func (r *request) Patch(url string, body []byte) (*http.Request, error) {
	return http.NewRequest("PATCH", r.BaseURL+url, bytes.NewBuffer(body))
}

func (r *request) Delete(url string) (*http.Request, error) {
	return http.NewRequest("DELETE", url, nil)
}
