package model

import "time"

type Book struct {
	ID            string    `gorm:"primaryKey;column:id" json:"id"`
	Title         string    `gorm:"column:title" json:"title"`
	Author        string    `gorm:"column:author" json:"author"`
	PublishedDate time.Time `gorm:"column:published_date" json:"published_date"`
	Status        string    `gorm:"column:status" json:"status"`
	UserID        string    `gorm:"column:user_id" json:"user_id"`
}
