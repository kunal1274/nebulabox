package buildspec

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// BuildSpec represents a NebulaBox build specification
type BuildSpec struct {
	Version  string                 `yaml:"version" json:"version"`
	Name     string                 `yaml:"name" json:"name"`
	Tag      string                 `yaml:"tag" json:"tag"`
	Base     BaseImage              `yaml:"base" json:"base"`
	Steps    []BuildStep            `yaml:"steps" json:"steps"`
	Env      map[string]string      `yaml:"env,omitempty" json:"env,omitempty"`
	Workdir  string                 `yaml:"workdir,omitempty" json:"workdir,omitempty"`
	Expose   []int                  `yaml:"expose,omitempty" json:"expose,omitempty"`
	Labels   map[string]string      `yaml:"labels,omitempty" json:"labels,omitempty"`
	Health   *HealthCheck            `yaml:"health,omitempty" json:"health,omitempty"`
	User     string                 `yaml:"user,omitempty" json:"user,omitempty"`
	Metadata map[string]interface{} `yaml:"metadata,omitempty" json:"metadata,omitempty"`
}

// BaseImage defines the base image for the build
type BaseImage struct {
	Image string `yaml:"image" json:"image"`
	Tag   string `yaml:"tag,omitempty" json:"tag,omitempty"`
}

// BuildStep represents a single build step
type BuildStep struct {
	Type    string            `yaml:"type" json:"type"` // run, copy, add, cmd, arg, volume
	Command string            `yaml:"command,omitempty" json:"command,omitempty"`
	Source  string            `yaml:"source,omitempty" json:"source,omitempty"`
	Dest    string            `yaml:"dest,omitempty" json:"dest,omitempty"`
	Env     map[string]string `yaml:"env,omitempty" json:"env,omitempty"`
	Workdir string            `yaml:"workdir,omitempty" json:"workdir,omitempty"`
	User    string            `yaml:"user,omitempty" json:"user,omitempty"`
	Comment string            `yaml:"comment,omitempty" json:"comment,omitempty"`
}

// HealthCheck defines container health check configuration
type HealthCheck struct {
	Type    string `yaml:"type" json:"type"` // http, tcp, cmd
	Path    string `yaml:"path,omitempty" json:"path,omitempty"`
	Port    int    `yaml:"port,omitempty" json:"port,omitempty"`
	Command string `yaml:"command,omitempty" json:"command,omitempty"`
	Interval int    `yaml:"interval,omitempty" json:"interval,omitempty"` // seconds
	Timeout  int    `yaml:"timeout,omitempty" json:"timeout,omitempty"`   // seconds
	Retries int    `yaml:"retries,omitempty" json:"retries,omitempty"`
}

// ParseSpec parses a build specification from JSON
func ParseSpec(data []byte) (*BuildSpec, error) {
	var spec BuildSpec
	if err := json.Unmarshal(data, &spec); err != nil {
		return nil, fmt.Errorf("failed to parse build spec: %w", err)
	}

	// Validate required fields
	if spec.Version == "" {
		spec.Version = "1.0"
	}
	if spec.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if spec.Tag == "" {
		spec.Tag = "latest"
	}
	if spec.Base.Image == "" {
		return nil, fmt.Errorf("base image is required")
	}
	if spec.Base.Tag == "" {
		spec.Base.Tag = "latest"
	}

	return &spec, nil
}

// ParseSpecFromFile reads and parses a build specification from a file
func ParseSpecFromFile(filename string) (*BuildSpec, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return ParseSpec(data)
}

