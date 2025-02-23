package senders

import (
	"fmt"
	"net/http"
)

type HTTPSender struct {
	serverAddr string
	client     *http.Client
}

func NewHTTPSender(serverAddr string, client *http.Client) *HTTPSender {
	return &HTTPSender{
		serverAddr: serverAddr,
		client:     client,
	}
}

func (s *HTTPSender) Send(mType string, mName string, mValue string) error {
	url := fmt.Sprintf("%s/update/%s/%s/%s", s.serverAddr, mType, mName, mValue)
	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	r.Header.Set("Content-Type", "text/plain")

	res, err := s.client.Do(r)
	if err != nil {
		return err
	}
	defer func() {
		_ = res.Body.Close()
	}()

	return nil
}
