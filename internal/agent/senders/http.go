package senders

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
	"github.com/invinciblewest/metrics/internal/logger"
	"github.com/invinciblewest/metrics/internal/models"
	"net/http"
	"net/url"
	"time"
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

	retryDelays := []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}

	restyClient.
		SetRetryCount(3).
		SetRetryWaitTime(1 * time.Second).
		SetRetryMaxWaitTime(5 * time.Second).
		AddRetryCondition(
			func(r *resty.Response, err error) bool {
				return err != nil || r.StatusCode() >= http.StatusInternalServerError
			},
		).
		SetRetryAfter(func(client *resty.Client, response *resty.Response) (time.Duration, error) {
			attempt := response.Request.Attempt
			if attempt <= len(retryDelays) {
				return retryDelays[attempt-1], nil
			}
			return 0, nil
		}).
		AddRetryHook(func(r *resty.Response, err error) {
			logger.Log.Info("retrying request...")
		})

	return &HTTPSender{
		serverAddr: serverAddr,
		client:     restyClient,
	}
}

func (s *HTTPSender) SendMetric(ctx context.Context, metrics []models.Metric) error {
	path, err := url.JoinPath(s.serverAddr, "updates")
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if err = json.NewEncoder(gz).Encode(metrics); err != nil {
		return err
	}
	if err = gz.Close(); err != nil {
		return err
	}

	resp, err := s.client.R().
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Accept-Encoding", "gzip").
		SetHeader("Content-Type", "application/json").
		SetBody(buf.Bytes()).
		SetContext(ctx).
		Post(path)

	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return errors.New("wrong response code")
	}
	return nil
}
