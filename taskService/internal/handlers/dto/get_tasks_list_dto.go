package dto

import "time"

type TaskToList struct {
	ID        int       `json:"id" gorm:"column:id"`
	Title     string    `json:"title" gorm:"column:title"`
	Status    string    `json:"status" gorm:"column:status"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
}
