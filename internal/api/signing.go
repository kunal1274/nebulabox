package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nebulabox/nebulabox/internal/security"
)

// SignImageRequest represents a request to sign an image
type SignImageRequest struct {
	Image    string `json:"image" binding:"required"`
	Tag      string `json:"tag" binding:"required"`
	Digest   string `json:"digest" binding:"required"`
	KeyID    string `json:"keyId" binding:"required"`
}

// VerifySignatureRequest represents a request to verify a signature
type VerifySignatureRequest struct {
	Image     string                `json:"image" binding:"required"`
	Tag       string                `json:"tag" binding:"required"`
	Digest    string                `json:"digest" binding:"required"`
	Signature security.ImageSignature `json:"signature" binding:"required"`
}

// GenerateKeyRequest represents a request to generate a signing key
type GenerateKeyRequest struct {
	KeyID string `json:"keyId" binding:"required"`
}

// KeyInfo represents information about a signing key
type KeyInfo struct {
	KeyID     string    `json:"keyId"`
	CreatedAt time.Time `json:"createdAt"`
	CreatedBy string    `json:"createdBy"`
	PublicKey string    `json:"publicKey"` // PEM format
}

// signImage handles POST /api/security/sign
func (s *Server) signImage(c *gin.Context) {
	var req SignImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Get current user from context (set by auth middleware)
	username := "system"
	if user, exists := c.Get("username"); exists {
		if u, ok := user.(string); ok {
			username = u
		}
	}

	// Sign the image
	signature, err := s.keyManager.SignImage(req.KeyID, req.Image, req.Tag, req.Digest, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to sign image",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, signature)
}

// verifySignature handles POST /api/security/verify
func (s *Server) verifySignature(c *gin.Context) {
	var req VerifySignatureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Verify the signature
	valid, err := s.keyManager.VerifySignature(&req.Signature, req.Image, req.Tag, req.Digest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to verify signature",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid": valid,
		"image": req.Image,
		"tag":   req.Tag,
		"digest": req.Digest,
	})
}

// generateKey handles POST /api/security/keys/generate
func (s *Server) generateKey(c *gin.Context) {
	var req GenerateKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Get current user
	username := "system"
	if user, exists := c.Get("username"); exists {
		if u, ok := user.(string); ok {
			username = u
		}
	}

	// Generate key
	key, err := s.keyManager.GenerateKey(req.KeyID, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate key",
			"details": err.Error(),
		})
		return
	}

	// Export public key
	publicKeyPEM, err := key.ExportPublicKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to export public key",
			"details": err.Error(),
		})
		return
	}

	keyInfo := KeyInfo{
		KeyID:     key.KeyID,
		CreatedAt: key.CreatedAt,
		CreatedBy: key.CreatedBy,
		PublicKey: publicKeyPEM,
	}

	c.JSON(http.StatusOK, keyInfo)
}

// listKeys handles GET /api/security/keys
func (s *Server) listKeys(c *gin.Context) {
	keyIDs := s.keyManager.ListKeys()
	keys := make([]KeyInfo, 0, len(keyIDs))

	for _, keyID := range keyIDs {
		key, exists := s.keyManager.GetKey(keyID)
		if !exists {
			continue
		}

		publicKeyPEM, err := key.ExportPublicKey()
		if err != nil {
			continue
		}

		keys = append(keys, KeyInfo{
			KeyID:     key.KeyID,
			CreatedAt: key.CreatedAt,
			CreatedBy: key.CreatedBy,
			PublicKey: publicKeyPEM,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"keys": keys,
		"count": len(keys),
	})
}

// getKey handles GET /api/security/keys/:keyId
func (s *Server) getKey(c *gin.Context) {
	keyID := c.Param("keyId")
	key, exists := s.keyManager.GetKey(keyID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Key not found",
		})
		return
	}

	publicKeyPEM, err := key.ExportPublicKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to export public key",
		})
		return
	}

	keyInfo := KeyInfo{
		KeyID:     key.KeyID,
		CreatedAt: key.CreatedAt,
		CreatedBy: key.CreatedBy,
		PublicKey: publicKeyPEM,
	}

	c.JSON(http.StatusOK, keyInfo)
}

