package composition

import (
	"fmt"
	"sync"
	"time"
)

// CompositionSpec defines how to create a new container by mixing elements from sources
type CompositionSpec struct {
	Name        string                     `json:"name"`
	Description string                     `json:"description,omitempty"`
	Sources     []SourceContainer          `json:"sources"`
	Overrides   *ContainerOverrides         `json:"overrides,omitempty"`
	Strategy    string                     `json:"strategy,omitempty"` // merge, override, priority
	CreatedAt   time.Time                  `json:"createdAt"`
}

// SourceContainer specifies what to extract from a source container
type SourceContainer struct {
	ContainerID   string                    `json:"containerId"`
	Elements      ContainerElements          `json:"elements"`
	Priority      int                       `json:"priority,omitempty"` // Higher priority wins in conflicts
	Description   string                    `json:"description,omitempty"`
}

// ContainerElements specifies which elements to extract from a source
type ContainerElements struct {
	Image         bool                      `json:"image,omitempty"`
	EnvVars       bool                      `json:"envVars,omitempty"`
	SelectedEnv   []string                  `json:"selectedEnv,omitempty"` // Specific env vars to include
	Ports         bool                      `json:"ports,omitempty"`
	SelectedPorts []string                  `json:"selectedPorts,omitempty"` // Specific ports to include
	Volumes       bool                      `json:"volumes,omitempty"`
	SelectedVolumes []string                `json:"selectedVolumes,omitempty"`
	Network       bool                      `json:"network,omitempty"`
	Service       bool                      `json:"service,omitempty"`
	HealthCheck   bool                      `json:"healthCheck,omitempty"`
	Labels        bool                      `json:"labels,omitempty"`
	SelectedLabels []string                 `json:"selectedLabels,omitempty"`
	Command       bool                      `json:"command,omitempty"`
	WorkingDir    bool                      `json:"workingDir,omitempty"`
	Resources     bool                      `json:"resources,omitempty"`
}

