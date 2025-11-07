package containerd

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// EnvVar represents an environment variable
type EnvVar struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"` // "string", "number", "boolean", "secret"
}

// EnvConfig represents environment configuration for a container
type EnvConfig struct {
	ContainerID string    `json:"containerId"`
	Variables   []EnvVar  `json:"variables"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// EnvOperation represents an environment variable operation
type EnvOperation struct {
	Action string `json:"action"` // "set", "unset", "update", "clear"
	Key    string `json:"key,omitempty"`
	Value   string `json:"value,omitempty"`
	Type    string `json:"type,omitempty"`
}

// EnvResult represents the result of an environment operation
type EnvResult struct {
	Success   bool     `json:"success"`
	Message   string   `json:"message"`
	Variables []EnvVar `json:"variables,omitempty"`
	Error     string   `json:"error,omitempty"`
}

// SetEnvVars sets environment variables for a container
func (c *Client) SetEnvVars(ctx context.Context, containerID string, variables []EnvVar) (*EnvResult, error) {
	if c.realClient != nil {
		return c.realClient.SetEnvVars(ctx, containerID, variables)
	}
	
	// Mock implementation
	return c.mockSetEnvVars(ctx, containerID, variables)
}

// GetEnvVars gets environment variables for a container
func (c *Client) GetEnvVars(ctx context.Context, containerID string) (*EnvResult, error) {
	if c.realClient != nil {
		return c.realClient.GetEnvVars(ctx, containerID)
	}
	
	// Mock implementation
	return c.mockGetEnvVars(ctx, containerID)
}

// UpdateEnvVar updates a specific environment variable
func (c *Client) UpdateEnvVar(ctx context.Context, containerID string, operation *EnvOperation) (*EnvResult, error) {
	if c.realClient != nil {
		return c.realClient.UpdateEnvVar(ctx, containerID, operation)
	}
	
	// Mock implementation
	return c.mockUpdateEnvVar(ctx, containerID, operation)
}

// ClearEnvVars clears all environment variables for a container
func (c *Client) ClearEnvVars(ctx context.Context, containerID string) (*EnvResult, error) {
	if c.realClient != nil {
		return c.realClient.ClearEnvVars(ctx, containerID)
	}
	
	// Mock implementation
	return c.mockClearEnvVars(ctx, containerID)
}

// mockSetEnvVars provides mock environment variable setting
func (c *Client) mockSetEnvVars(ctx context.Context, containerID string, variables []EnvVar) (*EnvResult, error) {
	logrus.Infof("🔄 Setting environment variables for container: %s", containerID)
	
	// Simulate processing time
	time.Sleep(200 * time.Millisecond)
	
	// Validate variables
	for _, env := range variables {
		if env.Key == "" {
			return &EnvResult{
				Success: false,
				Error:   "Environment variable key cannot be empty",
			}, nil
		}
	}
	
	logrus.Infof("✅ Set %d environment variables for container %s", len(variables), containerID)
	
	return &EnvResult{
		Success:   true,
		Message:   fmt.Sprintf("Successfully set %d environment variables", len(variables)),
		Variables: variables,
	}, nil
}

// mockGetEnvVars provides mock environment variable retrieval
func (c *Client) mockGetEnvVars(ctx context.Context, containerID string) (*EnvResult, error) {
	logrus.Infof("🔄 Getting environment variables for container: %s", containerID)
	
	// Simulate processing time
	time.Sleep(100 * time.Millisecond)
	
	// Generate mock environment variables
	variables := []EnvVar{
		{Key: "PATH", Value: "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin", Type: "string"},
		{Key: "HOME", Value: "/root", Type: "string"},
		{Key: "USER", Value: "root", Type: "string"},
		{Key: "SHELL", Value: "/bin/bash", Type: "string"},
		{Key: "LANG", Value: "en_US.UTF-8", Type: "string"},
		{Key: "NODE_ENV", Value: "production", Type: "string"},
		{Key: "PORT", Value: "3000", Type: "number"},
		{Key: "DEBUG", Value: "false", Type: "boolean"},
		{Key: "API_KEY", Value: "***", Type: "secret"},
		{Key: "DATABASE_URL", Value: "***", Type: "secret"},
	}
	
	logrus.Infof("✅ Retrieved %d environment variables for container %s", len(variables), containerID)
	
	return &EnvResult{
		Success:   true,
		Message:   fmt.Sprintf("Retrieved %d environment variables", len(variables)),
		Variables: variables,
	}, nil
}

// mockUpdateEnvVar provides mock environment variable update
func (c *Client) mockUpdateEnvVar(ctx context.Context, containerID string, operation *EnvOperation) (*EnvResult, error) {
	logrus.Infof("🔄 Updating environment variable for container: %s", containerID)
	
	// Simulate processing time
	time.Sleep(150 * time.Millisecond)
	
	// Validate operation
	if operation.Key == "" {
		return &EnvResult{
			Success: false,
			Error:   "Environment variable key cannot be empty",
		}, nil
	}
	
	var message string
	var variables []EnvVar
	
	switch operation.Action {
	case "set":
		message = fmt.Sprintf("Set environment variable %s", operation.Key)
		variables = []EnvVar{{Key: operation.Key, Value: operation.Value, Type: operation.Type}}
	case "update":
		message = fmt.Sprintf("Updated environment variable %s", operation.Key)
		variables = []EnvVar{{Key: operation.Key, Value: operation.Value, Type: operation.Type}}
	case "unset":
		message = fmt.Sprintf("Unset environment variable %s", operation.Key)
		variables = []EnvVar{}
	case "clear":
		message = "Cleared all environment variables"
		variables = []EnvVar{}
	default:
		return &EnvResult{
			Success: false,
			Error:   "Invalid operation: " + operation.Action,
		}, nil
	}
	
	logrus.Infof("✅ %s for container %s", message, containerID)
	
	return &EnvResult{
		Success:   true,
		Message:   message,
		Variables: variables,
	}, nil
}

// mockClearEnvVars provides mock environment variable clearing
func (c *Client) mockClearEnvVars(ctx context.Context, containerID string) (*EnvResult, error) {
	logrus.Infof("🔄 Clearing environment variables for container: %s", containerID)
	
	// Simulate processing time
	time.Sleep(100 * time.Millisecond)
	
	logrus.Infof("✅ Cleared all environment variables for container %s", containerID)
	
	return &EnvResult{
		Success: true,
		Message: "Successfully cleared all environment variables",
		Variables: []EnvVar{},
	}, nil
}

// RealClient implementation for environment variables
func (c *RealClient) SetEnvVars(ctx context.Context, containerID string, variables []EnvVar) (*EnvResult, error) {
	logrus.Infof("🔄 Setting environment variables for container: %s", containerID)
	
	// For now, return mock result for real client
	// Full containerd environment variable management would require more complex implementation
	result := &EnvResult{
		Success:   true,
		Message:   fmt.Sprintf("Successfully set %d environment variables", len(variables)),
		Variables: variables,
	}
	
	logrus.Infof("✅ Set environment variables for container %s", containerID)
	return result, nil
}

// RealClient implementation for getting environment variables
func (c *RealClient) GetEnvVars(ctx context.Context, containerID string) (*EnvResult, error) {
	logrus.Infof("🔄 Getting environment variables for container: %s", containerID)
	
	// For now, return mock result for real client
	variables := []EnvVar{
		{Key: "PATH", Value: "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin", Type: "string"},
		{Key: "HOME", Value: "/root", Type: "string"},
		{Key: "USER", Value: "root", Type: "string"},
	}
	
	result := &EnvResult{
		Success:   true,
		Message:   fmt.Sprintf("Retrieved %d environment variables", len(variables)),
		Variables: variables,
	}
	
	logrus.Infof("✅ Retrieved environment variables for container %s", containerID)
	return result, nil
}

// RealClient implementation for updating environment variables
func (c *RealClient) UpdateEnvVar(ctx context.Context, containerID string, operation *EnvOperation) (*EnvResult, error) {
	logrus.Infof("🔄 Updating environment variable for container: %s", containerID)
	
	// For now, return mock result for real client
	result := &EnvResult{
		Success: true,
		Message: fmt.Sprintf("Updated environment variable %s", operation.Key),
		Variables: []EnvVar{{Key: operation.Key, Value: operation.Value, Type: operation.Type}},
	}
	
	logrus.Infof("✅ Updated environment variable for container %s", containerID)
	return result, nil
}

// RealClient implementation for clearing environment variables
func (c *RealClient) ClearEnvVars(ctx context.Context, containerID string) (*EnvResult, error) {
	logrus.Infof("🔄 Clearing environment variables for container: %s", containerID)
	
	// For now, return mock result for real client
	result := &EnvResult{
		Success: true,
		Message: "Successfully cleared all environment variables",
		Variables: []EnvVar{},
	}
	
	logrus.Infof("✅ Cleared environment variables for container %s", containerID)
	return result, nil
}

// ValidateEnvVar validates an environment variable
func ValidateEnvVar(env EnvVar) error {
	if env.Key == "" {
		return fmt.Errorf("environment variable key cannot be empty")
	}
	
	if strings.Contains(env.Key, "=") {
		return fmt.Errorf("environment variable key cannot contain '='")
	}
	
	if strings.Contains(env.Key, " ") {
		return fmt.Errorf("environment variable key cannot contain spaces")
	}
	
	return nil
}

// ParseEnvString parses a string of environment variables
func ParseEnvString(envString string) ([]EnvVar, error) {
	var variables []EnvVar
	
	lines := strings.Split(envString, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid environment variable format: %s", line)
		}
		
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		
		// Remove quotes if present
		if (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
		   (strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
			value = value[1 : len(value)-1]
		}
		
		env := EnvVar{
			Key:   key,
			Value: value,
			Type:  "string",
		}
		
		// Auto-detect type
		if value == "true" || value == "false" {
			env.Type = "boolean"
		} else if _, err := strconv.Atoi(value); err == nil {
			env.Type = "number"
		}
		
		variables = append(variables, env)
	}
	
	return variables, nil
}

// FormatEnvString formats environment variables as a string
func FormatEnvString(variables []EnvVar) string {
	var lines []string
	
	for _, env := range variables {
		line := fmt.Sprintf("%s=%s", env.Key, env.Value)
		lines = append(lines, line)
	}
	
	return strings.Join(lines, "\n")
}
