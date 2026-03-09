package types

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
)

// EncryptionKey represents the key encryption key (KEK)
// In production, this should be derived from a hardware security module (HSM)
// or a secure key management service (KMS)
type EncryptionKey struct {
	Key []byte `json:"key"`
}

// NewEncryptionKey creates a new 256-bit encryption key
func NewEncryptionKey() (*EncryptionKey, error) {
	key := make([]byte, 32) // AES-256
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, err
	}
	return &EncryptionKey{Key: key}, nil
}

// Encrypt encrypts plaintext using AES-256-GCM
func (ek *EncryptionKey) Encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(ek.Key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// Decrypt decrypts ciphertext using AES-256-GCM
func (ek *EncryptionKey) Decrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(ek.Key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

// HashKey creates a hash of the API key for storage verification
func HashKey(apiKey string) string {
	hash := sha256.Sum256([]byte(apiKey))
	return base64.URLEncoding.EncodeToString(hash[:])
}

// Zeroize securely wipes a byte slice
func Zeroize(data []byte) {
	for i := range data {
		data[i] = 0
	}
}

// SecureString represents a string that will be zeroized after use
type SecureString struct {
	data []byte
}

// NewSecureString creates a new secure string
func NewSecureString(s string) *SecureString {
	return &SecureString{data: []byte(s)}
}

// String returns the string content
func (ss *SecureString) String() string {
	return string(ss.data)
}

// Bytes returns the byte content
func (ss *SecureString) Bytes() []byte {
	return ss.data
}

// Zeroize wipes the secure string
func (ss *SecureString) Zeroize() {
	Zeroize(ss.data)
	ss.data = nil
}

// IsZeroized checks if the secure string has been zeroized
func (ss *SecureString) IsZeroized() bool {
	return ss.data == nil
}
