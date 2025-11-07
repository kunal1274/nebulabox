package templates

import (
	"encoding/json"
	"fmt"
	"sync"
)

// StackTemplate represents a pre-defined multi-container stack template
type StackTemplate struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Category    string            `json:"category"` // web, database, microservices, etc.
	Containers  []ContainerConfig `json:"containers"`
	Networks    []NetworkConfig   `json:"networks,omitempty"`
	Volumes     []VolumeConfig    `json:"volumes,omitempty"`
	EnvVars     map[string]string `json:"envVars,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Tags        []string          `json:"tags,omitempty"`
	CreatedAt   string            `json:"createdAt,omitempty"`
}

// ContainerConfig defines a container in a stack template
type ContainerConfig struct {
	Name        string            `json:"name"`
	Image       string            `json:"image"`
	Ports       map[string]string `json:"ports,omitempty"`
	EnvVars     map[string]string `json:"envVars,omitempty"`
	Volumes     []string          `json:"volumes,omitempty"`
	Network     string            `json:"network,omitempty"`
	Service     string            `json:"service,omitempty"`
	DependsOn   []string          `json:"dependsOn,omitempty"` // Container names this depends on
	HealthCheck *HealthCheckConfig `json:"healthCheck,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Command     []string          `json:"command,omitempty"`
	Resources   *ResourceConfig   `json:"resources,omitempty"`
}

// NetworkConfig defines a network in a stack template
type NetworkConfig struct {
	Name    string            `json:"name"`
	Driver  string            `json:"driver,omitempty"` // bridge, overlay, etc.
	Labels  map[string]string `json:"labels,omitempty"`
	Options map[string]string `json:"options,omitempty"`
}

// VolumeConfig defines a volume in a stack template
type VolumeConfig struct {
	Name     string            `json:"name"`
	Driver   string            `json:"driver,omitempty"`
	Labels   map[string]string `json:"labels,omitempty"`
	Options  map[string]string `json:"options,omitempty"`
}

// HealthCheckConfig defines health check configuration
type HealthCheckConfig struct {
	Type            string   `json:"type"` // http, tcp, cmd
	HTTPPath        string   `json:"httpPath,omitempty"`
	HTTPPort        string   `json:"httpPort,omitempty"`
	TCPPort         string   `json:"tcpPort,omitempty"`
	Command         []string `json:"command,omitempty"`
	IntervalSeconds int      `json:"intervalSeconds,omitempty"`
	TimeoutSeconds  int      `json:"timeoutSeconds,omitempty"`
	Retries         int      `json:"retries,omitempty"`
	StartPeriodSec  int      `json:"startPeriodSec,omitempty"`
}

// ResourceConfig defines resource limits
type ResourceConfig struct {
	CPUShares    int64 `json:"cpuShares,omitempty"`
	MemoryLimit  int64 `json:"memoryLimit,omitempty"` // bytes
	MemorySwap   int64 `json:"memorySwap,omitempty"`
	CPUQuota     int64 `json:"cpuQuota,omitempty"`
	CPUPeriod    int64 `json:"cpuPeriod,omitempty"`
	PidsLimit    int64 `json:"pidsLimit,omitempty"`
	IOPSRead     int64 `json:"iopsRead,omitempty"`
	IOPSWrite    int64 `json:"iopsWrite,omitempty"`
}

// TemplateManager manages stack templates
type TemplateManager struct {
	templates map[string]*StackTemplate
	mu        sync.RWMutex
}

// NewTemplateManager creates a new template manager with default templates
func NewTemplateManager() *TemplateManager {
	tm := &TemplateManager{
		templates: make(map[string]*StackTemplate),
	}
	tm.initDefaultTemplates()
	return tm
}

// GetTemplate retrieves a template by ID
func (tm *TemplateManager) GetTemplate(id string) (*StackTemplate, bool) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	template, exists := tm.templates[id]
	return template, exists
}

