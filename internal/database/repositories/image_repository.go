package repositories

import (
	"fmt"
	"time"

	"github.com/nebulabox/nebulabox/internal/database"
	"gorm.io/gorm"
)

// ImageData represents image data for repository operations (avoids import cycle)
type ImageData struct {
	ID      string
	Name    string
	Tag     string
	Size    string
	Created string
	Digest  string
}

// ImageRepository handles image database operations
type ImageRepository struct {
	db *gorm.DB
}

// NewImageRepository creates a new image repository
func NewImageRepository(db *gorm.DB) *ImageRepository {
	return &ImageRepository{db: db}
}

// CreateOrUpdate creates or updates an image in the database
func (r *ImageRepository) CreateOrUpdate(image *ImageData) error {
	if r.db == nil {
		return fmt.Errorf("database not initialized")
	}

	// Parse size from string (e.g., "123.45 MB")
	sizeBytes := int64(0)
	if image.Size != "" {
		// Simple parsing - in production, parse properly
		// For now, estimate or store 0
	}

	// Parse created time
	createdTime := time.Now()
	if image.Created != "" {
		if parsed, err := time.Parse(time.RFC3339, image.Created); err == nil {
			createdTime = parsed
		}
	}

	dbImage := &database.Image{
		ID:         image.ID,
		Name:       image.Name,
		Tag:        image.Tag,
		Digest:     image.Digest,
		Size:       sizeBytes,
		Created:    createdTime,
		Repository: image.Name,
	}

	// Set timestamps
	now := time.Now()
	if dbImage.CreatedAt.IsZero() {
		dbImage.CreatedAt = now
	}
	dbImage.UpdatedAt = now

	// Use Save to create or update
	return r.db.Save(dbImage).Error
}

// Get retrieves an image by ID
func (r *ImageRepository) Get(id string) (*ImageData, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	var dbImage database.Image
	if err := r.db.Where("id = ?", id).First(&dbImage).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return r.toImageData(&dbImage), nil
}

// List retrieves all images
func (r *ImageRepository) List() ([]*ImageData, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	var dbImages []database.Image
	if err := r.db.Find(&dbImages).Error; err != nil {
		return nil, err
	}

	images := make([]*ImageData, len(dbImages))
	for i, dbImage := range dbImages {
		images[i] = r.toImageData(&dbImage)
	}

	return images, nil
}

// FindByNameAndTag finds an image by name and tag
func (r *ImageRepository) FindByNameAndTag(name, tag string) (*ImageData, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	var dbImage database.Image
	if err := r.db.Where("name = ? AND tag = ?", name, tag).First(&dbImage).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return r.toImageData(&dbImage), nil
}

// Delete soft deletes an image
func (r *ImageRepository) Delete(id string) error {
	if r.db == nil {
		return fmt.Errorf("database not initialized")
	}

	return r.db.Delete(&database.Image{}, "id = ?", id).Error
}

// toImageData converts database.Image to ImageData
func (r *ImageRepository) toImageData(dbImage *database.Image) *ImageData {
	sizeStr := fmt.Sprintf("%.2f MB", float64(dbImage.Size)/1024/1024)
	if dbImage.Size == 0 {
		sizeStr = "0 MB"
	}

	return &ImageData{
		ID:      dbImage.ID,
		Name:    dbImage.Name,
		Tag:     dbImage.Tag,
		Digest:  dbImage.Digest,
		Size:    sizeStr,
		Created: dbImage.Created.Format(time.RFC3339),
	}
}
