package senders

import (
	"fmt"
	"net/http"
)

type HttpSender struct {
	serverAddr string
	client     *http.Client
}

func NewHttpSender(serverAddr string, client *http.Client) *HttpSender {
	return &HttpSender{
		serverAddr: serverAddr,
		client:     client,
	}
}

func (s *HttpSender) Send(mType string, mName string, mValue string) error {
	url := fmt.Sprintf("%s/update/%s/%s/%s", s.serverAddr, mType, mName, mValue)
	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	r.Header.Set("Content-Type", "text/plain")

	if _, err := s.client.Do(r); err != nil {
		return err
	}

	return nil
}