// ListTemplates returns all templates, optionally filtered by category or tag
func (tm *TemplateManager) ListTemplates(category, tag string) []*StackTemplate {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	var result []*StackTemplate
	for _, template := range tm.templates {
		if category != "" && template.Category != category {
			continue
		}
		if tag != "" {
			found := false
			for _, t := range template.Tags {
				if t == tag {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		result = append(result, template)
	}
	return result
}

// SaveTemplate saves a custom template
func (tm *TemplateManager) SaveTemplate(template *StackTemplate) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if template.ID == "" {
		return fmt.Errorf("template ID is required")
	}
	if template.Name == "" {
		return fmt.Errorf("template name is required")
	}

	tm.templates[template.ID] = template
	return nil
}

// DeleteTemplate removes a custom template
func (tm *TemplateManager) DeleteTemplate(id string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if _, exists := tm.templates[id]; !exists {
		return fmt.Errorf("template not found: %s", id)
	}

	// Don't allow deleting default templates
	if id == "lamp" || id == "mean" || id == "lemp" || id == "wordpress" || id == "microservices" {
		return fmt.Errorf("cannot delete default template: %s", id)
	}

	delete(tm.templates, id)
	return nil
}

// initDefaultTemplates initializes common stack templates
func (tm *TemplateManager) initDefaultTemplates() {
	// LAMP Stack (Linux, Apache, MySQL, PHP)
	tm.templates["lamp"] = &StackTemplate{
		ID:          "lamp",
		Name:        "LAMP Stack",
		Description: "Linux, Apache, MySQL, PHP stack for web applications",
		Category:    "web",
		Tags:        []string{"php", "apache", "mysql", "web"},
		Containers: []ContainerConfig{
			{
				Name:  "apache",
				Image: "php:8.2-apache",
				Ports: map[string]string{"80": "8080"},
				Volumes: []string{"./www:/var/www/html"},
				DependsOn: []string{"mysql"},
			},
			{
				Name:  "mysql",
				Image: "mysql:8.0",
				Ports: map[string]string{"3306": "3306"},
				EnvVars: map[string]string{
					"MYSQL_ROOT_PASSWORD": "rootpass",
					"MYSQL_DATABASE":      "appdb",
					"MYSQL_USER":          "appuser",
					"MYSQL_PASSWORD":      "apppass",
				},
				HealthCheck: &HealthCheckConfig{
					Type:            "cmd",
					Command:         []string{"mysqladmin", "ping", "-h", "localhost"},
					IntervalSeconds: 10,
					TimeoutSeconds:  5,
					Retries:         3,
				},
			},
		},
	}

	// MEAN Stack (MongoDB, Express, Angular, Node.js)
	tm.templates["mean"] = &StackTemplate{
		ID:          "mean",
		Name:        "MEAN Stack",
		Description: "MongoDB, Express, Angular, Node.js full-stack JavaScript framework",
		Category:    "web",
		Tags:        []string{"node", "mongodb", "angular", "javascript"},
		Containers: []ContainerConfig{
			{
				Name:  "mongodb",
				Image: "mongo:7",
				Ports: map[string]string{"27017": "27017"},
				EnvVars: map[string]string{
					"MONGO_INITDB_ROOT_USERNAME": "admin",
					"MONGO_INITDB_ROOT_PASSWORD": "adminpass",
				},
				HealthCheck: &HealthCheckConfig{
					Type:            "tcp",
					TCPPort:         "27017",
					IntervalSeconds: 10,
					TimeoutSeconds:  5,
					Retries:         3,
				},
			},
			{
				Name:  "node-app",
				Image: "node:20-alpine",
				Ports: map[string]string{"3000": "3000"},
				DependsOn: []string{"mongodb"},
				HealthCheck: &HealthCheckConfig{
					Type:            "http",
					HTTPPath:        "/health",
					HTTPPort:        "3000",
					IntervalSeconds: 10,
					TimeoutSeconds:  5,
					Retries:         3,
				},
			},
		},
	}

	// LEMP Stack (Linux, Nginx, MySQL, PHP)
	tm.templates["lemp"] = &StackTemplate{
		ID:          "lemp",
		Name:        "LEMP Stack",
		Description: "Linux, Nginx, MySQL, PHP stack with Nginx reverse proxy",
		Category:    "web",
		Tags:        []string{"php", "nginx", "mysql", "web"},
		Containers: []ContainerConfig{
			{
				Name:  "nginx",
				Image: "nginx:alpine",
				Ports: map[string]string{"80": "8080"},
				DependsOn: []string{"php", "mysql"},
			},
			{
				Name:  "php",
				Image: "php:8.2-fpm",
				DependsOn: []string{"mysql"},
			},
			{
				Name:  "mysql",
				Image: "mysql:8.0",
				Ports: map[string]string{"3306": "3306"},
				EnvVars: map[string]string{
					"MYSQL_ROOT_PASSWORD": "rootpass",
					"MYSQL_DATABASE":      "appdb",
				},
			},
		},
	}

	// WordPress Stack
	tm.templates["wordpress"] = &StackTemplate{
		ID:          "wordpress",
		Name:        "WordPress",
		Description: "WordPress with MySQL database",
		Category:    "cms",
		Tags:        []string{"wordpress", "mysql", "php", "cms"},
		Containers: []ContainerConfig{
			{
				Name:  "wordpress",
				Image: "wordpress:latest",
				Ports: map[string]string{"80": "8080"},
				EnvVars: map[string]string{
					"WORDPRESS_DB_HOST":     "mysql",
					"WORDPRESS_DB_USER":      "wordpress",
					"WORDPRESS_DB_PASSWORD":  "wordpress",
					"WORDPRESS_DB_NAME":      "wordpress",
				},
				DependsOn: []string{"mysql"},
				HealthCheck: &HealthCheckConfig{
					Type:            "http",
					HTTPPath:        "/",
					HTTPPort:        "80",
					IntervalSeconds: 30,
					TimeoutSeconds:  10,
					Retries:         3,
				},
			},
			{
				Name:  "mysql",
				Image: "mysql:8.0",
				EnvVars: map[string]string{
					"MYSQL_ROOT_PASSWORD": "rootpass",
					"MYSQL_DATABASE":      "wordpress",
					"MYSQL_USER":          "wordpress",
					"MYSQL_PASSWORD":      "wordpress",
				},
				HealthCheck: &HealthCheckConfig{
					Type:            "cmd",
					Command:         []string{"mysqladmin", "ping", "-h", "localhost"},
					IntervalSeconds: 10,
					TimeoutSeconds:  5,
					Retries:         3,
				},
			},
		},
	}

	// Microservices Stack (API Gateway + Services + Database)
	tm.templates["microservices"] = &StackTemplate{
		ID:          "microservices",
		Name:        "Microservices Stack",
		Description: "API Gateway with multiple microservices and shared database",
		Category:    "microservices",
		Tags:        []string{"microservices", "api-gateway", "redis", "postgres"},
		Containers: []ContainerConfig{
			{
				Name:  "api-gateway",
				Image: "nginx:alpine",
				Ports: map[string]string{"80": "8080"},
				DependsOn: []string{"user-service", "product-service", "postgres", "redis"},
			},
			{
				Name:  "user-service",
				Image: "node:20-alpine",
				Ports: map[string]string{"3001": "3001"},
				DependsOn: []string{"postgres", "redis"},
				HealthCheck: &HealthCheckConfig{
					Type:            "http",
					HTTPPath:        "/health",
					HTTPPort:        "3001",
					IntervalSeconds: 10,
					TimeoutSeconds:  5,
					Retries:         3,
				},
			},
			{
				Name:  "product-service",
				Image: "node:20-alpine",
				Ports: map[string]string{"3002": "3002"},
				DependsOn: []string{"postgres", "redis"},
				HealthCheck: &HealthCheckConfig{
					Type:            "http",
					HTTPPath:        "/health",
					HTTPPort:        "3002",
					IntervalSeconds: 10,
					TimeoutSeconds:  5,
					Retries:         3,
				},
			},
			{
				Name:  "postgres",
				Image: "postgres:15-alpine",
				Ports: map[string]string{"5432": "5432"},
				EnvVars: map[string]string{
					"POSTGRES_USER":     "postgres",
					"POSTGRES_PASSWORD": "postgres",
					"POSTGRES_DB":       "appdb",
				},
				HealthCheck: &HealthCheckConfig{
					Type:            "cmd",
					Command:         []string{"pg_isready", "-U", "postgres"},
					IntervalSeconds: 10,
					TimeoutSeconds:  5,
					Retries:         3,
				},
			},
			{
				Name:  "redis",
				Image: "redis:7-alpine",
				Ports: map[string]string{"6379": "6379"},
				HealthCheck: &HealthCheckConfig{
					Type:            "cmd",
					Command:         []string{"redis-cli", "ping"},
					IntervalSeconds: 10,
					TimeoutSeconds:  5,
					Retries:         3,
				},
			},
		},
	}

	// Django + PostgreSQL Stack
	tm.templates["django-postgres"] = &StackTemplate{
		ID:          "django-postgres",
		Name:        "Django + PostgreSQL",
		Description: "Django web framework with PostgreSQL database",
		Category:    "web",
		Tags:        []string{"python", "django", "postgres", "web"},
		Containers: []ContainerConfig{
			{
				Name:  "django",
				Image: "python:3.11",
				Ports: map[string]string{"8000": "8000"},
				EnvVars: map[string]string{
					"DJANGO_DB_HOST":     "postgres",
					"DJANGO_DB_NAME":     "djangodb",
					"DJANGO_DB_USER":     "django",
					"DJANGO_DB_PASSWORD": "djangopass",
				},
				DependsOn: []string{"postgres"},
				HealthCheck: &HealthCheckConfig{
					Type:            "http",
					HTTPPath:        "/",
					HTTPPort:        "8000",
					IntervalSeconds: 30,
					TimeoutSeconds:  10,
					Retries:         3,
				},
			},
			{
				Name:  "postgres",
				Image: "postgres:15-alpine",
				Ports: map[string]string{"5432": "5432"},
				EnvVars: map[string]string{
					"POSTGRES_USER":     "django",
					"POSTGRES_PASSWORD": "djangopass",
					"POSTGRES_DB":       "djangodb",
				},
			},
		},
	}
}

// ToJSON converts a template to JSON string
func (tm *TemplateManager) ToJSON(templateID string) (string, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	template, exists := tm.templates[templateID]
	if !exists {
		return "", fmt.Errorf("template not found: %s", templateID)
	}

	data, err := json.MarshalIndent(template, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}