// ToDockerfile converts a BuildSpec to Dockerfile format
func (s *BuildSpec) ToDockerfile() string {
	var lines []string

	// FROM
	baseImage := s.Base.Image
	if s.Base.Tag != "" {
		baseImage += ":" + s.Base.Tag
	}
	lines = append(lines, fmt.Sprintf("FROM %s", baseImage))

	// WORKDIR
	if s.Workdir != "" {
		lines = append(lines, fmt.Sprintf("WORKDIR %s", s.Workdir))
	}

	// ENV (spec level)
	for k, v := range s.Env {
		lines = append(lines, fmt.Sprintf("ENV %s=%s", k, v))
	}

	// USER
	if s.User != "" {
		lines = append(lines, fmt.Sprintf("USER %s", s.User))
	}

	// LABELS
	for k, v := range s.Labels {
		lines = append(lines, fmt.Sprintf("LABEL %s=%s", k, v))
	}

	// EXPOSE
	for _, port := range s.Expose {
		lines = append(lines, fmt.Sprintf("EXPOSE %d", port))
	}

	// Steps
	for _, step := range s.Steps {
		if step.Comment != "" {
			lines = append(lines, fmt.Sprintf("# %s", step.Comment))
		}

		switch strings.ToLower(step.Type) {
		case "run":
			if step.Command != "" {
				lines = append(lines, fmt.Sprintf("RUN %s", step.Command))
			}
		case "copy", "add":
			if step.Source != "" && step.Dest != "" {
				lines = append(lines, fmt.Sprintf("%s %s %s", strings.ToUpper(step.Type), step.Source, step.Dest))
			}
		case "cmd":
			if step.Command != "" {
				lines = append(lines, fmt.Sprintf("CMD %s", step.Command))
			}
		case "arg":
			// ARG handling
			if step.Command != "" {
				lines = append(lines, fmt.Sprintf("ARG %s", step.Command))
			}
		case "volume":
			if step.Dest != "" {
				lines = append(lines, fmt.Sprintf("VOLUME %s", step.Dest))
			}
		}

		// Step-level workdir
		if step.Workdir != "" {
			lines = append(lines, fmt.Sprintf("WORKDIR %s", step.Workdir))
		}

		// Step-level env
		for k, v := range step.Env {
			lines = append(lines, fmt.Sprintf("ENV %s=%s", k, v))
		}

		// Step-level user
		if step.User != "" {
			lines = append(lines, fmt.Sprintf("USER %s", step.User))
		}
	}

	// HEALTHCHECK
	if s.Health != nil {
		switch strings.ToLower(s.Health.Type) {
		case "http":
			if s.Health.Path != "" && s.Health.Port > 0 {
				interval := s.Health.Interval
				if interval == 0 {
					interval = 30
				}
				timeout := s.Health.Timeout
				if timeout == 0 {
					timeout = 10
				}
				retries := s.Health.Retries
				if retries == 0 {
					retries = 3
				}
				lines = append(lines, fmt.Sprintf("HEALTHCHECK --interval=%ds --timeout=%ds --retries=%d CMD curl -f http://localhost:%d%s || exit 1", interval, timeout, retries, s.Health.Port, s.Health.Path))
			}
		case "tcp":
			if s.Health.Port > 0 {
				interval := s.Health.Interval
				if interval == 0 {
					interval = 30
				}
				timeout := s.Health.Timeout
				if timeout == 0 {
					timeout = 10
				}
				retries := s.Health.Retries
				if retries == 0 {
					retries = 3
				}
				lines = append(lines, fmt.Sprintf("HEALTHCHECK --interval=%ds --timeout=%ds --retries=%d CMD nc -z localhost %d || exit 1", interval, timeout, retries, s.Health.Port))
			}
		case "cmd":
			if s.Health.Command != "" {
				interval := s.Health.Interval
				if interval == 0 {
					interval = 30
				}
				timeout := s.Health.Timeout
				if timeout == 0 {
					timeout = 10
				}
				retries := s.Health.Retries
				if retries == 0 {
					retries = 3
				}
				lines = append(lines, fmt.Sprintf("HEALTHCHECK --interval=%ds --timeout=%ds --retries=%d %s", interval, timeout, retries, s.Health.Command))
			}
		}
	}

	return strings.Join(lines, "\n")
}

// Validate validates the build specification
func (s *BuildSpec) Validate() error {
	if s.Name == "" {
		return fmt.Errorf("name is required")
	}
	if s.Base.Image == "" {
		return fmt.Errorf("base image is required")
	}

	// Validate step types
	validTypes := map[string]bool{
		"run": true, "copy": true, "add": true, "cmd": true, "arg": true, "volume": true,
	}
	for i, step := range s.Steps {
		if !validTypes[strings.ToLower(step.Type)] {
			return fmt.Errorf("invalid step type '%s' at step %d", step.Type, i+1)
		}
		if step.Type == "run" && step.Command == "" {
			return fmt.Errorf("RUN step requires a command at step %d", i+1)
		}
		if (step.Type == "copy" || step.Type == "add") && (step.Source == "" || step.Dest == "") {
			return fmt.Errorf("%s step requires source and dest at step %d", strings.ToUpper(step.Type), i+1)
		}
	}

	return nil
}

