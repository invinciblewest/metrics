package encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
)

func TestCryptor_EncryptDecrypt(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Ошибка генерации ключа: %v", err)
	}
	publicKey := &privateKey.PublicKey

	cryptor := &Cryptor{
		privateKey: privateKey,
		publicKey:  publicKey,
	}

	original := []byte("секретное сообщение")

	var encrypted []byte
	encrypted, err = cryptor.Encrypt(original)
	if err != nil {
		t.Fatalf("Ошибка шифрования: %v", err)
	}

	decrypted, err := cryptor.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Ошибка дешифрования: %v", err)
	}

	if string(decrypted) != string(original) {
		t.Errorf("Ожидалось: %s, получено: %s", original, decrypted)
	}
}
