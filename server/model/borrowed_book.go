package model

import "time"

type BorrowedBook struct {
	ID           string    `gorm:"primaryKey;column:id" json:"id"`
	BookID       string    `gorm:"column:book_id" json:"book_id"`
	UserID       string    `gorm:"column:user_id" json:"user_id"`
	BorrowedDate time.Time `gorm:"column:borrowed_date" json:"borrowed_date"`
	ReturnDate   time.Time `gorm:"column:return_date" json:"return_date"`
}
