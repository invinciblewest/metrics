package encryption

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDecryptBodyMiddleware(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Ошибка генерации ключа: %v", err)
	}
	publicKey := &privateKey.PublicKey

	cryptor := &Cryptor{
		privateKey: privateKey,
		publicKey:  publicKey,
	}

	original := []byte("тестовое сообщение")

	// Шифруем тело запроса
	encrypted, err := cryptor.Encrypt(original)
	if err != nil {
		t.Fatalf("Ошибка шифрования: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(encrypted))
	rec := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Ошибка чтения тела: %v", err)
		}
		if !bytes.Equal(body, original) {
			t.Errorf("Ожидалось: %s, получено: %s", original, body)
		}
	})

	middleware := DecryptBodyMiddleware(cryptor)
	middleware(handler).ServeHTTP(rec, req)
}
