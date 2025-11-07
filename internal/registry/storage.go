package registry

import (
    "bytes"
    "fmt"
    "sync"
)

type blobData struct {
    content []byte
    size    int64
}

type manifestData struct {
    content []byte
    mediaType string
}

// memoryStore is a minimal in-memory registry storage
type memoryStore struct {
    mu sync.RWMutex
    // repo -> digest -> blob
    blobs map[string]map[string]*blobData
    // repo -> tag -> manifest digest
    tags map[string]map[string]string
    // repo -> digest -> manifest
    manifests map[string]map[string]*manifestData
    // uploads uuid -> buffer
    uploads map[string]*bytes.Buffer
}

func newMemoryStore() *memoryStore {
    return &memoryStore{
        blobs:     make(map[string]map[string]*blobData),
        tags:      make(map[string]map[string]string),
        manifests: make(map[string]map[string]*manifestData),
        uploads:   make(map[string]*bytes.Buffer),
    }
}

func (s *memoryStore) ensureRepo(repo string) {
    if _, ok := s.blobs[repo]; !ok {
        s.blobs[repo] = make(map[string]*blobData)
    }
    if _, ok := s.tags[repo]; !ok {
        s.tags[repo] = make(map[string]string)
    }
    if _, ok := s.manifests[repo]; !ok {
        s.manifests[repo] = make(map[string]*manifestData)
    }
}

func (s *memoryStore) ListRepositories() []string {
    s.mu.RLock(); defer s.mu.RUnlock()
    repos := make([]string, 0, len(s.tags))
    for r := range s.tags {
        repos = append(repos, r)
    }
    return repos
}

func (s *memoryStore) ListTags(repo string) []string {
    s.mu.RLock(); defer s.mu.RUnlock()
    out := []string{}
    if m, ok := s.tags[repo]; ok {
        for t := range m {
            out = append(out, t)
        }
    }
    return out
}

// Ensure memoryStore implements Storage interface
var _ Storage = (*memoryStore)(nil)

// Add missing methods to memoryStore to implement Storage interface
func (s *memoryStore) PutBlob(repo, digest string, data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.ensureRepo(repo)
	s.blobs[repo][digest] = &blobData{
		content: data,
		size:    int64(len(data)),
	}
	return nil
}

func (s *memoryStore) FinalizeUpload(repo, uuid, digest, createdBy string, metadata map[string]string, description string) (int64, error) {
	size, err := s.finalizeUpload(repo, uuid, digest)
	return size, err
}

func (s *memoryStore) GetVersionMetadata(repo string) *VersionMetadata {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	// Return empty metadata for memory store
	return NewVersionMetadata(repo)
}

func (s *memoryStore) SetVersionMetadata(repo string, vm *VersionMetadata) error {
	// No-op for memory store
	return nil
}

func (s *memoryStore) startUpload(uuid string) {
    s.mu.Lock(); defer s.mu.Unlock()
    s.uploads[uuid] = &bytes.Buffer{}
}

func (s *memoryStore) AppendUpload(uuid string, chunk []byte) error {
    s.mu.Lock(); defer s.mu.Unlock()
    buf, ok := s.uploads[uuid]
    if !ok { return fmt.Errorf("upload not found") }
    _, _ = buf.Write(chunk)
    return nil
}

func (s *memoryStore) finalizeUpload(repo, uuid, digest string) (int64, error) {
    s.mu.Lock(); defer s.mu.Unlock()
    buf, ok := s.uploads[uuid]
    if !ok { return 0, fmt.Errorf("upload not found") }
    s.ensureRepo(repo)
    s.blobs[repo][digest] = &blobData{content: buf.Bytes(), size: int64(buf.Len())}
    delete(s.uploads, uuid)
    return int64(len(s.blobs[repo][digest].content)), nil
}

func (s *memoryStore) GetBlob(repo, digest string) (*blobData, bool) {
    s.mu.RLock(); defer s.mu.RUnlock()
    if m, ok := s.blobs[repo]; ok {
        b, ok2 := m[digest]
        return b, ok2
    }
    return nil, false
}

func (s *memoryStore) PutManifest(repo, reference, digest, mediaType string, content []byte) error {
    s.mu.Lock(); defer s.mu.Unlock()
    s.ensureRepo(repo)
    s.manifests[repo][digest] = &manifestData{content: content, mediaType: mediaType}
    // if reference is a tag, bind tag->digest
    if len(reference) > 0 && (len(reference) < 7 || reference[:7] != "sha256:") {
        s.tags[repo][reference] = digest
    }
    return nil
}

func (s *memoryStore) GetManifest(repo, reference string) (*manifestData, bool) {
    s.mu.RLock(); defer s.mu.RUnlock()
    // if reference is tag, resolve to digest
    digest := reference
    if len(reference) > 0 && (len(reference) < 7 || reference[:7] != "sha256:") {
        if m, ok := s.tags[repo]; ok {
            if d, ok2 := m[reference]; ok2 { digest = d }
        }
    }
    if m, ok := s.manifests[repo]; ok {
        mf, ok2 := m[digest]
        return mf, ok2
    }
    return nil, false
}

func (s *memoryStore) DeleteTag(repo, tag string) bool {
    s.mu.Lock(); defer s.mu.Unlock()
    if m, ok := s.tags[repo]; ok {
        if _, ok2 := m[tag]; ok2 {
            delete(m, tag)
            return true
        }
    }
    return false
}

func (s *memoryStore) StartUpload(uuid string) {
    s.startUpload(uuid)
}


