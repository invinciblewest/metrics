package senders

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
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
	r := resty.New().R()
	r.Method = http.MethodPost
	r.URL = fmt.Sprintf("%s/update/%s/%s/%s", s.serverAddr, mType, mName, mValue)
	r.Header.Set("Content-Type", "text/plain")
	resp, err := r.Send()
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return errors.New("wrong response code")
	}
	return nil
}
