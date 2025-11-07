package engine

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ImageManager manages container images
type ImageManager struct {
	storagePath string
	images      map[string]*Image
}

// NewImageManager creates a new image manager
func NewImageManager(storagePath string) (*ImageManager, error) {
	imagesPath := filepath.Join(storagePath, "images")
	if err := os.MkdirAll(imagesPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create images directory: %w", err)
	}

	return &ImageManager{
		storagePath: storagePath,
		images:      make(map[string]*Image),
	}, nil
}

// BuildImage builds an image from a BuildSpec
func (im *ImageManager) BuildImage(ctx context.Context, spec *BuildSpec, runtime *NebulaRuntime) (*Image, error) {
	imageID := generateImageID(spec.Name, spec.Tag)
	imagePath := filepath.Join(im.storagePath, "images", imageID)

	// Create image directory
	if err := os.MkdirAll(imagePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create image directory: %w", err)
	}

	// Build layers from BuildSpec steps
	layers := make([]ImageLayer, 0)
	for i, step := range spec.Steps {
		layer, err := im.buildLayer(ctx, step, imagePath, i)
		if err != nil {
			return nil, fmt.Errorf("failed to build layer %d: %w", i, err)
		}
		layers = append(layers, *layer)
	}

	// Create image manifest
	manifest := &ImageManifest{
		SchemaVersion: 2,
		MediaType:      "application/vnd.oci.image.manifest.v1+json",
		Layers:         make([]*ManifestDescriptor, len(layers)),
	}

	for i, layer := range layers {
		manifest.Layers[i] = &ManifestDescriptor{
			MediaType: layer.MediaType,
			Size:      layer.Size,
			Digest:    layer.Digest,
		}
	}

	// Create image config
	config := &ImageConfig{
		Architecture: "amd64",
		OS:           "linux",
		Config:       make(map[string]interface{}),
		RootFS: &RootFS{
			Type:    "layers",
			DiffIDs: make([]string, len(layers)),
		},
	}

	for i, layer := range layers {
		config.RootFS.DiffIDs[i] = layer.Digest
	}

	// Set environment variables
	if len(spec.Env) > 0 {
		config.Config["Env"] = spec.Env
	}

	// Set working directory
	if spec.Workdir != "" {
		config.Config["WorkingDir"] = spec.Workdir
	}

	// Set exposed ports
	if len(spec.Expose) > 0 {
		config.Config["ExposedPorts"] = spec.Expose
	}

	// Calculate image size
	var totalSize int64
	for _, layer := range layers {
		totalSize += layer.Size
	}

	// Create image object
	image := &Image{
		ID:        imageID,
		Name:      spec.Name,
		Tag:       spec.Tag,
		Digest:    calculateImageDigest(manifest, config),
		Size:      totalSize,
		Layers:    layers,
		Config:    config,
		CreatedAt: time.Now(),
		Manifest:  manifest,
	}

	// Save image metadata
	if err := im.saveImageMetadata(image, imagePath); err != nil {
		return nil, fmt.Errorf("failed to save image metadata: %w", err)
	}

	// Store image
	imageRef := fmt.Sprintf("%s:%s", spec.Name, spec.Tag)
	im.images[imageRef] = image

	return image, nil
}

// buildLayer builds a single layer from a build step
func (im *ImageManager) buildLayer(ctx context.Context, step BuildStep, imagePath string, index int) (*ImageLayer, error) {
	layerID := fmt.Sprintf("layer-%d-%s", index, generateRandomID(12))
	layerPath := filepath.Join(imagePath, layerID)

	if err := os.MkdirAll(layerPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create layer directory: %w", err)
	}

	// Process step based on type
	switch step.Type {
	case "run":
		// Execute command and capture changes
		// For POC, we'll create a placeholder
		if err := im.executeRunStep(step, layerPath); err != nil {
			return nil, err
		}
	case "copy":
		// Copy files
		if err := im.executeCopyStep(step, layerPath); err != nil {
			return nil, err
		}
	case "volume":
		// Create volume mount point
		if err := os.MkdirAll(filepath.Join(layerPath, step.Dest), 0755); err != nil {
			return nil, err
		}
	}

	// Calculate layer size and digest
	size, err := calculateDirectorySize(layerPath)
	if err != nil {
		return nil, err
	}

	digest, err := calculateLayerDigest(layerPath)
	if err != nil {
		return nil, err
	}

	return &ImageLayer{
		ID:        layerID,
		Digest:    digest,
		Size:      size,
		MediaType: "application/vnd.oci.image.layer.v1.tar+gzip",
		Path:      layerPath,
	}, nil
}

