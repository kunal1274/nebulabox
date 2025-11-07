package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// FileStore provides persistent file-based storage for the registry
type FileStore struct {
	baseDir   string
	mu        sync.RWMutex
	blobs     map[string]map[string]*blobData    // repo -> digest -> blob
	tags      map[string]map[string]string        // repo -> tag -> digest
	manifests map[string]map[string]*manifestData  // repo -> digest -> manifest
	uploads   map[string]*bytes.Buffer             // uuid -> buffer
	metadata  map[string]*VersionMetadata         // Version metadata per repo
}

// NewFileStore creates a new file-based store
func NewFileStore(baseDir string) (*FileStore, error) {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}
	
	fs := &FileStore{
		baseDir:   baseDir,
		blobs:     make(map[string]map[string]*blobData),
		tags:      make(map[string]map[string]string),
		manifests: make(map[string]map[string]*manifestData),
		uploads:   make(map[string]*bytes.Buffer),
		metadata:  make(map[string]*VersionMetadata),
	}
	
	// Load existing metadata
	if err := fs.loadMetadata(); err != nil {
		// Non-fatal, continue with empty metadata
		fmt.Printf("Warning: failed to load metadata: %v\n", err)
	}
	
	return fs, nil
}

// loadMetadata loads version metadata from disk
func (fs *FileStore) loadMetadata() error {
	metadataFile := filepath.Join(fs.baseDir, "metadata.json")
	
	data, err := os.ReadFile(metadataFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	
	var metadata map[string]*VersionMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return fmt.Errorf("failed to unmarshal metadata: %w", err)
	}
	
	fs.mu.Lock()
	fs.metadata = metadata
	fs.mu.Unlock()
	
	return nil
}

// saveMetadata saves version metadata to disk
func (fs *FileStore) saveMetadata() error {
	fs.mu.RLock()
	metadata := fs.metadata
	fs.mu.RUnlock()
	
	metadataFile := filepath.Join(fs.baseDir, "metadata.json")
	
	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}
	
	if err := os.WriteFile(metadataFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write metadata: %w", err)
	}
	
	return nil
}

// getRepoDir returns the directory path for a repository
func (fs *FileStore) getRepoDir(repo string) string {
	return filepath.Join(fs.baseDir, "repos", repo)
}

// ensureRepoDir ensures the repository directory exists
func (fs *FileStore) ensureRepoDir(repo string) error {
	repoDir := fs.getRepoDir(repo)
	return os.MkdirAll(repoDir, 0755)
}

// GetVersionMetadata returns version metadata for a repository
func (fs *FileStore) GetVersionMetadata(repo string) *VersionMetadata {
	fs.mu.RLock()
	defer fs.mu.RUnlock()
	
	if vm, ok := fs.metadata[repo]; ok {
		return vm
	}
	
	return NewVersionMetadata(repo)
}

// SetVersionMetadata sets version metadata for a repository
func (fs *FileStore) SetVersionMetadata(repo string, vm *VersionMetadata) error {
	fs.mu.Lock()
	fs.metadata[repo] = vm
	fs.mu.Unlock()
	
	return fs.saveMetadata()
}

// ensureRepo ensures repository maps exist
func (fs *FileStore) ensureRepo(repo string) {
	if _, ok := fs.blobs[repo]; !ok {
		fs.blobs[repo] = make(map[string]*blobData)
	}
	if _, ok := fs.tags[repo]; !ok {
		fs.tags[repo] = make(map[string]string)
	}
	if _, ok := fs.manifests[repo]; !ok {
		fs.manifests[repo] = make(map[string]*manifestData)
	}
}

// ListRepositories returns all repository names
func (fs *FileStore) ListRepositories() []string {
	fs.mu.RLock()
	defer fs.mu.RUnlock()
	
	repos := make([]string, 0, len(fs.tags))
	for r := range fs.tags {
		repos = append(repos, r)
	}
	return repos
}

// ListTags returns all tags for a repository
func (fs *FileStore) ListTags(repo string) []string {
	fs.mu.RLock()
	defer fs.mu.RUnlock()
	
	out := []string{}
	if m, ok := fs.tags[repo]; ok {
		for t := range m {
			out = append(out, t)
		}
	}
	return out
}

// GetBlob retrieves a blob
func (fs *FileStore) GetBlob(repo, digest string) (*blobData, bool) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()
	
	if m, ok := fs.blobs[repo]; ok {
		b, ok2 := m[digest]
		return b, ok2
	}
	
	// Try to load from disk
	blobPath := filepath.Join(fs.getRepoDir(repo), "blobs", digest)
	if data, err := os.ReadFile(blobPath); err == nil {
		return &blobData{
			content: data,
			size:    int64(len(data)),
		}, true
	}
	
	return nil, false
}

