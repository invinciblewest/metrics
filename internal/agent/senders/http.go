package senders

import (
	"errors"
	"github.com/go-resty/resty/v2"
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

func (s *HTTPSender) Send(mType string, mName string, mValue string) error {
	path, err := url.JoinPath(s.serverAddr, "update", mType, mName, mValue)
	if err != nil {
		return err
	}

	resp, err := s.client.R().
		SetHeader("Content-Type", "text/plain").
		Post(path)

	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return errors.New("wrong response code")
	}
	return nil
}