// executeRunStep executes a RUN step
func (im *ImageManager) executeRunStep(step BuildStep, layerPath string) error {
	// In a full implementation, we'd:
	// 1. Create a temporary container
	// 2. Execute the command
	// 3. Capture filesystem changes
	// 4. Create a layer from the changes

	// For POC, create a placeholder file
	placeholder := filepath.Join(layerPath, ".run-step")
	return os.WriteFile(placeholder, []byte(step.Command), 0644)
}

// executeCopyStep executes a COPY step
func (im *ImageManager) executeCopyStep(step BuildStep, layerPath string) error {
	destPath := filepath.Join(layerPath, step.Dest)
	destDir := filepath.Dir(destPath)

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Copy source to destination
	// For POC, we'll create a placeholder
	// In full implementation, copy actual files
	placeholder := filepath.Join(layerPath, ".copy-step")
	return os.WriteFile(placeholder, []byte(fmt.Sprintf("%s -> %s", step.Source, step.Dest)), 0644)
}

// saveImageMetadata saves image metadata to disk
func (im *ImageManager) saveImageMetadata(image *Image, imagePath string) error {
	// Save manifest
	manifestPath := filepath.Join(imagePath, "manifest.json")
	manifestData, err := json.MarshalIndent(image.Manifest, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(manifestPath, manifestData, 0644); err != nil {
		return err
	}

	// Save config
	configPath := filepath.Join(imagePath, "config.json")
	configData, err := json.MarshalIndent(image.Config, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(configPath, configData, 0644); err != nil {
		return err
	}

	// Save image metadata
	metadataPath := filepath.Join(imagePath, "image.json")
	imageData, err := json.MarshalIndent(image, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(metadataPath, imageData, 0644)
}

// PullImage pulls an image from a registry
func (im *ImageManager) PullImage(ctx context.Context, imageRef string) error {
	// In a full implementation, we'd:
	// 1. Parse image reference
	// 2. Connect to registry
	// 3. Download manifest
	// 4. Download layers
	// 5. Store locally

	// For POC, create a placeholder image
	image := &Image{
		ID:        generateImageID("pulled", imageRef),
		Name:      imageRef,
		Tag:       "latest",
		Digest:    generateRandomID(64),
		Size:      0,
		Layers:    []ImageLayer{},
		Config:    &ImageConfig{},
		CreatedAt: time.Now(),
	}

	im.images[imageRef] = image
	return nil
}

// GetImage retrieves an image by reference
func (im *ImageManager) GetImage(imageRef string) (*Image, error) {
	image, exists := im.images[imageRef]
	if !exists {
		return nil, fmt.Errorf("image %s not found", imageRef)
	}
	return image, nil
}

// ListImages returns all images
func (im *ImageManager) ListImages() []*Image {
	images := make([]*Image, 0, len(im.images))
	for _, image := range im.images {
		images = append(images, image)
	}
	return images
}

// DeleteImage removes an image
func (im *ImageManager) DeleteImage(imageRef string) error {
	image, exists := im.images[imageRef]
	if !exists {
		return fmt.Errorf("image %s not found", imageRef)
	}

	// Remove image directory
	imagePath := filepath.Join(im.storagePath, "images", image.ID)
	if err := os.RemoveAll(imagePath); err != nil {
		return fmt.Errorf("failed to remove image directory: %w", err)
	}

	delete(im.images, imageRef)
	return nil
}

// Helper functions

func generateImageID(name, tag string) string {
	data := fmt.Sprintf("%s:%s:%d", name, tag, time.Now().UnixNano())
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("sha256:%x", hash[:16])
}

func calculateImageDigest(manifest *ImageManifest, config *ImageConfig) string {
	data, _ := json.Marshal(map[string]interface{}{
		"manifest": manifest,
		"config":   config,
	})
	hash := sha256.Sum256(data)
	return fmt.Sprintf("sha256:%x", hash)
}

func calculateLayerDigest(layerPath string) (string, error) {
	// Calculate digest of layer directory
	// In full implementation, create tar and calculate digest
	hash := sha256.Sum256([]byte(layerPath))
	return fmt.Sprintf("sha256:%x", hash[:16]), nil
}

func calculateDirectorySize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
}

func generateRandomID(length int) string {
	// Simplified ID generation
	// In production, use crypto/rand
	return "abc123def456"[:length]
}