// ContainerOverrides allows manual overrides for the final container
type ContainerOverrides struct {
	Image       *string            `json:"image,omitempty"`
	EnvVars     map[string]string `json:"envVars,omitempty"`
	Ports       map[string]string `json:"ports,omitempty"`
	Volumes     map[string]string `json:"volumes,omitempty"`
	Network     *string           `json:"network,omitempty"`
	Service     *string           `json:"service,omitempty"`
	HealthCheck *HealthCheckOverride `json:"healthCheck,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Command     []string          `json:"command,omitempty"`
	WorkingDir  *string           `json:"workingDir,omitempty"`
}

// HealthCheckOverride allows overriding health check configuration
type HealthCheckOverride struct {
	Type            string   `json:"type,omitempty"`
	HTTPPath        string   `json:"httpPath,omitempty"`
	HTTPPort        string   `json:"httpPort,omitempty"`
	TCPPort         string   `json:"tcpPort,omitempty"`
	Command         []string `json:"command,omitempty"`
	IntervalSeconds int      `json:"intervalSeconds,omitempty"`
	TimeoutSeconds  int      `json:"timeoutSeconds,omitempty"`
	Retries         int      `json:"retries,omitempty"`
	StartPeriodSec  int      `json:"startPeriodSec,omitempty"`
}

// ComposedContainerSpec represents the final container specification after composition
type ComposedContainerSpec struct {
	Name        string                 `json:"name"`
	Image       string                 `json:"image"`
	EnvVars     map[string]string     `json:"envVars"`
	Ports       map[string]string     `json:"ports"`
	Volumes     map[string]string    `json:"volumes"`
	Network     string                 `json:"network,omitempty"`
	Service     string                 `json:"service,omitempty"`
	HealthCheck *HealthCheckOverride   `json:"healthCheck,omitempty"`
	Labels      map[string]string     `json:"labels"`
	Command     []string               `json:"command,omitempty"`
	WorkingDir  string                 `json:"workingDir,omitempty"`
	Conflicts   []ConflictResolution   `json:"conflicts,omitempty"`
	Sources     []string               `json:"sources"` // Container IDs used
}

// ConflictResolution describes how a conflict was resolved
type ConflictResolution struct {
	Type        string   `json:"type"`        // image, env, port, volume, network
	Source1     string   `json:"source1"`     // Container ID
	Source2     string   `json:"source2"`    // Container ID
	Value1      interface{} `json:"value1"`
	Value2      interface{} `json:"value2"`
	Resolution  string   `json:"resolution"` // merged, first, last, priority, custom
	FinalValue  interface{} `json:"finalValue"`
	Message     string   `json:"message,omitempty"`
}

// CompositionManager manages container compositions
type CompositionManager struct {
	specs map[string]*CompositionSpec
	mu    sync.RWMutex
}

// NewCompositionManager creates a new composition manager
func NewCompositionManager() *CompositionManager {
	return &CompositionManager{
		specs: make(map[string]*CompositionSpec),
	}
}

// SaveSpec saves a composition specification
func (cm *CompositionManager) SaveSpec(spec *CompositionSpec) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if spec.Name == "" {
		return fmt.Errorf("composition spec name is required")
	}

	spec.CreatedAt = time.Now()
	cm.specs[spec.Name] = spec
	return nil
}

// GetSpec retrieves a composition spec by name
func (cm *CompositionManager) GetSpec(name string) (*CompositionSpec, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	spec, exists := cm.specs[name]
	return spec, exists
}

// ListSpecs returns all composition specs
func (cm *CompositionManager) ListSpecs() []*CompositionSpec {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	specs := make([]*CompositionSpec, 0, len(cm.specs))
	for _, spec := range cm.specs {
		specs = append(specs, spec)
	}
	return specs
}

// DeleteSpec removes a composition spec
func (cm *CompositionManager) DeleteSpec(name string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if _, exists := cm.specs[name]; !exists {
		return fmt.Errorf("composition spec not found: %s", name)
	}

	delete(cm.specs, name)
	return nil
}

// ContainerSourceData represents the extracted data from a source container
type ContainerSourceData struct {
	ContainerID string
	Image       string
	EnvVars     map[string]string
	Ports       map[string]string
	Volumes     map[string]string
	Network     string
	Service     string
	HealthCheck *HealthCheckOverride
	Labels      map[string]string
	Command     []string
	WorkingDir  string
	Priority    int
}

// ComposeContainer creates a composed container specification from sources
func (cm *CompositionManager) ComposeContainer(
	spec *CompositionSpec,
	sourceDataGetter func(containerID string) (*ContainerSourceData, error),
) (*ComposedContainerSpec, error) {
	if spec == nil {
		return nil, fmt.Errorf("composition spec is required")
	}

	if len(spec.Sources) == 0 {
		return nil, fmt.Errorf("at least one source container is required")
	}

	composed := &ComposedContainerSpec{
		Name:      spec.Name,
		EnvVars:   make(map[string]string),
		Ports:     make(map[string]string),
		Volumes:   make(map[string]string),
		Labels:    make(map[string]string),
		Sources:   make([]string, 0),
		Conflicts: make([]ConflictResolution, 0),
	}

	var imageSources []*ContainerSourceData
	var healthCheckSources []*ContainerSourceData
	networkSources := make(map[string]*ContainerSourceData)
	serviceSources := make(map[string]*ContainerSourceData)

	// Collect source data
	sources := make([]*ContainerSourceData, 0, len(spec.Sources))
	for _, src := range spec.Sources {
		data, err := sourceDataGetter(src.ContainerID)
		if err != nil {
			return nil, fmt.Errorf("failed to get source data for %s: %w", src.ContainerID, err)
		}

		// Apply element filters
		data = filterSourceData(data, &src.Elements)
		data.Priority = src.Priority
		if data.Priority == 0 {
			data.Priority = 1 // Default priority
		}

		sources = append(sources, data)
		composed.Sources = append(composed.Sources, src.ContainerID)

		if src.Elements.Image && data.Image != "" {
			imageSources = append(imageSources, data)
		}
		if src.Elements.HealthCheck && data.HealthCheck != nil {
			healthCheckSources = append(healthCheckSources, data)
		}
		if src.Elements.Network && data.Network != "" {
			if existing, ok := networkSources[data.Network]; !ok || data.Priority > existing.Priority {
				networkSources[data.Network] = data
			}
		}
		if src.Elements.Service && data.Service != "" {
			if existing, ok := serviceSources[data.Service]; !ok || data.Priority > existing.Priority {
				serviceSources[data.Service] = data
			}
		}
	}

	// Resolve image (must have one, highest priority wins)
	if len(imageSources) > 0 {
		highestPriority := imageSources[0]
		for _, src := range imageSources[1:] {
			if src.Priority > highestPriority.Priority {
				highestPriority = src
			} else if src.Priority == highestPriority.Priority {
				// Same priority - check for conflict
				if src.Image != highestPriority.Image {
					composed.Conflicts = append(composed.Conflicts, ConflictResolution{
						Type:       "image",
						Source1:    highestPriority.ContainerID,
						Source2:    src.ContainerID,
						Value1:     highestPriority.Image,
						Value2:     src.Image,
						Resolution: "priority",
						FinalValue: highestPriority.Image,
						Message:    fmt.Sprintf("Image conflict: using %s from %s (priority %d)", highestPriority.Image, highestPriority.ContainerID, highestPriority.Priority),
					})
				}
			}
		}
		composed.Image = highestPriority.Image
	} else {
		return nil, fmt.Errorf("no image source specified")
	}

	// Merge environment variables
	for _, src := range sources {
		if src.EnvVars != nil {
			for k, v := range src.EnvVars {
				if existing, exists := composed.EnvVars[k]; exists && existing != v {
					// Conflict resolution based on priority
					if src.Priority > 1 { // If higher priority, override
						composed.EnvVars[k] = v
						composed.Conflicts = append(composed.Conflicts, ConflictResolution{
							Type:       "env",
							Source1:    "existing",
							Source2:    src.ContainerID,
							Value1:     existing,
							Value2:     v,
							Resolution: "priority",
							FinalValue: v,
							Message:    fmt.Sprintf("Env var %s conflict: using %s from %s", k, v, src.ContainerID),
						})
					}
				} else {
					composed.EnvVars[k] = v
				}
			}
		}
	}

	// Merge ports (handle conflicts)
	for _, src := range sources {
		if src.Ports != nil {
			for containerPort, hostPort := range src.Ports {
				if existing, exists := composed.Ports[containerPort]; exists && existing != hostPort {
					composed.Conflicts = append(composed.Conflicts, ConflictResolution{
						Type:       "port",
						Source1:    "existing",
						Source2:    src.ContainerID,
						Value1:     existing,
						Value2:     hostPort,
						Resolution: "priority",
						FinalValue: hostPort,
						Message:    fmt.Sprintf("Port %s conflict: using %s from %s", containerPort, hostPort, src.ContainerID),
					})
				}
				composed.Ports[containerPort] = hostPort
			}
		}
	}

	// Merge volumes
	for _, src := range sources {
		if src.Volumes != nil {
			for source, dest := range src.Volumes {
				if existing, exists := composed.Volumes[source]; exists && existing != dest {
					composed.Conflicts = append(composed.Conflicts, ConflictResolution{
						Type:       "volume",
						Source1:    "existing",
						Source2:    src.ContainerID,
						Value1:     existing,
						Value2:     dest,
						Resolution: "priority",
						FinalValue: dest,
					})
				}
				composed.Volumes[source] = dest
			}
		}
	}

	// Merge labels
	for _, src := range sources {
		if src.Labels != nil {
			for k, v := range src.Labels {
				composed.Labels[k] = v
			}
		}
	}

	// Resolve network (highest priority)
	if len(networkSources) > 0 {
		var selected *ContainerSourceData
		for _, src := range networkSources {
			if selected == nil || src.Priority > selected.Priority {
				selected = src
			}
		}
		composed.Network = selected.Network
	}

	// Resolve service (highest priority)
	if len(serviceSources) > 0 {
		var selected *ContainerSourceData
		for _, src := range serviceSources {
			if selected == nil || src.Priority > selected.Priority {
				selected = src
			}
		}
		composed.Service = selected.Service
	}

	// Resolve health check (highest priority)
	if len(healthCheckSources) > 0 {
		highestPriority := healthCheckSources[0]
		for _, src := range healthCheckSources[1:] {
			if src.Priority > highestPriority.Priority {
				highestPriority = src
			}
		}
		composed.HealthCheck = highestPriority.HealthCheck
	}

	// Apply overrides (manual overrides take precedence)
	if spec.Overrides != nil {
		if spec.Overrides.Image != nil {
			composed.Image = *spec.Overrides.Image
		}
		if spec.Overrides.EnvVars != nil {
			for k, v := range spec.Overrides.EnvVars {
				composed.EnvVars[k] = v
			}
		}
		if spec.Overrides.Ports != nil {
			for k, v := range spec.Overrides.Ports {
				composed.Ports[k] = v
			}
		}
		if spec.Overrides.Volumes != nil {
			for k, v := range spec.Overrides.Volumes {
				composed.Volumes[k] = v
			}
		}
		if spec.Overrides.Network != nil {
			composed.Network = *spec.Overrides.Network
		}
		if spec.Overrides.Service != nil {
			composed.Service = *spec.Overrides.Service
		}
		if spec.Overrides.HealthCheck != nil {
			composed.HealthCheck = spec.Overrides.HealthCheck
		}
		if spec.Overrides.Labels != nil {
			for k, v := range spec.Overrides.Labels {
				composed.Labels[k] = v
			}
		}
		if spec.Overrides.Command != nil {
			composed.Command = spec.Overrides.Command
		}
		if spec.Overrides.WorkingDir != nil {
			composed.WorkingDir = *spec.Overrides.WorkingDir
		}
	}

	return composed, nil
}

// filterSourceData applies element filters to source data
func filterSourceData(data *ContainerSourceData, elements *ContainerElements) *ContainerSourceData {
	filtered := &ContainerSourceData{
		ContainerID: data.ContainerID,
		Priority:    data.Priority,
	}

	if elements.Image {
		filtered.Image = data.Image
	}

	if elements.EnvVars {
		if len(elements.SelectedEnv) > 0 {
			filtered.EnvVars = make(map[string]string)
			for _, key := range elements.SelectedEnv {
				if val, exists := data.EnvVars[key]; exists {
					filtered.EnvVars[key] = val
				}
			}
		} else {
			filtered.EnvVars = data.EnvVars
		}
	}

	if elements.Ports {
		if len(elements.SelectedPorts) > 0 {
			filtered.Ports = make(map[string]string)
			for _, port := range elements.SelectedPorts {
				if val, exists := data.Ports[port]; exists {
					filtered.Ports[port] = val
				}
			}
		} else {
			filtered.Ports = data.Ports
		}
	}

	if elements.Volumes {
		if len(elements.SelectedVolumes) > 0 {
			filtered.Volumes = make(map[string]string)
			for _, vol := range elements.SelectedVolumes {
				if val, exists := data.Volumes[vol]; exists {
					filtered.Volumes[vol] = val
				}
			}
		} else {
			filtered.Volumes = data.Volumes
		}
	}

	if elements.Network {
		filtered.Network = data.Network
	}

	if elements.Service {
		filtered.Service = data.Service
	}

	if elements.HealthCheck {
		filtered.HealthCheck = data.HealthCheck
	}

	if elements.Labels {
		if len(elements.SelectedLabels) > 0 {
			filtered.Labels = make(map[string]string)
			for _, key := range elements.SelectedLabels {
				if val, exists := data.Labels[key]; exists {
					filtered.Labels[key] = val
				}
			}
		} else {
			filtered.Labels = data.Labels
		}
	}

	if elements.Command {
		filtered.Command = data.Command
	}

	if elements.WorkingDir {
		filtered.WorkingDir = data.WorkingDir
	}

	return filtered
}

