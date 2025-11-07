package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// RegistryClient handles communication with the Nebula Registry
type RegistryClient struct {
	baseURL string
	token   string
	client  *http.Client
}

// NewRegistryClient creates a new registry client
func NewRegistryClient(baseURL string) *RegistryClient {
	return &RegistryClient{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		client:  http.DefaultClient,
	}
}

// SetToken sets the authentication token
func (rc *RegistryClient) SetToken(token string) {
	rc.token = token
}

// Login authenticates with the registry and sets the token
func (rc *RegistryClient) Login(username, password string) (string, error) {
	loginURL := rc.baseURL + "/auth/login"
	
	body := map[string]string{
		"username": username,
		"password": password,
	}
	jsonBody, _ := json.Marshal(body)
	
	req, err := http.NewRequest("POST", loginURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := rc.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("registry unavailable: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("login failed: %s", string(bodyBytes))
	}
	
	var result struct {
		Token     string `json:"token"`
		TokenType string `json:"token_type"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}
	
	rc.token = result.Token
	return result.Token, nil
}

// ListRepositories lists all repositories in the registry
func (rc *RegistryClient) ListRepositories() ([]string, error) {
	apiURL := rc.baseURL + "/api/registry/repositories"
	
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	
	if rc.token != "" {
		req.Header.Set("Authorization", "Bearer "+rc.token)
	}
	
	resp, err := rc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("registry unavailable: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list repositories: status %d", resp.StatusCode)
	}
	
	var result struct {
		Repositories []string `json:"repositories"`
		Count        int      `json:"count"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	return result.Repositories, nil
}

// ListVersions lists all versions/tags for a repository
func (rc *RegistryClient) ListVersions(repo string) ([]map[string]interface{}, error) {
	apiURL := rc.baseURL + "/api/registry/repositories/" + url.PathEscape(repo) + "/versions"
	
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	
	if rc.token != "" {
		req.Header.Set("Authorization", "Bearer "+rc.token)
	}
	
	resp, err := rc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("registry unavailable: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list versions: status %d", resp.StatusCode)
	}
	
	var result struct {
		Repository string                   `json:"repository"`
		Versions   []map[string]interface{} `json:"versions"`
		Count      int                      `json:"count"`
		Latest     string                   `json:"latest"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	return result.Versions, nil
}

// GetRepositorySummary gets a summary of repository information
func (rc *RegistryClient) GetRepositorySummary(repo string) (map[string]interface{}, error) {
	apiURL := rc.baseURL + "/api/registry/repositories/" + url.PathEscape(repo) + "/summary"
	
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	
	if rc.token != "" {
		req.Header.Set("Authorization", "Bearer "+rc.token)
	}
	
	resp, err := rc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("registry unavailable: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get summary: status %d", resp.StatusCode)
	}
	
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	return result, nil
}

// makeRequest makes an authenticated request to the registry
func (rc *RegistryClient) makeRequest(method, path string, body io.Reader, headers map[string]string) (*http.Response, error) {
	fullURL := rc.baseURL + path
	
	req, err := http.NewRequest(method, fullURL, body)
	if err != nil {
		return nil, err
	}
	
	if rc.token != "" {
		req.Header.Set("Authorization", "Bearer "+rc.token)
	}
	
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	
	return rc.client.Do(req)
}