// PutBlob stores a blob
func (fs *FileStore) PutBlob(repo, digest string, data []byte) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	
	fs.ensureRepo(repo)
	
	blobPath := filepath.Join(fs.getRepoDir(repo), "blobs", digest)
	blobDir := filepath.Dir(blobPath)
	
	if err := os.MkdirAll(blobDir, 0755); err != nil {
		return fmt.Errorf("failed to create blob directory: %w", err)
	}
	
	if err := os.WriteFile(blobPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write blob: %w", err)
	}
	
	fs.blobs[repo][digest] = &blobData{
		content: data,
		size:    int64(len(data)),
	}
	
	return nil
}

// GetManifest retrieves a manifest
func (fs *FileStore) GetManifest(repo, reference string) (*manifestData, bool) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()
	
	// If reference is tag, resolve to digest
	digest := reference
	if len(reference) > 0 && (len(reference) < 7 || reference[:7] != "sha256:") {
		if m, ok := fs.tags[repo]; ok {
			if d, ok2 := m[reference]; ok2 {
				digest = d
			}
		}
	}
	
	if m, ok := fs.manifests[repo]; ok {
		if mf, ok2 := m[digest]; ok2 {
			return mf, true
		}
	}
	
	// Try to load from disk
	manifestPath := filepath.Join(fs.getRepoDir(repo), "manifests", reference)
	if data, err := os.ReadFile(manifestPath); err == nil {
		return &manifestData{
			content:   data,
			mediaType: "application/vnd.oci.image.manifest.v1+json",
		}, true
	}
	
	return nil, false
}

// PutManifest stores a manifest
func (fs *FileStore) PutManifest(repo, reference, digest, mediaType string, content []byte) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	
	fs.ensureRepo(repo)
	fs.manifests[repo][digest] = &manifestData{
		content:   content,
		mediaType: mediaType,
	}
	
	// If reference is a tag, bind tag->digest
	if len(reference) > 0 && (len(reference) < 7 || reference[:7] != "sha256:") {
		fs.tags[repo][reference] = digest
	}
	
	manifestPath := filepath.Join(fs.getRepoDir(repo), "manifests", reference)
	manifestDir := filepath.Dir(manifestPath)
	
	if err := os.MkdirAll(manifestDir, 0755); err != nil {
		return fmt.Errorf("failed to create manifest directory: %w", err)
	}
	
	if err := os.WriteFile(manifestPath, content, 0644); err != nil {
		return fmt.Errorf("failed to write manifest: %w", err)
	}
	
	return nil
}

// DeleteTag deletes a tag
func (fs *FileStore) DeleteTag(repo, tag string) bool {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	
	if m, ok := fs.tags[repo]; ok {
		if _, ok2 := m[tag]; ok2 {
			delete(m, tag)
			
			// Update metadata
			vm := fs.GetVersionMetadata(repo)
			vm.DeleteVersion(tag)
			fs.SetVersionMetadata(repo, vm)
			
			// Delete manifest file if it exists
			manifestPath := filepath.Join(fs.getRepoDir(repo), "manifests", tag)
			os.Remove(manifestPath)
			
			return true
		}
	}
	return false
}

// StartUpload initiates an upload
func (fs *FileStore) StartUpload(uuid string) {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	
	fs.uploads[uuid] = &bytes.Buffer{}
}

// AppendUpload appends data to an upload
func (fs *FileStore) AppendUpload(uuid string, chunk []byte) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	
	buf, ok := fs.uploads[uuid]
	if !ok {
		return fmt.Errorf("upload not found")
	}
	_, _ = buf.Write(chunk)
	return nil
}

// FinalizeUpload finalizes an upload
func (fs *FileStore) FinalizeUpload(repo, uuid, digest, createdBy string, metadata map[string]string, description string) (int64, error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	
	buf, ok := fs.uploads[uuid]
	if !ok {
		return 0, fmt.Errorf("upload not found")
	}
	
	data := buf.Bytes()
	size := int64(len(data))
	
	// Ensure repo exists
	fs.ensureRepo(repo)
	fs.blobs[repo][digest] = &blobData{
		content: data,
		size:    size,
	}
	
	delete(fs.uploads, uuid)
	
	// Save to disk (call internal method to avoid recursion)
	blobPath := filepath.Join(fs.getRepoDir(repo), "blobs", digest)
	blobDir := filepath.Dir(blobPath)
	
	if err := os.MkdirAll(blobDir, 0755); err != nil {
		return 0, fmt.Errorf("failed to create blob directory: %w", err)
	}
	
	if err := os.WriteFile(blobPath, data, 0644); err != nil {
		return 0, fmt.Errorf("failed to write blob: %w", err)
	}
	
	return size, nil
}

