package models

import "time"

type File struct {
	ID         int       `gorm:"primaryKey;serial" json:"id"`
	Name       string    `gorm:"type:varchar(255);not null" json:"name"`
	FileTypeID int       `gorm:"column:file_type;not null" json:"file_type_id"`
	URL        string    `gorm:"type:text;not null" json:"url"`
	CreatedAt  time.Time `gorm:"default:now()" json:"created_at"`
	FileType   FileType  `gorm:"foreignKey:FileTypeID" json:"file_type"`
}

func (File) TableName() string {
	return "file_service.files"
}
