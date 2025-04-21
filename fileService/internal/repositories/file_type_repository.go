package repositories

import (
	"fileService/internal/models"
	"gorm.io/gorm"
)

type FileTypeRepository struct {
	db *gorm.DB
}

func NewFileTypeRepository(db *gorm.DB) *FileTypeRepository {
	return &FileTypeRepository{db: db}
}

// CreateFileType создает новый тип файла
func (r *FileTypeRepository) CreateFileType(fileType *models.FileType) error {
	return r.db.Create(fileType).Error
}

// GetFileTypeByID получает тип файла по ID
func (r *FileTypeRepository) GetFileTypeByID(id int) (*models.FileType, error) {
	var fileType models.FileType
	err := r.db.First(&fileType, "id = ?", id).Error
	return &fileType, err
}

// GetFileTypeByName получает тип файла по названию
func (r *FileTypeRepository) GetFileTypeByName(name string) (*models.FileType, error) {
	var fileType models.FileType
	err := r.db.First(&fileType, "name = ?", name).Error
	return &fileType, err
}

// DeleteFileTypeByID удаляет тип файла по ID
func (r *FileTypeRepository) DeleteFileTypeByID(id int) error {
	return r.db.Delete(&models.FileType{}, id).Error
}
