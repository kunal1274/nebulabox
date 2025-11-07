package registry

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"
)

// ImageVersion represents a versioned image tag
type ImageVersion struct {
	Tag         string            `json:"tag"`
	Digest      string            `json:"digest"`
	CreatedAt   time.Time         `json:"createdAt"`
	CreatedBy   string            `json:"createdBy"`
	Size        int64             `json:"size"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	Description string            `json:"description,omitempty"`
}

// VersionMetadata tracks version history and metadata
type VersionMetadata struct {
	Repository string                 `json:"repository"`
	Versions   map[string]*ImageVersion `json:"versions"`
	Latest     string                 `json:"latest"`
	UpdatedAt  time.Time              `json:"updatedAt"`
}

// NewVersionMetadata creates new version metadata
func NewVersionMetadata(repo string) *VersionMetadata {
	return &VersionMetadata{
		Repository: repo,
		Versions:   make(map[string]*ImageVersion),
		UpdatedAt:  time.Now(),
	}
}

// AddVersion adds a new version to the metadata
func (vm *VersionMetadata) AddVersion(tag, digest, createdBy string, size int64, metadata map[string]string, description string) {
	version := &ImageVersion{
		Tag:         tag,
		Digest:      digest,
		CreatedAt:   time.Now(),
		CreatedBy:   createdBy,
		Size:        size,
		Metadata:    metadata,
		Description: description,
	}
	
	vm.Versions[tag] = version
	
	// Update latest if this is a semantic version or 'latest'
	if vm.isLatestTag(tag) {
		vm.Latest = tag
	}
	
	vm.UpdatedAt = time.Now()
}

// GetVersion retrieves version information for a tag
func (vm *VersionMetadata) GetVersion(tag string) (*ImageVersion, bool) {
	v, ok := vm.Versions[tag]
	return v, ok
}

// ListVersions returns all versions sorted by creation time (newest first)
func (vm *VersionMetadata) ListVersions() []*ImageVersion {
	versions := make([]*ImageVersion, 0, len(vm.Versions))
	for _, v := range vm.Versions {
		versions = append(versions, v)
	}
	
	sort.Slice(versions, func(i, j int) bool {
		return versions[i].CreatedAt.After(versions[j].CreatedAt)
	})
	
	return versions
}

// GetLatestVersion returns the latest version
func (vm *VersionMetadata) GetLatestVersion() *ImageVersion {
	if vm.Latest == "" {
		// Find the most recently created version
		versions := vm.ListVersions()
		if len(versions) > 0 {
			return versions[0]
		}
		return nil
	}
	return vm.Versions[vm.Latest]
}

// DeleteVersion removes a version
func (vm *VersionMetadata) DeleteVersion(tag string) bool {
	if _, ok := vm.Versions[tag]; ok {
		delete(vm.Versions, tag)
		
		// If deleted version was latest, update latest
		if vm.Latest == tag {
			vm.Latest = ""
			versions := vm.ListVersions()
			if len(versions) > 0 {
				vm.Latest = versions[0].Tag
			}
		}
		
		vm.UpdatedAt = time.Now()
		return true
	}
	return false
}

// isLatestTag determines if a tag should be considered "latest"
func (vm *VersionMetadata) isLatestTag(tag string) bool {
	tag = strings.ToLower(tag)
	return tag == "latest" || 
		   tag == "main" || 
		   tag == "master" ||
		   strings.HasPrefix(tag, "v") && isSemanticVersion(tag[1:])
}

// isSemanticVersion checks if a string is a semantic version
func isSemanticVersion(v string) bool {
	parts := strings.Split(v, ".")
	if len(parts) != 3 {
		return false
	}
	for _, part := range parts {
		if len(part) == 0 {
			return false
		}
		for _, c := range part {
			if c < '0' || c > '9' {
				return false
			}
		}
	}
	return true
}

// CompareVersions compares two semantic versions
// Returns: -1 if v1 < v2, 0 if v1 == v2, 1 if v1 > v2
func CompareVersions(v1, v2 string) int {
	v1 = strings.TrimPrefix(v1, "v")
	v2 = strings.TrimPrefix(v2, "v")
	
	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")
	
	maxLen := len(parts1)
	if len(parts2) > maxLen {
		maxLen = len(parts2)
	}
	
	for i := 0; i < maxLen; i++ {
		var n1, n2 int
		if i < len(parts1) {
			fmt.Sscanf(parts1[i], "%d", &n1)
		}
		if i < len(parts2) {
			fmt.Sscanf(parts2[i], "%d", &n2)
		}
		
		if n1 < n2 {
			return -1
		}
		if n1 > n2 {
			return 1
		}
	}
	
	return 0
}

// VersionSummary provides a summary of version information
type VersionSummary struct {
	Repository  string   `json:"repository"`
	TotalVersions int    `json:"totalVersions"`
	Latest      string   `json:"latest"`
	Tags        []string `json:"tags"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// GetVersionSummary returns a summary of version metadata
func (vm *VersionMetadata) GetVersionSummary() *VersionSummary {
	tags := make([]string, 0, len(vm.Versions))
	for tag := range vm.Versions {
		tags = append(tags, tag)
	}
	
	sort.Strings(tags)
	
	return &VersionSummary{
		Repository:    vm.Repository,
		TotalVersions: len(vm.Versions),
		Latest:        vm.Latest,
		Tags:          tags,
		UpdatedAt:     vm.UpdatedAt,
	}
}

// MarshalJSON custom marshaling for VersionMetadata
func (vm *VersionMetadata) MarshalJSON() ([]byte, error) {
	type Alias VersionMetadata
	return json.Marshal(&struct {
		*Alias
		VersionsJSON map[string]*ImageVersion `json:"versions"`
	}{
		Alias:        (*Alias)(vm),
		VersionsJSON: vm.Versions,
	})
}

