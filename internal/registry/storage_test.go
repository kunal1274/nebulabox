package registry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorage_ListRepositories(t *testing.T) {
	storage := NewStorage()
	
	repos := storage.ListRepositories()
	assert.NotNil(t, repos)
	assert.IsType(t, []string{}, repos)
}

func TestStorage_AddTag(t *testing.T) {
	storage := NewStorage()
	
	digest := "sha256:abc123"
	err := storage.PutManifest("test/repo", "latest", digest, []byte(`{"config":{}}`), "application/vnd.docker.distribution.manifest.v2+json")
	assert.NoError(t, err)
	
	tags := storage.ListTags("test/repo")
	assert.Contains(t, tags, "latest")
	
	manifest := storage.GetManifest("test/repo", "latest")
	assert.NotNil(t, manifest)
	assert.NotNil(t, manifest.Content)
}

func TestStorage_DeleteTag(t *testing.T) {
	storage := NewStorage()
	
	digest := "sha256:abc123"
	storage.PutManifest("test/repo", "v1", digest, []byte(`{"config":{}}`), "application/vnd.docker.distribution.manifest.v2+json")
	
	tags := storage.ListTags("test/repo")
	assert.Contains(t, tags, "v1")
	
	err := storage.DeleteTag("test/repo", "v1")
	assert.NoError(t, err)
	
	tags = storage.ListTags("test/repo")
	assert.NotContains(t, tags, "v1")
}

func TestStorage_StoreBlob(t *testing.T) {
	storage := NewStorage()
	
	digest := "sha256:def456"
	data := []byte("test blob data")
	storage.PutBlob("test/repo", digest, data)
	
	blob := storage.GetBlob("test/repo", digest)
	assert.NotNil(t, blob)
	assert.Equal(t, data, blob)
}

