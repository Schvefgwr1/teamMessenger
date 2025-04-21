package models

type FileType struct {
	ID   int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"type:varchar(255);not null" json:"name"`
}

func (FileType) TableName() string {
	return "file_service.file_types"
}
