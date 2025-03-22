package senders

import (
	"errors"
	"github.com/go-resty/resty/v2"
	"github.com/invinciblewest/metrics/internal/models"
	"net/http"
	"net/url"
)

type HTTPSender struct {
	serverAddr string
	client     *resty.Client
}

func NewHTTPSender(serverAddr string, client *http.Client) *HTTPSender {
	restyClient := resty.New()
	if client != nil {
		restyClient.SetTransport(client.Transport)
	}

	return &HTTPSender{
		serverAddr: serverAddr,
		client:     restyClient,
	}
}

func (s *HTTPSender) Send(metrics models.Metrics) error {
	path, err := url.JoinPath(s.serverAddr, "update")
	if err != nil {
		return err
	}

	resp, err := s.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(&metrics).
		Post(path)

	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return errors.New("wrong response code")
	}
	return nil
}
