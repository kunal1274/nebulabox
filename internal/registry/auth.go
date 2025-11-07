package registry

import (
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"
	"sync"
	"time"
)

// AuthConfig represents authentication configuration
type AuthConfig struct {
	// Token-based auth
	Tokens map[string]*TokenInfo `json:"tokens"`
	
	// Basic auth (username/password)
	Users map[string]*UserInfo `json:"users"`
	
	// Token expiration
	TokenExpiry time.Duration `json:"tokenExpiry"`
	
	mu sync.RWMutex
}

// TokenInfo contains token metadata
type TokenInfo struct {
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiresAt time.Time `json:"expiresAt"`
	Scopes    []string  `json:"scopes"`
}

// UserInfo contains user authentication information
type UserInfo struct {
	Username string   `json:"username"`
	Password string   `json:"password"` // Should be hashed in production
	Roles    []string `json:"roles"`
}

// NewAuthConfig creates a new authentication configuration
func NewAuthConfig() *AuthConfig {
	return &AuthConfig{
		Tokens:     make(map[string]*TokenInfo),
		Users:      make(map[string]*UserInfo),
		TokenExpiry: 24 * time.Hour,
	}
}

// AddUser adds a user to the auth configuration
func (ac *AuthConfig) AddUser(username, password string, roles []string) {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	
	ac.Users[username] = &UserInfo{
		Username: username,
		Password: password, // In production, hash this
		Roles:    roles,
	}
}

// AuthenticateUser authenticates a user with username/password
func (ac *AuthConfig) AuthenticateUser(username, password string) (*UserInfo, error) {
	ac.mu.RLock()
	defer ac.mu.RUnlock()
	
	user, ok := ac.Users[username]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}
	
	// Constant-time comparison
	if subtle.ConstantTimeCompare([]byte(user.Password), []byte(password)) != 1 {
		return nil, fmt.Errorf("invalid password")
	}
	
	return user, nil
}

// GenerateToken generates a new token for a user
func (ac *AuthConfig) GenerateToken(username string, scopes []string) string {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	
	token := generateSecureToken()
	
	ac.Tokens[token] = &TokenInfo{
		Username:  username,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(ac.TokenExpiry),
		Scopes:    scopes,
	}
	
	return token
}

// ValidateToken validates a token and returns user info
func (ac *AuthConfig) ValidateToken(token string) (*TokenInfo, error) {
	ac.mu.RLock()
	defer ac.mu.RUnlock()
	
	info, ok := ac.Tokens[token]
	if !ok {
		return nil, fmt.Errorf("invalid token")
	}
	
	if time.Now().After(info.ExpiresAt) {
		return nil, fmt.Errorf("token expired")
	}
	
	return info, nil
}

// RevokeToken revokes a token
func (ac *AuthConfig) RevokeToken(token string) bool {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	
	if _, ok := ac.Tokens[token]; ok {
		delete(ac.Tokens, token)
		return true
	}
	return false
}

// ParseBasicAuth parses HTTP Basic Authentication header
func ParseBasicAuth(authHeader string) (username, password string, ok bool) {
	if !strings.HasPrefix(authHeader, "Basic ") {
		return "", "", false
	}
	
	encoded := strings.TrimPrefix(authHeader, "Basic ")
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", "", false
	}
	
	parts := strings.SplitN(string(decoded), ":", 2)
	if len(parts) != 2 {
		return "", "", false
	}
	
	return parts[0], parts[1], true
}

// ParseBearerAuth parses HTTP Bearer token from Authorization header
func ParseBearerAuth(authHeader string) (token string, ok bool) {
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", false
	}
	
	token = strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	return token, len(token) > 0
}

// generateSecureToken generates a secure random token
func generateSecureToken() string {
	// In production, use crypto/rand for better security
	return fmt.Sprintf("nb-%d-%d", time.Now().UnixNano(), time.Now().Unix()%10000)
}

// HasScope checks if a token has a required scope
func (ti *TokenInfo) HasScope(scope string) bool {
	for _, s := range ti.Scopes {
		if s == scope || s == "*" {
			return true
		}
	}
	return false
}

// HasRole checks if a user has a required role
func (ui *UserInfo) HasRole(role string) bool {
	for _, r := range ui.Roles {
		if r == role || r == "admin" {
			return true
		}
	}
	return false
}

