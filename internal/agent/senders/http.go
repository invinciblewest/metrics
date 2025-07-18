package senders

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/invinciblewest/metrics/pkg/encryption"

	"github.com/go-resty/resty/v2"
	"github.com/invinciblewest/metrics/internal/logger"
	"github.com/invinciblewest/metrics/internal/models"
	"go.uber.org/zap"
)

// HTTPSender отправляет метрики на сервер через HTTP с использованием RESTy клиента.
type HTTPSender struct {
	serverAddr string
	client     *resty.Client
	hashKey    string
	gzipPool   *sync.Pool
	bufPool    *sync.Pool
	cryptor    *encryption.Cryptor
}

// NewHTTPSender создает новый экземпляр HTTPSender с заданным адресом сервера, ключом хеширования и HTTP клиентом.
func NewHTTPSender(serverAddr string, hashKey string, client *http.Client, cryptor *encryption.Cryptor) *HTTPSender {
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
		hashKey:    hashKey,
		gzipPool: &sync.Pool{
			New: func() interface{} {
				return gzip.NewWriter(io.Discard)
			},
		},
		bufPool: &sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
		cryptor: cryptor,
	}
}

// SendMetric отправляет список метрик на сервер в формате JSON, сжимаемом с помощью gzip.
func (s *HTTPSender) SendMetric(ctx context.Context, metrics []models.Metric) error {
	path, err := url.JoinPath(s.serverAddr, "updates")
	if err != nil {
		return err
	}

	buf := s.bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer s.bufPool.Put(buf)

	gz := s.gzipPool.Get().(*gzip.Writer)
	gz.Reset(buf)
	defer func() {
		if err = gz.Close(); err != nil {
			logger.Log.Error("failed to close gzip writer", zap.Error(err))
			s.gzipPool.Put(gz)
		}
	}()

	if err = json.NewEncoder(gz).Encode(metrics); err != nil {
		return err
	}
	if err = gz.Close(); err != nil {
		return err
	}

	if s.cryptor != nil {
		var encryptedData []byte
		encryptedData, err = s.cryptor.Encrypt(buf.Bytes())
		if err != nil {
			return err
		}
		buf.Reset()
		if _, err = buf.Write(encryptedData); err != nil {
			return err
		}
	}

	realIP := getLocalIP()

	req := s.client.R().
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Accept-Encoding", "gzip").
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Real-IP", realIP).
		SetBody(buf.Bytes()).
		SetContext(ctx)

	if s.hashKey != "" {
		hash := hmac.New(sha256.New, []byte(s.hashKey))
		hash.Write(buf.Bytes())
		req.SetHeader("HashSHA256", base64.StdEncoding.EncodeToString(hash.Sum(nil)))
	}

	var resp *resty.Response
	resp, err = req.Post(path)

	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return errors.New("wrong response code")
	}
	return nil
}

// getLocalIP возвращает первый не-loopback IPv4 адрес хоста
func getLocalIP() string {
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, addr := range addresses {
		if IPNet, ok := addr.(*net.IPNet); ok && !IPNet.IP.IsLoopback() && IPNet.IP.To4() != nil {
			return IPNet.IP.String()
		}
	}
	return ""
}
