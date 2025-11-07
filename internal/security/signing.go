package security

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"time"
)

// ImageSignature represents a cryptographic signature for an image
type ImageSignature struct {
	Image      string    `json:"image"`
	Tag        string    `json:"tag"`
	Digest     string    `json:"digest"`
	SignedBy   string    `json:"signedBy"`   // Username or key ID
	SignedAt   time.Time `json:"signedAt"`
	Signature  string    `json:"signature"`   // Base64-encoded signature
	PublicKey  string    `json:"publicKey"`  // Base64-encoded public key (PEM)
	Algorithm  string    `json:"algorithm"`  // RSA, ECDSA, etc.
	KeyID      string    `json:"keyId"`      // Key identifier
}

// SigningKey represents a signing key pair
type SigningKey struct {
	KeyID      string
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
	CreatedAt  time.Time
	CreatedBy  string
}

// KeyManager manages signing keys
type KeyManager struct {
	keys map[string]*SigningKey
}

// NewKeyManager creates a new key manager
func NewKeyManager() *KeyManager {
	return &KeyManager{
		keys: make(map[string]*SigningKey),
	}
}

// GenerateKey generates a new RSA key pair for signing
func (km *KeyManager) GenerateKey(keyID, createdBy string) (*SigningKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("failed to generate key: %w", err)
	}

	key := &SigningKey{
		KeyID:      keyID,
		PrivateKey: privateKey,
		PublicKey:  &privateKey.PublicKey,
		CreatedAt:  time.Now(),
		CreatedBy:  createdBy,
	}

	km.keys[keyID] = key
	return key, nil
}

// GetKey retrieves a key by ID
func (km *KeyManager) GetKey(keyID string) (*SigningKey, bool) {
	key, exists := km.keys[keyID]
	return key, exists
}

// ListKeys returns all key IDs
func (km *KeyManager) ListKeys() []string {
	ids := make([]string, 0, len(km.keys))
	for id := range km.keys {
		ids = append(ids, id)
	}
	return ids
}

// ExportPublicKey exports the public key in PEM format
func (key *SigningKey) ExportPublicKey() (string, error) {
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(key.PublicKey)
	if err != nil {
		return "", fmt.Errorf("failed to marshal public key: %w", err)
	}

	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	return string(publicKeyPEM), nil
}

// ExportPrivateKey exports the private key in PEM format (for backup only)
func (key *SigningKey) ExportPrivateKey() (string, error) {
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(key.PrivateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})
	return string(privateKeyPEM), nil
}

// SignImage signs an image manifest
func (km *KeyManager) SignImage(keyID, image, tag, digest, signedBy string) (*ImageSignature, error) {
	key, exists := km.keys[keyID]
	if !exists {
		return nil, fmt.Errorf("key not found: %s", keyID)
	}

	// Create signature payload
	payload := map[string]interface{}{
		"image":  image,
		"tag":    tag,
		"digest": digest,
		"signedBy": signedBy,
		"signedAt": time.Now().Unix(),
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Sign the payload
	hash := sha256.Sum256(payloadJSON)
	signature, err := rsa.SignPKCS1v15(rand.Reader, key.PrivateKey, crypto.SHA256, hash[:])
	if err != nil {
		return nil, fmt.Errorf("failed to sign: %w", err)
	}

	// Export public key
	publicKeyPEM, err := key.ExportPublicKey()
	if err != nil {
		return nil, fmt.Errorf("failed to export public key: %w", err)
	}

	sig := &ImageSignature{
		Image:     image,
		Tag:       tag,
		Digest:    digest,
		SignedBy:  signedBy,
		SignedAt:  time.Now(),
		Signature: base64.StdEncoding.EncodeToString(signature),
		PublicKey: base64.StdEncoding.EncodeToString([]byte(publicKeyPEM)),
		Algorithm: "RSA-SHA256",
		KeyID:     keyID,
	}

	return sig, nil
}

// VerifySignature verifies an image signature
func (km *KeyManager) VerifySignature(sig *ImageSignature, image, tag, digest string) (bool, error) {
	// Reconstruct the payload
	payload := map[string]interface{}{
		"image":  image,
		"tag":    tag,
		"digest": digest,
		"signedBy": sig.SignedBy,
		"signedAt": sig.SignedAt.Unix(),
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return false, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Decode signature
	signatureBytes, err := base64.StdEncoding.DecodeString(sig.Signature)
	if err != nil {
		return false, fmt.Errorf("invalid signature encoding: %w", err)
	}

	// Decode public key
	publicKeyPEMBytes, err := base64.StdEncoding.DecodeString(sig.PublicKey)
	if err != nil {
		return false, fmt.Errorf("invalid public key encoding: %w", err)
	}

	block, _ := pem.Decode(publicKeyPEMBytes)
	if block == nil {
		return false, fmt.Errorf("failed to decode PEM block")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false, fmt.Errorf("failed to parse public key: %w", err)
	}

	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return false, fmt.Errorf("not an RSA public key")
	}

	// Verify signature
	hash := sha256.Sum256(payloadJSON)
	err = rsa.VerifyPKCS1v15(rsaPublicKey, crypto.SHA256, hash[:], signatureBytes)
	if err != nil {
		return false, nil // Signature invalid
	}

	// Verify payload matches
	if sig.Image != image || sig.Tag != tag || sig.Digest != digest {
		return false, nil
	}

	return true, nil
}

