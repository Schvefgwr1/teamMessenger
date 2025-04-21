package repositories

import (
	"fileService/internal/dto"
	"fileService/internal/models"
	"gorm.io/gorm"
)

type FileRepository struct {
	db *gorm.DB
}

func NewFileRepository(db *gorm.DB) *FileRepository {
	return &FileRepository{db: db}
}

func (r *FileRepository) CreateFile(file *models.File) error {
	return r.db.Create(file).Error
}

func (r *FileRepository) GetFileByID(id int) (*models.File, error) {
	var file models.File
	err := r.db.Preload("FileType").First(&file, "id = ?", id).Error
	return &file, err
}

func (r *FileRepository) GetFileByName(name string) (*models.File, error) {
	var file models.File
	err := r.db.Preload("FileType").First(&file, "name = ?", name).Error
	return &file, err
}

func (r *FileRepository) UpdateFile(file *models.File) error {
	return r.db.Omit("FileType").Save(file).Error
}

func (r *FileRepository) GetFileNamesWithPagination(limit, offset int) (*[]dto.FileInformation, error) {
	var results []dto.FileInformation
	err := r.db.Model(&models.File{}).
		Select("id, name").
		Limit(limit).
		Offset(offset).
		Find(&results).Error
	return &results, err
}
