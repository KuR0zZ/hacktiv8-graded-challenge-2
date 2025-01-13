package model

type User struct {
	ID       string `gorm:"primaryKey;column:id" json:"id"`
	Username string `gorm:"column:username" json:"username"`
	Password string `gorm:"column:password" json:"password"`
}
