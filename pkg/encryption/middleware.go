package encryption

import (
	"bytes"
	"io"
	"net/http"

	"github.com/invinciblewest/metrics/internal/logger"
	"go.uber.org/zap"
)

func DecryptBodyMiddleware(cryptor *Cryptor) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			encryptedBody, err := io.ReadAll(r.Body)
			if err != nil {
				logger.Log.Error("failed to read request body", zap.Error(err))
				http.Error(w, "failed to read request body", http.StatusBadRequest)
				return
			}

			var decryptedBody []byte
			decryptedBody, err = cryptor.Decrypt(encryptedBody)
			if err != nil {
				logger.Log.Error("failed to decrypt request body", zap.Error(err))
				http.Error(w, "failed to decrypt request body", http.StatusBadRequest)
				return
			}

			r.Body = io.NopCloser(bytes.NewReader(decryptedBody))
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
