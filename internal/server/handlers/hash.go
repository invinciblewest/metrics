package handlers

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/http"
	"strings"

	"github.com/invinciblewest/metrics/internal/logger"
	"go.uber.org/zap"
)

type responseRecorder struct {
	http.ResponseWriter
	body *strings.Builder
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// hashMiddleware создает middleware для проверки и добавления SHA256 хеша к запросам и ответам.
func hashMiddleware(hashKey string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if hashKey != "" {
				receivedHash := r.Header.Get("HashSHA256")
				if receivedHash != "" {
					decodedHash, err := base64.StdEncoding.DecodeString(receivedHash)
					if err != nil {
						logger.Log.Info("failed to decode hash", zap.Error(err))
						w.WriteHeader(http.StatusBadRequest)
						return
					}
					body, err := io.ReadAll(r.Body)
					if err != nil {
						logger.Log.Info("failed to read request body", zap.Error(err))
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					r.Body = io.NopCloser(bytes.NewBuffer(body))

					hash := hmac.New(sha256.New, []byte(hashKey))
					hash.Write(body)
					computedHash := hash.Sum(nil)
					if !hmac.Equal(computedHash, decodedHash) {
						logger.Log.Info(
							"hash mismatch",
							zap.String("expected", base64.StdEncoding.EncodeToString(computedHash)),
							zap.String("received", r.Header.Get("HashSHA256")),
						)
						w.WriteHeader(http.StatusBadRequest)
						return
					}
				}
			}

			rec := &responseRecorder{
				ResponseWriter: w,
				body:           new(strings.Builder),
			}
			next.ServeHTTP(rec, r)

			if hashKey != "" {
				hash := hmac.New(sha256.New, []byte(hashKey))
				hash.Write([]byte(rec.body.String()))
				computedHash := hash.Sum(nil)
				w.Header().Set("HashSHA256", base64.StdEncoding.EncodeToString(computedHash))
			}
		})
	}
}
